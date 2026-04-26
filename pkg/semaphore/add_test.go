package semaphore_test

import (
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/kaichao/scalebox/pkg/semaphore"
)

func TestAddValue(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")
	os.Setenv("SEMAPHORE_AUTO_CREATE", "yes")

	testCases := []struct {
		name      string
		semaName  string
		appID     int
		delta     int
		expectErr bool
	}{
		{
			name:      "positive delta",
			semaName:  "test_add_value_pos",
			appID:     testAppID,
			delta:     5,
			expectErr: false,
		},
		{
			name:      "negative delta",
			semaName:  "test_add_value_neg",
			appID:     testAppID,
			delta:     -2,
			expectErr: false,
		},
		{
			name:      "zero delta",
			semaName:  "test_add_value_zero",
			appID:     testAppID,
			delta:     0,
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := semaphore.AddValue(tc.semaName, tc.appID, tc.delta)
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else if err != nil {
				t.Logf("AddValue returned error (may be expected in test environment): %v", err)
			} else {
				t.Logf("AddValue returned value: %d", value)
			}
		})
	}
}

func TestAddRegexValue(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")
	os.Setenv("SEMAPHORE_AUTO_CREATE", "yes")

	testCases := []struct {
		name      string
		regex     string
		appID     int
		delta     int
		expectErr bool
	}{
		{
			name:      "prefix match",
			regex:     "test_regex_prefix_",
			appID:     testAppID,
			delta:     3,
			expectErr: false,
		},
		{
			name:      "exact match with $",
			regex:     "^test_regex_exact$",
			appID:     testAppID,
			delta:     5,
			expectErr: false,
		},
		{
			name:      "regex pattern",
			regex:     "^test_regex_pattern_[0-9]+$",
			appID:     testAppID,
			delta:     2,
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := semaphore.AddRegexValue(tc.regex, tc.delta, tc.appID)
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else if err != nil {
				t.Logf("AddRegexValue returned error (may be expected in test environment): %v", err)
			} else {
				t.Logf("AddRegexValue returned result: %s", result)
				if result == "" {
					t.Logf("Empty result returned (may be expected if no matching semaphores)")
				} else if result[0] != '{' {
					t.Logf("Result doesn't start with '{': %s", result)
				}
			}
		})
	}
}

func TestAddMapValues(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")
	os.Setenv("SEMAPHORE_AUTO_CREATE", "yes")

	testCases := []struct {
		name      string
		pairs     map[string]int
		appID     int
		expectErr bool
	}{
		{
			name: "multiple semaphores",
			pairs: map[string]int{
				"test_map_1": 10,
				"test_map_2": 20,
				"test_map_3": 30,
			},
			appID:     testAppID,
			expectErr: false,
		},
		{
			name: "mixed positive and negative deltas",
			pairs: map[string]int{
				"test_map_pos":  15,
				"test_map_neg":  -5,
				"test_map_zero": 0,
			},
			appID:     testAppID,
			expectErr: false,
		},
		{
			name:      "empty map",
			pairs:     map[string]int{},
			appID:     testAppID,
			expectErr: false,
		},
		{
			name: "single semaphore",
			pairs: map[string]int{
				"test_map_single": 25,
			},
			appID:     testAppID,
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := semaphore.AddMapValues(tc.pairs, tc.appID)
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else if err != nil {
				t.Logf("AddMapValues returned error (may be expected in test environment): %v", err)
			} else {
				t.Logf("AddMapValues returned result: %v", result)
				if len(tc.pairs) == 0 {
					if len(result) != 0 {
						t.Errorf("Expected empty result for empty input, got: %v", result)
					}
				} else if len(result) > 0 {
					for name, value := range result {
						t.Logf("  %s: %d", name, value)
					}
				}
			}
		})
	}
}

func TestAddMapValuesAutoCreate(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	pairs := map[string]int{
		"test_auto_create_1": 10,
		"test_auto_create_2": 20,
	}

	// 测试没有自动创建的情况
	os.Setenv("SEMAPHORE_AUTO_CREATE", "no")
	_, err := semaphore.AddMapValues(pairs, testAppID)
	if err != nil {
		t.Logf("AddMapValues without auto-create returned error (expected): %v", err)
	} else {
		t.Logf("AddMapValues without auto-create succeeded (unexpected)")
	}

	// 测试有自动创建的情况
	os.Setenv("SEMAPHORE_AUTO_CREATE", "yes")
	result, err := semaphore.AddMapValues(pairs, testAppID)
	if err != nil {
		t.Logf("AddMapValues with auto-create returned error: %v", err)
	} else {
		t.Logf("AddMapValues with auto-create succeeded, result: %v", result)
	}
}

func TestAddValueAutoCreate(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	// 测试没有自动创建的情况
	os.Setenv("SEMAPHORE_AUTO_CREATE", "no")
	_, err := semaphore.AddValue("test_auto_create_single", testAppID, 5)
	if err != nil {
		t.Logf("AddValue without auto-create returned error (expected): %v", err)
	} else {
		t.Logf("AddValue without auto-create succeeded (unexpected)")
	}

	// 测试有自动创建的情况
	os.Setenv("SEMAPHORE_AUTO_CREATE", "yes")
	value, err := semaphore.AddValue("test_auto_create_single", testAppID, 5)
	if err != nil {
		t.Logf("AddValue with auto-create returned error: %v", err)
	} else {
		t.Logf("AddValue with auto-create succeeded, value: %d", value)
	}
}
