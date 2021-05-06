package agent

import (
	"time"

	util "github.com/bstr9/gpu_manager/util"
	dcevents "github.com/docker/docker/api/types/events"
	log "github.com/sirupsen/logrus"
)

type ContainerStatus uint8

// "created","running","paused","restarting","exited","dead","removed"
const (
	ContainerStatusUnknown   ContainerStatus = 0
	ContainerStatusCreated                   = 1
	ContainerStatusRunning                   = 2
	ContainerStatusPaused                    = 3
	ContainerStatusRetarting                 = 4
	ContainerStatusExited                    = 5 
	ContainerStatusDead                      = 6 
	ContainerStatusRemoved                   = 255 
)

type Container struct {
	lock   *util.RWMutex
	agent  *DockerAgent
	task   *Task
	id     string
	status ContainerStatus
}

func NewContainer(agent *DockerAgent, task *Task, id string) *Container {
	return &Container{
		lock:   util.NewRWMutex("container", 100*time.Millisecond),
		agent:  agent,
		task:   task,
		id:     id,
		status: ContainerStatusCreated,
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
// if use `docker stop xxx` command will recv these messages: kill -> die -> stop 
// if use `docker pause xxx` command will recv message: pause, and if use `docker unpause xxx` command will recv messgae: unpause
func (c *Container) OnMessage(msg dcevents.Message) error {
	log.WithFields(log.Fields{
		"container": c.Id,
		"message":   msg,
	}).Info("container recv message")
	c.lock.Lock()
	defer c.lock.Unlock()
	switch msg.Action {
	case "create":
		c.status = ContainerStatusCreated
	case "start":
		c.status = ContainerStatusRunning
	case "pause":
		c.status = ContainerStatusPaused
	case "unpause":
		c.status = ContainerStatusRunning
	case "restart":
		c.status = ContainerStatusRetarting
	case "destroy":
		c.status = ContainerStatusRemoved
		event := c.NewEvent(ContainerStatusRemoved, EventActionUpdateContainers)
		c.agent.OnMessage(event)
	case "stop":
		c.status = ContainerStatusExited
	case "die":
		event := c.NewEvent(ContainerStatusDead, EventActionUpdateContainers)
		c.status = ContainerStatusDead
		c.agent.OnMessage(event)
	}
	return nil
}

func (c *Container) NewEvent(toStatus ContainerStatus, action DockerEventAction) ContainerEvent {
	event := ContainerEvent{
		id:         c.id,
		fromStatus: c.status,
		toStatus:   toStatus,
		timestamp:  time.Now(),
		action:     action,
	}
	return event
}

//func (c *Container) Start() error {
//}
//
//func (c *Container) Stop() error {
//}
//
//func (c *Container) Kill() error {
//}

type ContainerEvent struct {
	id         string
	fromStatus ContainerStatus
	toStatus   ContainerStatus
	timestamp  time.Time
	action     DockerEventAction
}

func (e ContainerEvent) Id() string {
	return e.id
}

func (e ContainerEvent) FromStatus() ContainerStatus {
	return e.fromStatus
}

func (e ContainerEvent) Status() ContainerStatus {
	return e.toStatus
}

func (e ContainerEvent) Type() DockerEventType {
	return DockerEventTypeContainer
}

func (e ContainerEvent) Action() DockerEventAction {
	return e.action
}
