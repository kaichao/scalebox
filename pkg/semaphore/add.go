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
func AddValue(name string, appID int, delta int) (string, error) {
	sqlFmt := `
		WITH updated_rows AS (
    		UPDATE t_semaphore
    		SET value = value + $3
    		WHERE (name %s $1) AND app = $2
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
	sqlText := fmt.Sprintf(sqlFmt, op)
	v := ""
	if err := postgres.GetDB().QueryRow(sqlText, name, appID, delta).Scan(&v); err != nil {
		errInfo := fmt.Sprintf("[ERROR]db-error in semaphore-op (%s,%d), err-t=%T,err=%v",
			name, appID, err, err)
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
			if err := Create(name, delta, appID); err != nil {
				logrus.Errorf(" Semaphore (name:%s,app-id:%d), create error,err-info:%v\n",
					name, appID, err)
				return "", err
			}
			return fmt.Sprintf("%d", delta), nil
		}
		errInfo := fmt.Sprintf("[ERROR]Semaphore(name:%s,app-id:%d) not-found", name, appID)
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

// AddListValue ...
func AddListValue(names []string, appID int, delta int) (string, error) {
	sqlText := `
		WITH updated_rows AS (
    		UPDATE t_semaphore
    		SET value = value + $3
    		WHERE name = ANY($1) AND app = $2
    		RETURNING name,value
		)
		SELECT COALESCE(JSON_OBJECT_AGG(name, value), '{}') AS aggregated_values
		FROM updated_rows
	`
	v := ""
	if err := postgres.GetDB().QueryRow(sqlText, names, appID, delta).Scan(&v); err != nil {
		errInfo := fmt.Sprintf("[ERROR]db-error in semaphore-op (%s,%d), err-t=%T,err=%v",
			names, appID, err, err)
		logrus.Errorln(errInfo)
		return "", err
	}

	return regexp.MustCompile(`\s+`).ReplaceAllString(v, ""), nil
}

// AddMultiValues ...
// 用一条sql语句，或者用一个transaction，完成以下功能。如果更新出错，报错。
// pairs中存放着name、delta的对应值，delta是value的增减值。
// 返回结果为修改后的name及最终值。
func AddMultiValues(pairs map[string]int, appID int) (map[string]int, error) {
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

	// 使用一条SQL语句同时检查和更新
	// 通过检查更新的行数是否等于输入的数量来判断是否有name不存在
	sqlText := `
		WITH data AS (
			SELECT name, delta
			FROM unnest($1::text[], $2::int[]) AS t(name, delta)
		),
		updated_rows AS (
			UPDATE t_semaphore s
			SET value = s.value + d.delta
			FROM data d
			WHERE s.name = d.name AND s.app = $3
			RETURNING s.name, s.value
		)
		SELECT COALESCE(JSON_OBJECT_AGG(name, value), '{}') AS aggregated_values,
		       COUNT(*) as updated_count
		FROM updated_rows
	`
	var v string
	var updatedCount int
	if err := postgres.GetDB().QueryRow(sqlText, names, deltas, appID).Scan(&v, &updatedCount); err != nil {
		errInfo := fmt.Sprintf("[ERROR]db-error in semaphore-op (%v,%d), err-t=%T,err=%v",
			pairs, appID, err, err)
		logrus.Errorln(errInfo)
		return map[string]int{}, err
	}

	// 检查更新的行数是否等于输入的数量
	if updatedCount != len(names) {
		errInfo := fmt.Sprintf("[ERROR]Some semaphores not found for pairs:%v,app-id:%d, updated:%d, expected:%d",
			pairs, appID, updatedCount, len(names))
		logrus.Errorln(errInfo)
		return map[string]int{}, errors.New(errInfo)
	}

	// 解析JSON字符串
	packed := regexp.MustCompile(`\s+`).ReplaceAllString(v, "")
	if packed == "{}" {
		// 这种情况不应该发生，因为updatedCount > 0
		errInfo := fmt.Sprintf("[ERROR]Unexpected empty result after update for pairs:%v,app-id:%d",
			pairs, appID)
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
