package agent 

import (
	"context"
	dc "github.com/docker/docker/client"
	dctypes "github.com/docker/docker/api/types"
	dcfilters "github.com/docker/docker/api/types/filters"
	dcevents "github.com/docker/docker/api/types/events"
	util "gpu_manager/util"
)


type DockerStatus uint8

const (
	DockerStatusUnknown 	DockerStatus = 0
	DockerStatusInitail 	DockerStatus = 1
	DockerStatusRunning 	DockerStatus = 2 
	DockerStatusStopped 	DockerStatus = 3
	DockerStatusMayCrash 	DockerStatus = 4 // <- not sure what's happened
)

type DockerAgent struct {
	client 				*dc.Client
	ctx 				context.Context	
	cancel 				context.CancelFunc
	lock 				*util.RWMutex
	containers 			map[string]dctypes.Container
	status 				DockerStatus
	lastValidTimestamp 	time.Time
}

func NewDockerAgent() (*DockerAgent, error) {
	client, err := dc.NewEnvironmentClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &DockerAgent{
		client: client,
		ctx: ctx,
		cancel: cancel,
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

    for {
		events, errors := a.client.Events(a.ctx, options)

    WATCH:
        for {
            select {
			case msg := <-events:
            	fmt.Printf("get event %v action %v \n", msg.Actor.Attributes["name"], msg.Action)
				a.lastValidTimestamp = msg.Time
				id := msg.Actor.ID
				a.lock.Lock()
				if container, ok := a.containers[id]; ok {
					go container.OnMessage(msg)
				} else {
					container := &NewContainer(a) // TODO@bjw
					a.containers[id] = container
					go container.OnMessage(msg)
				}
				a.lock.Unlock()

				// create / update
                if event.Action == "create" || event.Action == "update" {
                    name := event.Actor.Attributes["name"]
                    image := event.Actor.Attributes["image"]
                    delete(event.Actor.Attributes, "name")
                    delete(event.Actor.Attributes, "image")
                    container := &Container{
                        ID:     event.Actor.ID,
                        Name:   name,
                        Image:  image,
                        Labels: event.Actor.Attributes,
                        Env:    make(map[string]string),
                    }
					info, err := a.client.ContainerInspect(w.ctx, event.Actor.ID)
                    if err == nil {
                        for _, env := range info.Config.Env {
                            kv := strings.SplitN(env, "=", 2)
                            if len(kv) >= 2 {
                                container.Env[kv[0]] = kv[1]
                            }
                        }
                    }
					a.containers[container.ID] = container
					a.containers[container.Name] = container
                }

                // Delete
                if event.Action == "die" || event.Action == "kill" {
                    delete(w.containers, event.Actor.ID)
                    delete(w.containers, event.Actor.Attributes["name"])
                }

            case err := <-errors:
                // Restart watch call
                logp.Err("Error watching for docker events: %v", err)
                time.Sleep(1 * time.Second)
                break WATCH
			case <-a.ctx.Done():
                logp.Debug("docker", "Watcher stopped")
                return
            }
        }
    }
}

func (a *DockerAgent) Container(containerID string) {}

func (a *DockerAgent) Containers(options types.ContainerListOptions) ([]types.Container, error) {
	return a.client.ContainerList(context.Background(), options)
}

func (a *DockerAgent) Stop(containerID string) {}

func (a *DockerAgent) Kill(containerID string) {}

func (a *DockerAgent) Clear(containerID string) {}
