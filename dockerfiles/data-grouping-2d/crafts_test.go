package main

import (
	"fmt"
	"os"
	"regexp"
	"testing"
)

func TestCraftsParse(t *testing.T) {

	os.Setenv("COORD_TYPE", "integer")
	// 去除所有空字符
	dss := regexp.MustCompile("\\s+").ReplaceAllString(jsonDataSetCraftsFits, "")
	dataset := parseDataSet(dss)

	s := "fits,Dec+6007_09_03/20221019/Dec+6007_09_03_arcdrift-M04_0001.fits"
	entity := dataset.parseEntity(s)
	dataset.addEntity(entity)
	ss := dataset.getNewGroups(entity)
	fmt.Println("new-groups:", ss)

	s = "fits,Dec+6007_09_03/20221019/Dec+6007_09_03_arcdrift-M04_0002.fits"
	entity = dataset.parseEntity(s)
	dataset.addEntity(entity)
	ss = dataset.getNewGroups(entity)
	fmt.Println("new-groups:", ss)

	s = "fits,Dec+6007_09_03/20221019/Dec+6007_09_03_arcdrift-M04_0003.fits"
	entity = dataset.parseEntity(s)
	dataset.addEntity(entity)
	ss = dataset.getNewGroups(entity)
	fmt.Println("new-groups:", ss)

	s = "fits,Dec+6007_09_03/20221019/Dec+6007_09_03_arcdrift-M04_0004.fits"
	entity = dataset.parseEntity(s)
	dataset.addEntity(entity)
	ss = dataset.getNewGroups(entity)
	fmt.Println("new-groups:", ss)
}
