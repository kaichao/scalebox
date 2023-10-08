package main

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataSetParse(t *testing.T) {
	// 去除所有空字符
	s := regexp.MustCompile("\\s+").ReplaceAllString(jsonDataSet0, "")
	ds := parseDataSet(s)
	assert.Equal(t, dataset0, ds)

	s = regexp.MustCompile("\\s+").ReplaceAllString(jsonDataSet1, "")
	ds = parseDataSet(s)
	assert.Equal(t, dataset1, ds)

	s = regexp.MustCompile("\\s+").ReplaceAllString(jsonDataSet2, "")
	ds = parseDataSet(s)
	assert.Equal(t, dataset2, ds)
}

func TestEntityParse(t *testing.T) {
	s := "ds0/my_path_1234/M01.txt"
	entity := &Entity{
		name:      s,
		datasetID: "ds0",
		x:         "1",
		y:         "1234",
	}

	e := dataset0.parseEntity(s)
	assert.Equal(t, entity, e)
}
