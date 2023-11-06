package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestGetJumpServer(test *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"user@1.2.3.4", []string{"user@1.2.3.4"}},
		{"[user@1.2.3.4,1/user@1.2.3.5,3]", []string{"user@1.2.3.4", "user@1.2.3.5"}},
	}

	for i := 0; i < 1; i++ {
		for _, t := range tests {
			got := getJumpServer(t.input)
			fmt.Printf("%s\t%s\n", t.input, got)
			if !sliceContains(t.want, got) {
				test.Errorf("getJumpServer('%s'), got = '%s'; want '%s'\n", t.input, got, t.want)
			}
		}
	}
}

func sliceContains(slice []string, elem string) bool {
	return strings.Contains(strings.Join(slice, ","), elem)
}

func TestSplit(test *testing.T) {
	cases := []struct {
		input string
		want  []string
	}{
		{"", []string{""}},
		{"1.2.3.4,", []string{"1.2.3.4", ""}},
		{",1.2.3.4", []string{"", "1.2.3.4"}},
		{"1.2.3.4", []string{"1.2.3.4"}},
		{"user@1.2.3.4", []string{"user@1.2.3.4"}},
		{"user@1.2.3.4,user@1.2.3.5", []string{"user@1.2.3.4", "user@1.2.3.5"}},
		{"[user@1.2.3.4,1/user@1.2.3.5,2]", []string{"[user@1.2.3.4,1/user@1.2.3.5,2]"}},
		{"[user@1.2.3.4,1/user@1.2.3.5,2],", []string{"[user@1.2.3.4,1/user@1.2.3.5,2]", ""}},

		{",[user@1.2.3.4,1/user@1.2.3.5,2]", []string{"", "[user@1.2.3.4,1/user@1.2.3.5,2]"}},
		{"user@1.2.3.4,[user@1.2.3.4,1/user@1.2.3.5,2]", []string{"user@1.2.3.4", "[user@1.2.3.4,1/user@1.2.3.5,2]"}},
		{"[user@1.2.3.4,1/user@1.2.3.5,2],user@1.2.3.4", []string{"[user@1.2.3.4,1/user@1.2.3.5,2]", "user@1.2.3.4"}},
	}

	for _, t := range cases {
		if got := split(t.input); !reflect.DeepEqual(got, t.want) {
			test.Errorf("split(%s), got = %v; want %v", t.input, got, t.want)
		}
	}
}

func TestCheckFormat(test *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"user@1.2.3.4", true},
		{"user@1.2.3.4,", false},
		{",user@1.2.3.4", false},
		{"1.2.3.4", false},

		{"user@1.2.3.4,user@1.2.3.5", true},
		{"user@1.2.3.4,1.2.3.5", false},
		{"1.2.3.4,user@1.2.3.5", false},

		{"[user@1.2.3.4,1/user@1.2.3.5,2]", true},
		{"user@1.2.3.4,[user@1.2.3.4,1/user@1.2.3.5,2]", true},
		{"1.2.3.4,[user@1.2.3.4,1/user@1.2.3.5,2]", false},
		{"[user@1.2.3.4,1/user@1.2.3.5,2],1.2.3.4", false},
		{"[user@1.2.3.4,1/1.2.3.5,2]", false},
		{"[user@1.2.3.4,1/user@1.2.3.5]", false},
	}

	for _, t := range tests {
		got := checkFormat(t.input)

		if got != t.want {
			test.Errorf("checkFormat('%s'), got = %v; want %v", t.input, got, t.want)
		}
	}
}
