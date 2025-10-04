package main

import (
	_ "github.com/velocityapp/vlib-logger" // Initialize logger

	"vapp/controller/routes"

	_ "github.com/velocityapp/vlib-k8s" // Initialize k8s client
)

func main() {
	// Entry point for the vapi-controller application

	routes.Start()
}
