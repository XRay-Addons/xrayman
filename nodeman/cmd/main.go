package main

import (
	"context"
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/nodeclient"
	"github.com/XRay-Addons/xrayman/shared/models"
)

func main() {
	endpoint := "http://stepka.co.uk:8088"

	client, err := nodeclient.New(endpoint, 10)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	users := []models.User{
		models.User{
			Name: "UserA",
			UUID: "UserA-UUID",
		},
		models.User{
			Name: "UserB",
			UUID: "UserB-UUID",
		},
	}

	cfg, err := client.Start(context.Background(), users)
	if err != nil {
		fmt.Printf("start node error: %v", err)
		return
	}
	fmt.Printf("Node started successfully: %+v\n", cfg)
	status, err := client.Status(context.Background())
	if err != nil {
		fmt.Printf("get node status error: %v\n", err)
		return
	}
	fmt.Printf("Node status: %+v\n", status)
}
