package agent

import (
	"errors"
)

type Agent struct {
	docker *DockerAgent
}

func NewAgent() (*Agent, error) {
	docker, err := NewDockerAgent()
	if err != nil {
		return nil, err
	}
	return &Agent{
		docker: docker,
	}, nil
}

func (a *Agent) Run() error {
	if a == nil {
		return errors.New("agent is not initilized")
	}
	a.docker.Run()
	return nil
}
