package vtask_test

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/kaichao/scalebox/pkg/vtask"
)

func TestGetSemaphore(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	// Test with SEMAPHORE_AUTO_CREATE=yes, non-existent semaphore auto-created
	os.Setenv("SEMAPHORE_AUTO_CREATE", "yes")
	value1, err1 := vtask.GetSemaphore("test_vtask_getvalue_sema", 0, testAppID)
	if err1 != nil {
		fmt.Printf("GetSemaphore with auto-create error: %v\n", err1)
	} else {
		fmt.Printf("GetSemaphore with auto-create succeeded, value=%d\n", value1)
	}

	// Test with vtaskID
	value2, err2 := vtask.GetSemaphore("test_vtask_getvalue_vtask", 12345, testAppID)
	if err2 != nil {
		fmt.Printf("GetSemaphore with vtask error: %v\n", err2)
	} else {
		fmt.Printf("GetSemaphore with vtask succeeded, value=%d\n", value2)
	}

	// Test without auto-create, non-existent semaphore should error
	os.Unsetenv("SEMAPHORE_AUTO_CREATE")
	value3, err3 := vtask.GetSemaphore("non_existent_vtask_sema", 0, testAppID)
	if err3 != nil {
		fmt.Printf("GetSemaphore without auto-create error (expected): %v\n", err3)
	} else {
		fmt.Printf("GetSemaphore without auto-create succeeded (unexpected), value=%d\n", value3)
	}
}

func TestGetSemaphoreNotFoundLogic(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	testCases := []struct {
		name        string
		autoCreate  string
		expectError bool
		expectValue int
	}{
		{"auto_create_yes", "yes", false, 0},
		{"auto_create_no", "no", true, 0},
		{"auto_create_unset", "", true, 0},
	}

	for _, tc := range testCases {
		os.Setenv("SEMAPHORE_AUTO_CREATE", tc.autoCreate)

		value, err := vtask.GetSemaphore(tc.name, 0, testAppID)

		if tc.expectError {
			if err == nil {
				fmt.Printf("Test %s: Expected error but got none, value=%d\n", tc.name, value)
			} else {
				fmt.Printf("Test %s: Got expected error: %v\n", tc.name, err)
			}
		} else {
			if err != nil {
				fmt.Printf("Test %s: Unexpected error: %v\n", tc.name, err)
			} else if value != tc.expectValue {
				fmt.Printf("Test %s: Expected value %d but got %d\n", tc.name, tc.expectValue, value)
			} else {
				fmt.Printf("Test %s: Success, value=%d\n", tc.name, value)
			}
		}
	}
}
