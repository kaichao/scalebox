package semaphore

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"

	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/common"
	"github.com/kaichao/scalebox/pkg/postgres"
	"github.com/sirupsen/logrus"
)

// AddValue ...
func AddValue(name string, appID int, delta int) (v int, err error) {
	sqlText := `
		UPDATE t_semaphore
		SET value = value + $3
		WHERE name = $1 AND app = $2 AND vtask IS NULL
		RETURNING value
	`
	err = postgres.GetDB().QueryRow(sqlText,
		name, appID, delta).Scan(&v)
	logrus.Tracef("In semaphore.AddValue(),name=%s,app-id:%d,delta:%d,ret-value:%d,err:%v\n",
		name, appID, delta, v, err)

	if err == nil {
		return v, nil
	}
	if err != sql.ErrNoRows {
		return v, errors.WrapE(err, "update semaphore",
			"app-id", appID, "sema-name", name, "delta", delta)
	}
	// not-defined semaphore
	if os.Getenv("SEMAPHORE_AUTO_CREATE") == "yes" {
		// create semaphore first time
		if err := Create(name, delta, appID); err != nil {
			return -1, errors.WrapE(err, "create semaphore",
				"app-id", appID, "sema-name", name, "delta", delta)
		}
		return 0, nil
	}
	return -1, errors.WrapE(err, "semaphore not found",
		"app-id", appID, "sema-name", name, "delta", delta)
}

// AddRegexValue ...
// 根据name正则表达式前缀匹配，如需精确匹配，name末尾加上$
func AddRegexValue(name string, appID int, delta int) (v string, err error) {
	sqlText := `
		WITH updated_rows AS (
			UPDATE t_semaphore
			SET value = value + $3
			WHERE (name ~ $1) AND app = $2 AND vtask IS NULL
			RETURNING name,value
		)
		SELECT COALESCE(JSON_OBJECT_AGG(name, value), '{}') AS aggregated_values
		FROM updated_rows
	`

	if !common.IsRegexString(name[0:1]) {
		// 首字母不是regex元字符，自动添加^
		name = "^" + name
	}
	err = postgres.GetDB().QueryRow(sqlText, name, appID, delta).Scan(&v)
	v = regexp.MustCompile(`\s+`).ReplaceAllString(v, "")

	logrus.Tracef("In semaphore.AddRegexValue(),name=%s,app-id:%d,delta:%d,ret-value:%s,err:%v\n",
		name, appID, delta, v, err)
	return v, errors.WrapE(err, "semaphore-op failed",
		"app-id", appID, "sema-expr", name, "delta", delta)
}

// AddMapValues ...
// 用一条sql语句，或者用一个transaction，完成以下功能。如果更新出错，报错。
// pairs中存放着name、delta的对应值，delta是value的增减值。
// 返回结果为修改后的name及最终值。
func AddMapValues(pairs map[string]int, appID int) (map[string]int, error) {
	if len(pairs) == 0 {
		return map[string]int{}, nil
	}

	// 提取names和deltas数组
	names := make([]string, 0, len(pairs))
	deltas := make([]int, 0, len(pairs))
	for name, delta := range pairs {
		names = append(names, name)
		deltas = append(deltas, delta)
	}

	sqlText := `
		WITH data AS (
			SELECT name, delta
			FROM unnest($1::text[], $2::int[]) AS t(name, delta)
		),
		updated_rows AS (
			UPDATE t_semaphore s
			SET value = s.value + d.delta
			FROM data d
			WHERE s.name = d.name AND s.app = $3 AND s.vtask IS NULL
			RETURNING s.name, s.value
		)
		SELECT COALESCE(JSON_OBJECT_AGG(name, value), '{}') AS aggregated_values,
		   COUNT(*) as updated_count
		FROM updated_rows
	`

	// 使用一条SQL语句同时检查和更新
	// 通过检查更新的行数是否等于输入的数量来判断是否有name不存在
	var v string
	var updatedCount int
	err := postgres.GetDB().QueryRow(sqlText, names, deltas, appID).Scan(&v, &updatedCount)
	logrus.Tracef("In semaphore.AddMapValues(),pairs=%v,app-id:%d,,ret-value:%s,err:%v\n",
		pairs, appID, v, err)

	if err != nil {
		return map[string]int{}, errors.WrapE(err, "semaphore-op failed",
			"app-id", appID, "sema-pairs", pairs)
	}

	// 检查更新的行数是否等于输入的数量
	if updatedCount != len(names) {
		// 有些信号量不存在
		if os.Getenv("SEMAPHORE_AUTO_CREATE") == "yes" {
			// 自动创建不存在的信号量
			for name, delta := range pairs {
				// 尝试创建信号量（如果已经存在，Create函数会更新）
				if err := Create(name, delta, appID); err != nil {
					return map[string]int{}, errors.WrapE(err, "semaphore create",
						"app-id", appID, "sema-name", name, "delta", delta)
					// 继续尝试创建其他信号量
				}
			}

			// 重试更新
			err = postgres.GetDB().QueryRow(sqlText, names, deltas, appID).Scan(&v, &updatedCount)

			if err != nil {
				err = errors.WrapE(err, "semaphore-op failed",
					"app-id", appID, "sema-names", names, "deltas", deltas, "pairs", pairs)
				return map[string]int{}, err
			}

			if updatedCount != len(names) {
				err := errors.E("semaphores still not found after auto-create",
					"app-id", appID, "pairs", pairs, "updated", updatedCount, "expected", len(names))
				return map[string]int{}, err
			}
		} else {
			return map[string]int{}, errors.E("semaphores not found",
				"app-id", appID, "pairs", pairs, "updated", updatedCount, "expected", len(names))
		}
	}

	// 删除空字符
	packed := regexp.MustCompile(`\s+`).ReplaceAllString(v, "")
	if packed == "{}" {
		// 这种情况不应该发生，因为updatedCount > 0
		err := errors.E("Unexpected empty result after update semaphores",
			"app-id", appID, "pairs", pairs, "updated")
		return map[string]int{}, err
	}

	// 解析JSON到map
	result := make(map[string]int)
	// 使用正则表达式提取键值对
	re := regexp.MustCompile(`"([^"]+)":(-?[0-9]+)`)
	matches := re.FindAllStringSubmatch(packed, -1)
	for _, match := range matches {
		if len(match) == 3 {
			name := match[1]
			var value int
			fmt.Sscanf(match[2], "%d", &value)
			result[name] = value
		}
	}

	return result, nil
}
