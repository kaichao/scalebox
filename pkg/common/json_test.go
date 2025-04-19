package common_test

import (
	"testing"

	"github.com/kaichao/gopkg/common"
)

func TestSetJSONAttribute(t *testing.T) {
	tests := []struct {
		input    string
		name     string
		value    string
		expected string
	}{
		{"", "attr", "value", `{"attr":"value"}`},
		{"{}", "attr", "value", `{"attr":"value"}`},
	}

	for _, tc := range tests {
		result := common.SetJSONAttribute(tc.input, tc.name, tc.value)
		if result != tc.expected {
			t.Errorf(`common.SetJSONAttribute("%s","%s","%s"), expected "%s", got "%v"`,
				tc.input, tc.name, tc.value, tc.expected, result)
		}
	}
}
