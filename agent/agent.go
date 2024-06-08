package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/johanOA/server-monitoring/proto"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMonitoringServiceServer
}

func (s *server) GetMetrics(ctx context.Context, req *pb.MetricsRequest) (*pb.MetricsResponse, error) {
	cpuUsage, _ := cpu.Percent(time.Second, false)
	memStats, _ := mem.VirtualMemory()
	diskStats, _ := disk.Usage("/")

	return &pb.MetricsResponse{
		CpuUsage:    cpuUsage[0],
		MemoryUsage: memStats.UsedPercent,
		DiskUsage:   diskStats.UsedPercent,
	}, nil
}

func main() {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMonitoringServiceServer(s, &server{})

	fmt.Println("Agent is running on port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
