package vtask_test

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/kaichao/scalebox/pkg/vtask"
)

func TestCreateSemaphore(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	// Test creating a single semaphore without vtask
	err := vtask.CreateSemaphore("test_vtask_create_single", 100, 0, testAppID)
	if err != nil {
		fmt.Printf("CreateSemaphore single error (expected due to foreign key): %v\n", err)
	} else {
		fmt.Println("CreateSemaphore single succeeded")
	}

	// Test creating a semaphore with vtaskID
	err2 := vtask.CreateSemaphore("test_vtask_create_vtask", 200, 12345, testAppID)
	if err2 != nil {
		fmt.Printf("CreateSemaphore with vtask error: %v\n", err2)
	} else {
		fmt.Println("CreateSemaphore with vtask succeeded")
	}

	// Test CONFLICT_ACTION=IGNORE
	os.Setenv("CONFLICT_ACTION", "IGNORE")
	err3 := vtask.CreateSemaphore("test_vtask_create_single", 150, 0, testAppID)
	if err3 != nil {
		fmt.Printf("CreateSemaphore with IGNORE error: %v\n", err3)
	} else {
		fmt.Println("CreateSemaphore with IGNORE succeeded")
	}

	// Test CONFLICT_ACTION=OVERWRITE
	os.Setenv("CONFLICT_ACTION", "OVERWRITE")
	err4 := vtask.CreateSemaphore("test_vtask_create_single", 150, 0, testAppID)
	if err4 != nil {
		fmt.Printf("CreateSemaphore with OVERWRITE error: %v\n", err4)
	} else {
		fmt.Println("CreateSemaphore with OVERWRITE succeeded")
	}
}

func TestCreateSemaphores(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	lines := []string{
		`"test_vtask_sema_1":10`,
		`"test_vtask_sema_2":20`,
		`"test_vtask_sema_3":30`,
	}

	// Test batch creating semaphores without vtask
	err := vtask.CreateSemaphores(lines, 0, testAppID, 10)
	if err != nil {
		fmt.Printf("CreateSemaphores error (expected due to foreign key): %v\n", err)
	} else {
		fmt.Println("CreateSemaphores succeeded")
	}

	// Test empty list
	emptyLines := []string{}
	err2 := vtask.CreateSemaphores(emptyLines, 0, testAppID, 10)
	if err2 != nil {
		fmt.Printf("CreateSemaphores empty lines error: %v\n", err2)
	} else {
		fmt.Println("CreateSemaphores empty lines succeeded")
	}

	// Test with vtaskID
	lines2 := []string{
		`"test_vtask_sema_vtask_1":100`,
		`"test_vtask_sema_vtask_2":200`,
	}
	err3 := vtask.CreateSemaphores(lines2, 12345, testAppID, 10)
	if err3 != nil {
		fmt.Printf("CreateSemaphores with vtask error: %v\n", err3)
	} else {
		fmt.Println("CreateSemaphores with vtask succeeded")
	}
}

func TestCreateSemaphoreWithConflictAction(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	semaphoreName := "test_vtask_conflict_sema"

	// Test CONFLICT_ACTION=OVERWRITE
	os.Setenv("CONFLICT_ACTION", "OVERWRITE")

	err1 := vtask.CreateSemaphore(semaphoreName, 100, 0, testAppID)
	if err1 != nil {
		fmt.Printf("Create with OVERWRITE (first time) error: %v\n", err1)
	} else {
		fmt.Println("Create with OVERWRITE (first time) succeeded")
	}

	err2 := vtask.CreateSemaphore(semaphoreName, 200, 0, testAppID)
	if err2 != nil {
		fmt.Printf("Create with OVERWRITE (second time) error: %v\n", err2)
	} else {
		fmt.Println("Create with OVERWRITE (second time) succeeded")
	}

	// Test CONFLICT_ACTION=IGNORE
	os.Setenv("CONFLICT_ACTION", "IGNORE")

	err3 := vtask.CreateSemaphore("test_vtask_conflict_ignore", 300, 0, testAppID)
	if err3 != nil {
		fmt.Printf("Create with IGNORE (first time) error: %v\n", err3)
	} else {
		fmt.Println("Create with IGNORE (first time) succeeded")
	}

	err4 := vtask.CreateSemaphore("test_vtask_conflict_ignore", 400, 0, testAppID)
	if err4 != nil {
		fmt.Printf("Create with IGNORE (second time) error: %v\n", err4)
	} else {
		fmt.Println("Create with IGNORE (second time) succeeded")
	}

	// Test CONFLICT_ACTION unset (default behavior, should error on conflict)
	os.Unsetenv("CONFLICT_ACTION")

	err5 := vtask.CreateSemaphore("test_vtask_conflict_default", 500, 0, testAppID)
	if err5 != nil {
		fmt.Printf("Create without CONFLICT_ACTION (first time) error: %v\n", err5)
	} else {
		fmt.Println("Create without CONFLICT_ACTION (first time) succeeded")
	}

	err6 := vtask.CreateSemaphore("test_vtask_conflict_default", 600, 0, testAppID)
	if err6 != nil {
		fmt.Printf("Create without CONFLICT_ACTION (second time, should fail) error: %v\n", err6)
	} else {
		fmt.Println("Create without CONFLICT_ACTION (second time) succeeded (unexpected)")
	}
}
