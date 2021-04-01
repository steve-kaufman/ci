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
		deployDockerImage(w, "stevekaufman/eshop-frontend", 8085)
	case "/website":
		deployDockerImage(w, "stevekaufman/website", 5000)
	default:
		fmt.Println("--- Nothing to deploy ---")
		w.WriteHeader(404)
		fmt.Fprintf(w, "Nothing to deploy at '%s'", r.URL.Path)
	}
}

func deployDockerImage(w http.ResponseWriter, imageName string, port int) {
	fmt.Printf("Deploying '%s'", imageName)

	fmt.Println("Pulling new image")
	err := exec.Command("docker", "pull", imageName).Run()
	if err != nil {
		sendError(w, fmt.Sprintf("Failed to deploy '%s'", imageName))
		return
	}
	fmt.Println("Pulled new image")

	fmt.Println("Stopping old container")
	err = exec.Command("docker", "stop", imageName).Run()
	if err != nil {
		fmt.Println("No container already running")
	}
	fmt.Println("Stopped old container")

	fmt.Println("Removing old container just in case")
	err = exec.Command("docker", "rm", imageName).Run()
	if err != nil {
		fmt.Println("Container did not need removing")
	} else {
		fmt.Println("Container removed")
	}

	portMap := fmt.Sprintf("%d:8080", port)

	fmt.Println("Starting new container")
	err = exec.Command("docker", "run", "-d", "--rm", "--name", imageName, "-p", portMap, imageName).Run()
	if err != nil {
		sendError(w, err.Error())
		return
	}
	fmt.Println("Started new container")

	sendSuccess(w, fmt.Sprintf("Deployed '%s'", imageName))
}

func sendError(w http.ResponseWriter, msg string) {
	fmt.Println(msg)
	w.WriteHeader(400)
	fmt.Fprint(w, msg)
}

func sendSuccess(w http.ResponseWriter, msg string) {
	fmt.Fprint(w, msg)
}
