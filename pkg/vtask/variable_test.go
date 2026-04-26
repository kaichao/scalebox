package vtask_test

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/kaichao/scalebox/pkg/vtask"
)

func TestSetAndGetVariable(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	varName := "test_vtask_var"
	varValue := "hello_vtask"

	// Set variable without vtask
	err := vtask.SetVariable(varName, varValue, 0, testAppID)
	if err != nil {
		fmt.Printf("SetVariable error: %v\n", err)
	} else {
		fmt.Println("SetVariable succeeded")
	}

	// Get variable without vtask
	got, err := vtask.GetVariable(varName, 0, testAppID)
	if err != nil {
		fmt.Printf("GetVariable error: %v\n", err)
	} else if got != varValue {
		fmt.Printf("GetVariable: expected %s, got %s\n", varValue, got)
	} else {
		fmt.Printf("GetVariable succeeded, value=%s\n", got)
	}

	// Set variable with vtaskID
	varName2 := "test_vtask_var_vtask"
	err2 := vtask.SetVariable(varName2, "vtask_value", 12345, testAppID)
	if err2 != nil {
		fmt.Printf("SetVariable with vtask error: %v\n", err2)
	} else {
		fmt.Println("SetVariable with vtask succeeded")
	}

	// Get variable with vtaskID
	got2, err2 := vtask.GetVariable(varName2, 12345, testAppID)
	if err2 != nil {
		fmt.Printf("GetVariable with vtask error: %v\n", err2)
	} else {
		fmt.Printf("GetVariable with vtask succeeded, value=%s\n", got2)
	}

	// Update existing variable
	err3 := vtask.SetVariable(varName, "updated_value", 0, testAppID)
	if err3 != nil {
		fmt.Printf("SetVariable update error: %v\n", err3)
	} else {
		fmt.Println("SetVariable update succeeded")
	}

	// Verify update
	got3, err3 := vtask.GetVariable(varName, 0, testAppID)
	if err3 != nil {
		fmt.Printf("GetVariable after update error: %v\n", err3)
	} else if got3 != "updated_value" {
		fmt.Printf("GetVariable after update: expected updated_value, got %s\n", got3)
	} else {
		fmt.Printf("GetVariable after update succeeded, value=%s\n", got3)
	}
}

func TestGetNonExistentVariable(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	// Get non-existent variable should error
	_, err := vtask.GetVariable("non_existent_var_xyz", 0, testAppID)
	if err != nil {
		fmt.Printf("GetVariable non-existent error (expected): %v\n", err)
	} else {
		fmt.Println("GetVariable non-existent succeeded (unexpected)")
	}
}
