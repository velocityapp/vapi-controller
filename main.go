package main

import (
	_ "vapp/controller/logger" // This line is necessary for go-swagger to find docs
	"vapp/controller/routes"

	_ "github.com/velocityapp/vlib-k8s" // Initialize k8s client
)

func main() {
	// Entry point for the vapi-controller application
	routes.Start()
}
