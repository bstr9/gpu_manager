package agent 

type TaskStatus uint32

const (
	TaskStatusUnknown 			TaskStatus = 0
	TaskStatusInitial 			TaskStatus = 1
	TaskStatusCreated 			TaskStatus = 2
	TaskStatusRunning 			TaskStatus = 50 
	TaskStatusWaitingSuccess 	TaskStatus = 51
	TaskStatusWaitingKill 		TaskStatus = 52
	TaskStatusExited  			TaskStatus = 98
	TaskStatusSuccess 			TaskStatus = 99
	TaskStatusFinished 			TaskStatus = 100
)

type TaskError uint8

const (
	TaskErrorOk 				TaskError = 0
	TaskErrorUnknown 			TaskError = 1
	TaskErrorBadRequest 		TaskError = 2
	TaskErrorPermissionDenied 	TaskError = 3
)

type Task struct {
	ID 			string
	Name 		string
	Image 		string
	Commands    []string
	Volumes     []string
	Environment map[string]string
	ContainerID string
	Container   *Container
	Status 		TaskStatus
	CreatedTime time.Time
	UpdatedTime time.Time
	DeletedTime time.Time
}

type TaskList struct {
	ExceptedTasks []Task
	CurrentTasks  []Task
}
