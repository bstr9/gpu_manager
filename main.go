package main

import (
	"github.com/bstr9/gpu_manager/agent"
)

func main() {
	agent, err := agent.NewAgent()
	if err != nil {
		panic(err)
	}
	agent.Run()
}
