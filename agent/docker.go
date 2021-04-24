package agent

import (
	"context"
	"fmt"
	util "gpu_manager/util"
	"time"

	dctypes "github.com/docker/docker/api/types"
	dcfilters "github.com/docker/docker/api/types/filters"
	dc "github.com/docker/docker/client"
)

type DockerStatus uint8

const (
	DockerStatusUnknown  DockerStatus = 0
	DockerStatusInitail  DockerStatus = 1
	DockerStatusRunning  DockerStatus = 2
	DockerStatusStopped  DockerStatus = 3
	DockerStatusMayCrash DockerStatus = 4 // <- not sure what's happened
)

type DockerAgent struct {
	client             *dc.Client
	msgC               chan DockerEvent
	ticker             *time.Ticker
	ctx                context.Context
	cancel             context.CancelFunc
	lock               *util.RWMutex
	containers         map[string]*Container
	status             DockerStatus
	lastValidTimestamp time.Time
}

func NewDockerAgent() (*DockerAgent, error) {
	client, err := dc.NewEnvironmentClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &DockerAgent{
		client:             client,
		msgC:               make(chan DockerEvent, 10),
		ticker:             time.NewTicker(1 * time.Second),
		ctx:                ctx,
		cancel:             cancel,
		lock:               util.NewMutex("agentLock", 100*time.Millisecond),
		containers:         make(map[string]*Container),
		status:             DockerStatusInitail,
		lastValidTimestamp: time.Now(),
	}
}

func (a *DockerAgent) Run() error {
	if a.client == nil {
		return Error("docker agent not initilized")
	}
	filters := dcfilters.NewArgs()
	filters.Add("type", "container")
	filters.Add("label", "gpu=true")
	options := dctypes.EventsOptions{
		Since:   fmt.Sprintf("%d", a.lastValidTimestamp),
		Filters: filters,
	}

	events, errors := a.client.Events(a.ctx, options)
	for {
		select {
		case msg := <-events:
			id := msg.Actor.ID
			a.lock.Lock()
			switch msg.Action {
			case "create":
				a.containers[id] = NewContainer(a) // todo
			case "die", "kill":
				if container, err := a.Container(id); err == nil {
					go container.OnMessage(msg)
				}
			default:
				if container, ok := a.Container(id); err == nil {
					go container.OnMessage(msg)
				} else {
					a.containers[id] = NewContainer(a)
					go a.containers[id].OnMessage(msg)
				}
			}
			a.lastValidTimestamp = msg.Time
			a.lock.Unlock()
		case err := <-errors:
			// Restart watch call
			logp.Err("Error watching for docker events: %v", err)
			time.Sleep(1 * time.Second)
			events, errors = a.client.Events(a.ctx, options)
			break WATCH
		case event := <-a.msgC:
			a.OnMessage(event)
		case <-a.ctx.Done():
			logp.Debug("docker", "Watcher stopped")
			return
		}
	}
}

type DockerEventAction uint8

const (
	EventActionUnknown          DockerEventAction = 0
	EventActionUpdateContainers                   = 1
)

type DockerEventType uint8

const (
	DockerEventTypeUnknown   DockerEventType = 0
	DockerEventTypeContainer                 = 1
)

type DockerEvent interface {
	Type() DockerEventType
	Action() DockerEventAction
}

func (a *DockerAgent) RecvMessage(event DockerAgent) {
	a.msgC <- event
}

func (a *DockerAgent) OnMessage(event DockerEvent) {
}

func (a *DockerAgent) Container(containerID string) {}

func (a *DockerAgent) Containers(options types.ContainerListOptions) ([]types.Container, error) {
	return a.client.ContainerList(context.Background(), options)
}

func (a *DockerAgent) Stop(containerID string) {}

func (a *DockerAgent) Kill(containerID string) {}

func (a *DockerAgent) Clear(containerID string) {}
