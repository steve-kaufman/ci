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

	var err error

	switch r.URL.Path {
	case "/eshop-frontend":
		err = deployDockerImage("eshop-frontend", "stevekaufman/eshop-frontend", 8085)
	case "/website":
		err = deployDockerImage("website", "stevekaufman/website", 5000)
	default:
		fmt.Println("--- Nothing to deploy ---")
		w.WriteHeader(404)
		fmt.Fprintf(w, "Nothing to deploy at '%s'", r.URL.Path)
		return
	}
	if err != nil {
		sendError(w, err.Error())
		return
	}

	fmt.Fprint(w, "Deployed successfully!")
}

func deployDockerImage(name string, imageName string, port int) error {
	fmt.Println("Deploying:", imageName)

	fmt.Println("Pulling new image")
	err := exec.Command("docker", "pull", imageName).Run()
	if err != nil {
		return fmt.Errorf("couldn't pull '%s'", imageName)
	}
	fmt.Println("Pulled new image")

	fmt.Println("Stopping old container")
	err = exec.Command("docker", "stop", name).Run()
	if err != nil {
		fmt.Println("No container already running")
	}
	fmt.Println("Stopped old container")

	fmt.Println("Removing old container just in case")
	err = exec.Command("docker", "rm", name).Run()
	if err != nil {
		fmt.Println("Container did not need removing")
	} else {
		fmt.Println("Container removed")
	}

	portMap := fmt.Sprintf("%d:8080", port)

	fmt.Println("Starting new container")
	err = exec.Command("docker", "run", "-d", "--rm", "--name", name, "-p", portMap, imageName).Run()
	if err != nil {
		return err
	}
	fmt.Println("Started new container")

	return nil
}

func sendError(w http.ResponseWriter, msg string) {
	fmt.Println(msg)
	w.WriteHeader(500)
	fmt.Fprint(w, msg)
}
