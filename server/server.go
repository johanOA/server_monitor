package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/johanOA/server-monitoring/proto"
	"google.golang.org/grpc"
	"gopkg.in/gomail.v2"
)

func sendAlert(subject, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "you@example.com")
	m.SetHeader("To", "alert@example.com")
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.example.com", 587, "user", "password")

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
	}
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewMonitoringServiceClient(conn)

	for {
		req := &pb.MetricsRequest{ServerId: "server-1"}
		res, err := c.GetMetrics(context.Background(), req)
		if err != nil {
			log.Fatalf("Failed to get metrics: %v", err)
		}

		fmt.Printf("CPU Usage: %.2f%%, Memory Usage: %.2f%%, Disk Usage: %.2f%%\n",
			res.CpuUsage, res.MemoryUsage, res.DiskUsage)

		if res.CpuUsage > 80 {
			sendAlert("High CPU Usage", fmt.Sprintf("CPU Usage is %.2f%%", res.CpuUsage))
		}
		if res.MemoryUsage > 80 {
			sendAlert("High Memory Usage", fmt.Sprintf("Memory Usage is %.2f%%", res.MemoryUsage))
		}
		if res.DiskUsage > 80 {
			sendAlert("High Disk Usage", fmt.Sprintf("Disk Usage is %.2f%%", res.DiskUsage))
		}

		time.Sleep(10 * time.Second)
	}
}
