package main 

import (
	"context"
	dc "github.com/docker/docker/client"
	dctypes "github.com/docker/docker/api/types"
	dcevents "github.com/docker/docker/api/types/events"
)

type ContainerStatus uint8 

const (
	ContainerStatusUnknown 	ContainerStatus = 0
	ContainerStatusCreated 					= 1
	ContainerStatusRunning 					= 2
	ContainerStatusStopping 				= 3     
	ContainerStatusStopped                  = 4
	ContainerStatusDie     					= 5 
	ContainerStatusExited  					= 6 
)


type Container struct {
	Task 	    *Task
	ContainerID string
	Status 		ContainerStatus
}

func (c *Container) OnMessage(message dcevents.Message) error {
}

func (c *Container) Start() error {
}

func (c *Container) Stop() error {
}

func (c *Container) Kill() error {
}

// container report events:
// attach, commit, copy, create, destroy, 
// detach, die, exec_create, exec_detach, 
// exec_start, exec_die, export, 
// health_status, kill, oom, pause, rename, 
// resize, restart, start, stop, top, unpause, update, and prune
type DockerAgent struct {
	client *dc.Client
	containers []types.Container
}

func NewDockerAgent() (*DockerAgent, error) {
	client, err := dc.NewEnvironmentClient()
	if err != nil {
		return nil, err
	}

	return &DockerAgent{client: client}
}

func (agent *DockerAgent) Container(containerID string) {}

func (agent *DockerAgent) Containers(options types.ContainerListOptions) ([]types.Container, error) {
	return agent.client.ContainerList(context.Background(), options)
}

func (agent *DockerAgent) Stop(containerID string) {}

func (agent *DockerAgent) Kill(containerID string) {}

func (agent *DockerAgent) Clear(containerID string) {}
