package agent

import (
	"context"
	"errors"
	"fmt"
	"time"

	util "github.com/bstr9/gpu_manager/util"

	log "github.com/sirupsen/logrus"

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
	tasks              map[string]*Task
	status             DockerStatus
	lastValidTimestamp int64
}

func NewDockerAgent() (*DockerAgent, error) {
	client, err := dc.NewClientWithOpts(dc.FromEnv)
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
		lock:               util.NewRWMutex("agentLock", 100*time.Millisecond),
		containers:         make(map[string]*Container),
		tasks:              make(map[string]*Task),
		status:             DockerStatusInitail,
		lastValidTimestamp: time.Now().Unix(),
	}, nil
}

func (a *DockerAgent) Run() error {
	if a.client == nil {
		return errors.New("docker agent not initilized")
	}
	filters := dcfilters.NewArgs()
	filters.Add("type", "container")
	// filters.Add("label", "gpu=true")
	options := dctypes.EventsOptions{
		Since:   fmt.Sprintf("%d", a.lastValidTimestamp),
		Filters: filters,
	}

	events, errors := a.client.Events(a.ctx, options)
	for {
		select {
		case msg := <-events:
			log.WithFields(log.Fields{
				"message": msg,
			}).Debug("watching docker events")
			id := msg.Actor.ID
			a.lock.Lock()
			switch msg.Action {
			case "create":
				a.containers[id] = NewContainer(a, nil, id) // todo
			case "die", "kill":
				if container, ok := a.containers[id]; ok {
					go container.OnMessage(msg)
				}
			default:
				if container, ok := a.containers[id]; ok {
					go container.OnMessage(msg)
				} else {
					a.containers[id] = NewContainer(a, nil, id)
					go a.containers[id].OnMessage(msg)
				}
			}
			a.lastValidTimestamp = msg.Time
			a.lock.Unlock()
		case err := <-errors:
			// Restart watch call
			log.WithFields(log.Fields{
				"error": err,
			}).Error("watching docker events")
			time.Sleep(1 * time.Second)
			events, errors = a.client.Events(a.ctx, options)
		case event := <-a.msgC:
			a.OnMessage(event)
		case <-a.ctx.Done():
			log.Debug("watching docker stopped")
			return nil
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
	Type() 		DockerEventType
	Action() 	DockerEventAction
}

func (a *DockerAgent) OnMessage(event DockerEvent) {
	a.lock.Lock()
	defer a.lock.Unlock()
	switch event.Type() {
		case DockerEventTypeContainer:
			msg := event.(ContainerEvent)
			switch msg.Status() {
			case ContainerStatusDead, ContainerStatusRemoved:
				id := msg.Id()
				if _, ok := a.containers[id]; ok {
					delete(a.containers, id)
				}
			}
	}
	switch event.Action() {
	case EventActionUpdateContainers:
		return
	}
	return
}

func (a *DockerAgent) Container(containerID string) {}

func (a *DockerAgent) Containers(options dctypes.ContainerListOptions) ([]dctypes.Container, error) {
	return a.client.ContainerList(context.Background(), options)
}

func (a *DockerAgent) Stop(containerID string) {}

func (a *DockerAgent) Kill(containerID string) {}

func (a *DockerAgent) Clear(containerID string) {}
