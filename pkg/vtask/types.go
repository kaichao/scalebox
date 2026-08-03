// Package vtask provides Go wrappers for the vtask management RPCs
// exposed by the scalebox controld AppService.
//
// Layer 2 — vtask-scoped resource operations:
//   - GetVariable / SetVariable
//   - GetSemaphore / CreateSemaphore / AddSemaphoreValue / DeleteSemaphore
//
// Layer 3 — vtask advanced operations:
//   - BindResource / UnbindResource
//   - AddSubtask / GetInfo / Fail / List / ListSubtasks
package vtask

import "time"

// Info holds the result of GetVtask (derived from GetVtask RPC).
type Info struct {
	ID            int64
	AppID         int32
	Body          string
	Status        string
	ModuleName    string
	BoundResource string
	SemaName      string
	SubtaskCount  int32
	FinishedCount int32
	CreatedAt     time.Time
}

// SubtaskItem is a simplified view of a vtask subtask.
type SubtaskItem struct {
	ID         int64
	ModuleName string
	FromModule string
	Status     string
	Body       string
}
