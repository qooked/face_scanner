package scannerErrors

import "fmt"

var (
	ErrDuplicateTask      = fmt.Errorf("duplicate task")
	ErrTaskAlreadyStarted = fmt.Errorf("task already started")
	ErrTaskNotFound       = fmt.Errorf("task not found")
)
