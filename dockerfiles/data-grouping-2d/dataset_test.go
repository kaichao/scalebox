package main

import (
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestCheckNewGroupAvailable_Horizonal(t *testing.T) {
	s := "ds0/my_path_1234/M01.txt"
	e := dataset0.parseEntity(s)
	dataset0.addEntity(e)
	assert.Equal(t, false, dataset0.checkNewGroupAvailable(e))

	s = "ds0/my_path_1234/M02.txt"
	e = dataset0.parseEntity(s)
	dataset0.addEntity(e)
	assert.Equal(t, false, dataset0.checkNewGroupAvailable(e))

	s = "ds0/my_path_1234/M03.txt"
	e = dataset0.parseEntity(s)
	dataset0.addEntity(e)
	assert.Equal(t, true, dataset0.checkNewGroupAvailable(e))
}

func TestOutputNewGroups_Horizontal(t *testing.T) {
	fmt.Println("dataset0:", dataset0)

	names := []string{
		"ds0/my_path_1234/M03.txt",
		"ds0/my_path_1234/M02.txt",
		"ds0/my_path_1234/M01.txt",
	}
	for _, name := range names {
		entity := dataset0.parseEntity(name)
		dataset0.addEntity(entity)
		if dataset0.checkNewGroupAvailable(entity) {
			dataset0.outputNewGroups(entity)
		}
	}
}

func TestGetVerticalGroupRange(t *testing.T) {
	fmt.Println("dataset2:", dataset2)
	// y := 5
	for y := 1; y <= 10; y++ {
		arr := dataset2.getVerticalGroupRange(y)
		fmt.Println(y, arr)
	}
}

func TestOutputNewGroups_Vertical(t *testing.T) {
	fmt.Println("dataset1:", dataset0)

	names := []string{
		"ds1/my_path_0001/M01.txt",
		"ds1/my_path_0002/M01.txt",
		"ds1/my_path_0003/M01.txt",
		"ds1/my_path_0004/M01.txt",
		"ds1/my_path_0005/M01.txt",
		"ds1/my_path_0006/M01.txt",
		"ds1/my_path_0007/M01.txt",
		"ds1/my_path_0008/M01.txt",
		"ds1/my_path_0009/M01.txt",
		"ds1/my_path_0010/M01.txt",
	}
	for _, name := range names {
		entity := dataset1.parseEntity(name)
		fmt.Println("entity:", entity)
		dataset1.addEntity(entity)
		if dataset1.checkNewGroupAvailable(entity) {
			dataset1.outputNewGroups(entity)
		}
	}
}

func TestGetVerticalGroupRange_Interleaved(t *testing.T) {
	// y := 5
	for y := 1; y <= 10; y++ {
		arr := dataset2.getVerticalGroupRange(y)
		fmt.Println(y, arr)
	}
}
