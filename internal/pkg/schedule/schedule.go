package schedule

type Schedule interface {
	// AddTask adds a task to the schedule
	AddTask(task Task) error
	// RemoveTask removes a task from the schedule
	RemoveTask(task Task) error
	// GetTasks returns all tasks in the schedule
	GetTasks() []Task
	// Start starts the schedule
	Start() error
}

type Task struct {
	// Name is the name of the task
	Name string
	// Priority is the priority of the task
	Priority int
	// Func is the function to be executed
	Func func() error
	// CreateAt is the time when the task is created
	CreateAt int64
	// UpdateAt is the time when the task is updated
	UpdateAt int64
}
