package agent

import util "gpu_manager/util"

type ContainerStatus uint8

// "created","running","paused","restarting","removing","exited","dead"
const (
	ContainerStatusUnknown   ContainerStatus = 0
	ContainerStatusCreated                   = 1
	ContainerStatusRunning                   = 2
	ContainerStatusPaused                    = 3
	ContainerStatusRetarting                 = 4
	ContainerStatusRemoving                  = 5
	ContainerStatusExited                    = 6
	ContainerStatusDead                      = 7
)

type Container struct {
	lock   *util.RWMutex
	agent  *DockerAgent
	task   *Task
	id     string
	status ContainerStatus
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

// OnMessage: get event from docker agent, events are:
// Actions: attach, commit, copy, create, destroy, detach, die, exec_create, exec_detach, exec_start, exec_die, export, health_status, kill, oom, pause, rename, resize, restart, start, stop, top, unpause, update, and prune
func (c *Container) OnMessage(message dcevents.Message) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	switch msg.Action {
	case "create":
		c.status = ContainerStatusCreated
	case "start":
		c.status = ContainerStatusRunning
	case "pause":
		c.status = ContainerStatusPaused
	case "restart":
		c.status = ContainerStatusRetarting
	case "destroy":
		c.status = ContainerStatusRemoving
	case "stop":
		c.status = ContainerStatusExited
	case "die":
		c.status = ContainerStatusDead
	}
}

func (c *Container) Start() error {
}

func (c *Container) Stop() error {
}

func (c *Container) Kill() error {
}

type ContainerEvent struct {
}
