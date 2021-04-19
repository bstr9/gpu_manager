package agent

type ContainerStatus uint8 

// container has 4 status: restarting, running, paused, exited
const (
	ContainerStatusUnknown 		ContainerStatus = 0
	ContainerStatusRetarting 					= 1
	ContainerStatusRunning 						= 2 
	ContainerStatusPaused                   	= 3 
	ContainerStatusExited  						= 4 
)


type Container struct {
	lock 		*util.RWMutex
	agent       *DockerAgent
	task 	    *Task
	id 			string
	status 		ContainerStatus
}

func NewContainer(agent *DockerAgent) *Container {
	return &Container{
		agent: agent,
	}
}

func (c *Container) Task() *Task {
	return c.task
}

func (c *Container) Status() ContainerStatus {
	return c.status
}

func (c *Container) Id() string {
	return c.id
}

func (c *Container) OnMessage(message dcevents.Message) error {
}

func (c *Container) Start() error {
}

func (c *Container) Stop() error {
}

func (c *Container) Kill() error {
}

