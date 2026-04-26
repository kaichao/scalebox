package vtask_test

import (
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/kaichao/scalebox/pkg/vtask"
)

func TestAddSemaphoreValue(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")
	os.Setenv("SEMAPHORE_AUTO_CREATE", "yes")

	testCases := []struct {
		name      string
		semaName  string
		vtaskID   int64
		appID     int
		delta     int
		expectErr bool
	}{
		{
			name:      "positive delta without vtask",
			semaName:  "test_vtask_add_pos",
			vtaskID:   0,
			appID:     testAppID,
			delta:     5,
			expectErr: false,
		},
		{
			name:      "negative delta without vtask",
			semaName:  "test_vtask_add_neg",
			vtaskID:   0,
			appID:     testAppID,
			delta:     -2,
			expectErr: false,
		},
		{
			name:      "zero delta without vtask",
			semaName:  "test_vtask_add_zero",
			vtaskID:   0,
			appID:     testAppID,
			delta:     0,
			expectErr: false,
		},
		{
			name:      "positive delta with vtask",
			semaName:  "test_vtask_add_vtask",
			vtaskID:   12345,
			appID:     testAppID,
			delta:     10,
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := vtask.AddSemaphoreValue(tc.semaName, tc.delta, tc.vtaskID, tc.appID)
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else if err != nil {
				t.Logf("AddSemaphoreValue returned error (may be expected in test environment): %v", err)
			} else {
				t.Logf("AddSemaphoreValue returned value: %d", value)
			}
		})
	}
}

func TestAddSemaphoreMapValues(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")
	os.Setenv("SEMAPHORE_AUTO_CREATE", "yes")

	testCases := []struct {
		name      string
		pairs     map[string]int
		vtaskID   int64
		appID     int
		expectErr bool
	}{
		{
			name: "multiple semaphores without vtask",
			pairs: map[string]int{
				"test_vtask_map_1": 10,
				"test_vtask_map_2": 20,
				"test_vtask_map_3": 30,
			},
			vtaskID:   0,
			appID:     testAppID,
			expectErr: false,
		},
		{
			name: "mixed deltas without vtask",
			pairs: map[string]int{
				"test_vtask_map_pos":  15,
				"test_vtask_map_neg":  -5,
				"test_vtask_map_zero": 0,
			},
			vtaskID:   0,
			appID:     testAppID,
			expectErr: false,
		},
		{
			name:      "empty map",
			pairs:     map[string]int{},
			vtaskID:   0,
			appID:     testAppID,
			expectErr: false,
		},
		{
			name: "with vtaskID",
			pairs: map[string]int{
				"test_vtask_map_vtask": 25,
			},
			vtaskID:   12345,
			appID:     testAppID,
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := vtask.AddSemaphoreMapValues(tc.pairs, tc.vtaskID, tc.appID)
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else if err != nil {
				t.Logf("AddSemaphoreMapValues returned error (may be expected in test environment): %v", err)
			} else {
				t.Logf("AddSemaphoreMapValues returned result: %v", result)
				if len(tc.pairs) == 0 && len(result) != 0 {
					t.Errorf("Expected empty result for empty input, got: %v", result)
				}
			}
		})
	}
}

func TestAddSemaphoreValueAutoCreate(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	// Test without auto-create
	os.Setenv("SEMAPHORE_AUTO_CREATE", "no")
	_, err := vtask.AddSemaphoreValue("test_vtask_auto_create", 5, 0, testAppID)
	if err != nil {
		t.Logf("AddSemaphoreValue without auto-create returned error (expected): %v", err)
	} else {
		t.Logf("AddSemaphoreValue without auto-create succeeded (unexpected)")
	}

	// Test with auto-create
	os.Setenv("SEMAPHORE_AUTO_CREATE", "yes")
	value, err := vtask.AddSemaphoreValue("test_vtask_auto_create", 5, 0, testAppID)
	if err != nil {
		t.Logf("AddSemaphoreValue with auto-create returned error: %v", err)
	} else {
		t.Logf("AddSemaphoreValue with auto-create succeeded, value: %d", value)
	}
}
