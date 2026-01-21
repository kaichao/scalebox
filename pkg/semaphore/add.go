package semaphore

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/kaichao/scalebox/pkg/common"
	"github.com/kaichao/scalebox/pkg/postgres"
	"github.com/sirupsen/logrus"
)

// AddValue ...
// 根据name或正则表达式前缀匹配，如需精确匹配，name末尾加上$
func AddValue(name string, vtaskID int64, appID int, delta int) (v string, err error) {
	// 构建SQL查询，考虑vtaskID参数
	sqlFmt := `
		WITH updated_rows AS (
			UPDATE t_semaphore
			SET value = value + $3
			WHERE (name %s $1) AND app = $2 AND %s
			RETURNING name,value
		)
		SELECT COALESCE(JSON_OBJECT_AGG(name, value), '{}') AS aggregated_values
		FROM updated_rows
	`

	op := "="
	if common.IsRegexString(name) {
		op = "~"
		if !common.IsRegexString(name[0:1]) {
			// 首字母不是regex元字符，自动添加^
			name = "^" + name
		}
	}
	if vtaskID > 0 {
		// vtaskID > 0 时，需匹配vtask参数
		vtaskExpr := "vtask = $4"
		err = postgres.GetDB().QueryRow(fmt.Sprintf(sqlFmt, op, vtaskExpr),
			name, appID, delta, vtaskID).Scan(&v)
	} else {
		vtaskExpr := "vtask IS NULL"
		err = postgres.GetDB().QueryRow(fmt.Sprintf(sqlFmt, op, vtaskExpr),
			name, appID, delta).Scan(&v)
	}
	logrus.Debugf("In semaphore.AddValue(),name=%s,vtask-id:%d,app-id:%d,delta:%d,ret-value:%s,err:%v\n",
		name, vtaskID, appID, v, err)

	if err != nil {
		errInfo := fmt.Sprintf("[ERROR]db-error in semaphore-op (%s,%d,vtask:%d), err-t=%T,err=%v",
			name, appID, vtaskID, err, err)
		logrus.Errorln(errInfo)
		return "", err
	}

	packed := regexp.MustCompile(`\s+`).ReplaceAllString(v, "")
	if common.IsRegexString(name) {
		// regex name, return json-value
		return packed, nil
	}

	if packed == "{}" {
		// not-defined semaphore
		if os.Getenv("SEMAPHORE_AUTO_CREATE") == "yes" {
			// create semaphore first time
			if err := Create(name, delta, vtaskID, appID); err != nil {
				logrus.Errorf(" Semaphore (name:%s,app-id:%d,vtask:%d), create error,err-info:%v\n",
					name, appID, vtaskID, err)
				return "", err
			}
			return fmt.Sprintf("%d", delta), nil
		}
		errInfo := fmt.Sprintf("[ERROR]Semaphore(name:%s,app-id:%d,vtask:%d) not-found", name, appID, vtaskID)
		logrus.Errorln(errInfo)
		return "", errors.New(errInfo)
	}
	// non-regex, not null, return int-value
	re := regexp.MustCompile(`{".+":(-?[0-9]+)}`)
	ss := re.FindStringSubmatch(packed)
	if len(ss) == 0 {
		errInfo := fmt.Sprintf("[ERROR]Invalid JSON string, value=%s", packed)
		logrus.Errorln(errInfo)
		return "", errors.New(errInfo)
	}
	return ss[1], nil
}

// AddMultiValues ...
// 用一条sql语句，或者用一个transaction，完成以下功能。如果更新出错，报错。
// pairs中存放着name、delta的对应值，delta是value的增减值。
// 返回结果为修改后的name及最终值。
func AddMultiValues(pairs map[string]int, vtaskID int64, appID int) (map[string]int, error) {
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

	sqlFmt := `
		WITH data AS (
			SELECT name, delta
			FROM unnest($1::text[], $2::int[]) AS t(name, delta)
		),
		updated_rows AS (
			UPDATE t_semaphore s
			SET value = s.value + d.delta
			FROM data d
			WHERE s.name = d.name AND s.app = $3 AND %s
			RETURNING s.name, s.value
		)
		SELECT COALESCE(JSON_OBJECT_AGG(name, value), '{}') AS aggregated_values,
		   COUNT(*) as updated_count
		FROM updated_rows
	`

	vtaskExpr := "s.vtask IS NULL"
	if vtaskID > 0 {
		vtaskExpr = "s.vtask = $4"
	}
	sqlText := fmt.Sprintf(sqlFmt, vtaskExpr)

	// 使用一条SQL语句同时检查和更新
	// 通过检查更新的行数是否等于输入的数量来判断是否有name不存在
	var v string
	var err error
	var updatedCount int
	if vtaskID <= 0 {
		// vtaskID <= 0 时，查询vtask IS NULL的记录
		err = postgres.GetDB().QueryRow(sqlText, names, deltas, appID).Scan(&v, &updatedCount)
	} else {
		// vtaskID > 0 时，需要匹配vtask参数
		err = postgres.GetDB().QueryRow(sqlText, names, deltas, appID, vtaskID).Scan(&v, &updatedCount)
	}
	logrus.Debugf("In semaphore.AddMultiValues(),pairs=%v,vtask-id:%d,app-id:%d,,ret-value:%s,err:%v\n",
		pairs, vtaskID, appID, v, err)

	if err != nil {
		errInfo := fmt.Sprintf("[ERROR]db-error in semaphore-op (%v,%d,vtask:%d), err-t=%T,err=%v",
			pairs, appID, vtaskID, err, err)
		logrus.Errorln(errInfo)
		return map[string]int{}, err
	}

	// 检查更新的行数是否等于输入的数量
	if updatedCount != len(names) {
		// 有些信号量不存在
		if os.Getenv("SEMAPHORE_AUTO_CREATE") == "yes" {
			// 自动创建不存在的信号量
			for name, delta := range pairs {
				// 尝试创建信号量（如果已经存在，Create函数会更新）
				if err := Create(name, delta, vtaskID, appID); err != nil {
					logrus.Errorf(" Semaphore (name:%s,app-id:%d,vtask:%d), create error,err-info:%v\n",
						name, appID, vtaskID, err)
					// 继续尝试创建其他信号量
				}
			}

			// 重试更新
			if vtaskID <= 0 {
				err = postgres.GetDB().QueryRow(sqlText, names, deltas, appID).Scan(&v, &updatedCount)
			} else {
				err = postgres.GetDB().QueryRow(sqlText, names, deltas, appID, vtaskID).Scan(&v, &updatedCount)
			}

			if err != nil {
				errInfo := fmt.Sprintf("[ERROR]db-error in semaphore-op after auto-create (%v,%d,vtask:%d), err-t=%T,err=%v",
					pairs, appID, vtaskID, err, err)
				logrus.Errorln(errInfo)
				return map[string]int{}, err
			}

			if updatedCount != len(names) {
				errInfo := fmt.Sprintf("[ERROR]Some semaphores still not found after auto-create for pairs:%v,app-id:%d,vtask:%d, updated:%d, expected:%d",
					pairs, appID, vtaskID, updatedCount, len(names))
				logrus.Errorln(errInfo)
				return map[string]int{}, errors.New(errInfo)
			}
		} else {
			errInfo := fmt.Sprintf("[ERROR]Some semaphores not found for pairs:%v,app-id:%d,vtask:%d, updated:%d, expected:%d",
				pairs, appID, vtaskID, updatedCount, len(names))
			logrus.Errorln(errInfo)
			return map[string]int{}, errors.New(errInfo)
		}
	}

	// 解析JSON字符串
	packed := regexp.MustCompile(`\s+`).ReplaceAllString(v, "")
	if packed == "{}" {
		// 这种情况不应该发生，因为updatedCount > 0
		errInfo := fmt.Sprintf("[ERROR]Unexpected empty result after update for pairs:%v,app-id:%d,vtask:%d",
			pairs, appID, vtaskID)
		logrus.Errorln(errInfo)
		return map[string]int{}, errors.New(errInfo)
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
