package main

import (
	"context"
	pb "gpu_manager/proto"
)

type ApiService struct {}

func (*ApiService) Report(ctx context.Context, in *pb.ReportRequest) (out *pb.ReportResponse, err error) {
	meta := in.Meta
	tasks := in.Tasks
}
