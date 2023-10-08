package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	scalebox "github.com/kaichao/scalebox/golang/misc"
	"github.com/sirupsen/logrus"
)

func (dataset *DataSet) parseEntity(s string) *Entity {
	fmt.Println("KeyGroupRegex:", dataset.KeyGroupRegex)
	regex := regexp.MustCompile(dataset.KeyGroupRegex)
	var idx []int
	for _, s := range strings.Split(dataset.KeyGroupIndex, ",") {
		if i, err := strconv.Atoi(s); err == nil {
			idx = append(idx, i)
		}
	}

	if !regex.MatchString(s) {
		return nil
	}

	var kg0, kg1, kg2 string
	ss := regex.FindStringSubmatch(s)
	if len(idx) == 0 {
		// no key_group_index
		if len(ss) > 1 {
			kg0 = ss[1]
		}
	} else {
		kg0 = ss[idx[0]]
	}
	if len(idx) > 1 && len(ss) > idx[1] {
		kg1 = ss[idx[1]]
	}
	if len(idx) > 2 && len(ss) > idx[2] {
		kg2 = ss[idx[2]]
	}

	entity := &Entity{
		name:      s,
		datasetID: kg0,
		x:         kg1,
		y:         kg2,
	}

	return entity
}

// addEntity ...
func (dataset *DataSet) addEntity(entity *Entity) bool {
	sqlText := `
		INSERT INTO t_entity(name,dataset_id,x,y)
		VALUES($1,$2,$3,$4)
	`
	_, err := db.Exec(sqlText, entity.name, dataset.DatasetID,
		entity.x, entity.y)
	if err != nil {
		logrus.Errorf("add entity, err=%v\n", err)
		return false
	}

	return true
}

// isNewGroupAvailable ...
func (dataset *DataSet) checkNewGroupAvailable(entity *Entity) bool {
	var cnt int
	var err error
	if dataset.GroupType == "H" {
		sqlText := `
			SELECT count(*)
			FROM t_entity
			WHERE dataset_id=$1 AND y=$2
		`
		err := db.QueryRow(sqlText, dataset.DatasetID, entity.y).Scan(&cnt)
		if err != nil {
			logrus.Errorf("sum entity, err=%v\n", err)
			return false
		}
		logrus.Println("count=", cnt)
		printSqlite()
		if cnt == dataset.HorizontalWidth {
			return true
		}
	} else {
		sqlText := `
			SELECT count(*)
			FROM t_entity
			WHERE dataset_id=$1 AND x=$2 AND (y BETWEEN $3 AND $4)
		`
		x, _ := strconv.Atoi(entity.x)
		y, _ := strconv.Atoi(entity.y)
		arr := dataset.getVerticalGroupRange(y)
		for i := 0; i < len(arr)-1; i++ {
			n0 := arr[i]
			n1 := arr[i+1]
			err = db.QueryRow(sqlText, dataset.DatasetID, x, n0, n1).Scan(&cnt)
			fmt.Printf("x=%d,y0=%d,y1=%d,cnt=%d\n", x, n0, n1, cnt)
			if err != nil {
				logrus.Errorf("sum entity, err=%v\n", err)
				return false
			}
			length := n1 - n0 + 1
			if cnt == length {
				return true
			}
		}
	}

	return false
}

func printSqlite() {
	sqlText := `
	SELECT name,dataset_id,y
	FROM t_entity
	`
	rows, err := db.Query(sqlText)
	if err != nil {
		logrus.Printf("err:%v\n", err)
	}
	for rows.Next() {
		var (
			name, dataset string
			y             int
		)
		if rows.Scan(&name, &dataset, &y); err == nil {
			fmt.Printf("name:%s,dataset:%s,y:%d\n", name, dataset, y)
		}
	}
}

func (dataset *DataSet) getVerticalGroupRange(y int) []int {
	var n0, n1 int
	y -= dataset.VerticalStart
	if dataset.Interleaved {
		m := y % (dataset.GroupLength - 1)
		if m == 0 {
			if y == 0 {
				n0 = 0
				n1 = dataset.GroupLength - 1
			} else if y == dataset.VerticalLength-1 {
				n1 = y
				n0 = y - dataset.GroupLength + 1
			} else { // 2 ranges
				n1 = y
				n0 = y - dataset.GroupLength + 1
				n2 := y + dataset.GroupLength - 1
				if n2 > dataset.VerticalLength-1 {
					n2 = dataset.VerticalLength - 1
				}
				return []int{dataset.VerticalStart + n0, dataset.VerticalStart + n1, dataset.VerticalStart + n2}
			}
		} else {
			n0 = y - y%(dataset.GroupLength-1)
			n1 = n0 + dataset.GroupLength - 1
		}
	} else {
		// non-interleaved
		n0 = y - y%dataset.GroupLength
		n1 = n0 + dataset.GroupLength - 1
	}
	if n1 > dataset.VerticalLength-1 {
		n1 = dataset.VerticalLength - 1
	}
	return []int{dataset.VerticalStart + n0, dataset.VerticalStart + n1}
}

// outputNewGroups ...
func (dataset *DataSet) outputNewGroups(entity *Entity) {
	var txt string
	var err error
	if dataset.GroupType == "H" {
		sqlText := `
			SELECT GROUP_CONCAT(name)
			FROM (
				SELECT name
				FROM t_entity
				WHERE dataset_id=$1 AND y=$2
				ORDER BY 1
			)
		`
		err = db.QueryRow(sqlText, dataset.DatasetID, entity.y).Scan(&txt)
		if err != nil {
			logrus.Errorf("sum entity, err=%v\n", err)
		} else {
			scalebox.AppendToFile(messageFile, dataset.SinkJob+","+txt)
		}
	} else {
		sqlText := `
			SELECT GROUP_CONCAT(name)
			FROM (
				SELECT name 
				FROM t_entity
				WHERE dataset_id=$1 AND x=$2 AND (y BETWEEN $3 AND $4)
				ORDER BY 1
			)
		`
		x, _ := strconv.Atoi(entity.x)
		y, _ := strconv.Atoi(entity.y)
		arr := dataset.getVerticalGroupRange(y)
		for i := 0; i < len(arr)-1; i++ {
			n0 := arr[i]
			n1 := arr[i+1]

			err = db.QueryRow(sqlText, dataset.DatasetID, x, n0, n1).Scan(&txt)
			if err != nil {
				logrus.Errorf("sum entity, err=%v\n", err)
			} else {
				scalebox.AppendToFile(messageFile, dataset.SinkJob+","+txt)
			}
		}
	}
}

func (dataset *DataSet) loadExistedFiles() {
}

func parseDataSet(t string) *DataSet {
	var ds DataSet
	if err := json.Unmarshal([]byte(t), &ds); err != nil {
		// skip non-json format error
		if !strings.HasPrefix(err.Error(), "invalid character") {
			fmt.Printf("error parsing, err-info:%v\n", err)
		}
		// non-dataset definition
		return nil
	}
	ds.DatasetID = datasetPrefix + ":" + ds.DatasetID
	return &ds
}
