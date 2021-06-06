package agent

import (
	"errors"
	"time"
	pb "github.com/bstr9/gpu_manager/proto"
	util "github.com/bstr9/gpu_manager/util"
)

type Agent struct {
	docker 		*DockerAgent
	addr        string
	grpcClient 	*pb.GpuApiClient
	lock 		*util.RWMutex
}

func NewAgent(addr string) (*Agent, error) {
	docker, err := NewDockerAgent()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	conn, err := grpc.Dial(addr, grpc.WithInsecure)
	return &Agent{
		docker: docker,
		grpcClient: pb.NewGpuApiClient(conn),
		lock: util.NewRWMutex("agent", 10*time.Millisecond),
	}, nil
}

func (a *Agent) Run() error {
	if a == nil {
		return errors.New("agent is not initilized")
	}
	a.docker.Run()
	return nil
}

func (a *Agent) Report() error {

}
