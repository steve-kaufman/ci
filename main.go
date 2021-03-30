package main

import (
	"fmt"
	"net/http"
	"os/exec"
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/eshop-frontend":
		deployEshopFrontend(w)
	}
}

func deployEshopFrontend(w http.ResponseWriter) {
	err := exec.Command("docker", "pull", "stevekaufman/eshop-frontend").Run()
	if err != nil {
		sendError(w, "Failed to deploy eshop frontend")
	}

	err = exec.Command("docker", "stop", "eshop-frontend").Run()
	if err != nil {
		return
	}

	err = exec.Command("docker", "run", "-d", "--rm", "--name", "eshop-frontend", "-p", "8085:8080", "stevekaufman/eshop-frontend").Run()
	if err != nil {
		sendError(w, err.Error())
	}
}

func sendError(w http.ResponseWriter, msg string) {
	w.WriteHeader(400)
	fmt.Fprint(w, msg)
}
