package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func main() {
  log.Fatal(http.ListenAndServe(":5500", http.HandlerFunc(handler)))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received a request to deploy '%s'", r.URL.Path)

	switch r.URL.Path {
	case "/eshop-frontend":
		deployEshopFrontend(w)
	}
}

func deployEshopFrontend(w http.ResponseWriter) {
	fmt.Println("Deploying Eshop Frontend")

	fmt.Println("Pulling new image")
	err := exec.Command("docker", "pull", "stevekaufman/eshop-frontend").Run()
	if err != nil {
		sendError(w, "Failed to deploy eshop frontend")
	}
	fmt.Println("Pulled new image")

	fmt.Println("Stopping old container")
	err = exec.Command("docker", "stop", "eshop-frontend").Run()
	if err != nil {
		fmt.Println("No container already running")
	}
	fmt.Println("Stopped old container")

	fmt.Println("Starting new container")
	err = exec.Command("docker", "run", "-d", "--rm", "--name", "eshop-frontend", "-p", "8085:8080", "stevekaufman/eshop-frontend").Run()
	if err != nil {
		sendError(w, err.Error())
	}
	fmt.Println("Started new container")
}

func sendError(w http.ResponseWriter, msg string) {
	w.WriteHeader(400)
	fmt.Fprint(w, msg)
}
