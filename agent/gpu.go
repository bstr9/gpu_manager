package main 

import (
	"github.com/docker/docker/client"
)

type GPUStatus uint8

const (
	GPUStatusUnknown GPUStatus = 0
)

type GPU struct {
	ID 			string
	Status 		GPUStatus
	memory 		uint64
	availableMemory uint64
}

func (gpu *GPU) Apply(taskID string, memory uint8) error {
	return nil
}

func (gpu *GPU) Release(taskID string) error {
	return nil
}

func (gpu *GPU) Run() error {
	return nil
}

func (gpu *GPU) Memory() uint64 {
	return 0
}
