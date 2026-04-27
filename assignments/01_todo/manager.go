package todo

// Manager defines the contract for task business logic.
// Consumers (HTTP handlers, CLI, gRPC) depend on this interface.
//go:generate go-mock-gen --interface=Manager
type Manager interface {
	AddTask(cmd *AddCmd) (*Task, error)
	ListTasks(cmd *ListCmd) ([]Task, error)
	UpdateTasks(cmd *UpdateCmd) (*Task, error)
	DeleteTask(cmd *DeleteCmd) error
	ToggleDone(cmd *DoneCmd) (*Task, error)
}
