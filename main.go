package main

import (
	"github.com/bstr9/gpu_manager/agent"
)

var docker *agent.DockerAgent

func init() {
	var err error
	docker, err = agent.NewDockerAgent()
	if err != nil {
		panic(err)
	}
}

func main() {
	docker.Run()
}
