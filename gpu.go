package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

func main() {
	// Check if the socket exists
	socketPath := "/var/lib/kubelet/device-plugins/nvidia-gpu.sock"
	if _, err := os.Stat(socketPath); os.IsNotExist(err) {
		log.Fatalf("Socket not found at %s", socketPath)
	}

	// Create a gRPC connection
	conn, err := grpc.Dial(
		socketPath,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a device plugin client
	client := pluginapi.NewDevicePluginClient(conn)

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine to handle shutdown
	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		cancel()
	}()

	// List available devices
	resp, err := client.ListAndWatch(ctx, &pluginapi.Empty{})
	if err != nil {
		log.Fatalf("Failed to list devices: %v", err)
	}

	fmt.Println("Connected to NVIDIA device plugin. Press Ctrl+C to exit.")
	fmt.Println("Waiting for device updates...\n")

	// Print device information
	for {
		response, err := resp.Recv()
		if err != nil {
			if ctx.Err() == context.Canceled {
				fmt.Println("Gracefully shutting down...")
				return
			}
			log.Printf("Error receiving device list: %v", err)
			break
		}

		// Group GPUs by NUMA node
		numaNodes := make(map[int][]string)
		for _, device := range response.Devices {
			if device.Topology != nil && len(device.Topology.Nodes) > 0 {
				numaID := int(device.Topology.Nodes[0].ID)
				numaNodes[numaID] = append(numaNodes[numaID], device.ID)
			}
		}

		fmt.Printf("\nGPU Status Update (%s):\n", time.Now().Format("15:04:05"))
		fmt.Printf("Total GPUs: %d\n", len(response.Devices))
		fmt.Printf("Healthy GPUs: %d\n", len(response.Devices))

		fmt.Println("\nGPUs by NUMA Node:")
		for numaID, gpus := range numaNodes {
			fmt.Printf("NUMA Node %d: %d GPUs\n", numaID, len(gpus))
			for _, gpu := range gpus {
				fmt.Printf("  - %s\n", gpu)
			}
		}
		fmt.Println("----------------------------------------")
	}
}
