package main

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/johanOA/server-monitoring/proto"
	"google.golang.org/grpc"
)

type Metrics struct {
	CPU    float64
	Memory float64
	Disk   float64
}

func getMetrics() Metrics {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewMonitoringServiceClient(conn)

	req := &pb.MetricsRequest{ServerId: "server-1"}
	res, err := c.GetMetrics(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to get metrics: %v", err)
	}

	return Metrics{
		CPU:    res.CpuUsage,
		Memory: res.MemoryUsage,
		Disk:   res.DiskUsage,
	}
}

func main() {
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"getMetrics": getMetrics,
	})

	r.LoadHTMLFiles("./templates/index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.Run(":8081")
}
