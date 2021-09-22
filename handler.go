package main

import (
	"fmt"
	"log"
	"net/http"
)

type Deployer interface {
	Deploy(container Container) error
}

type Handler struct {
	deployer Deployer
	config   Config
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received a request to deploy '%s'\n", r.URL.Path)

	if r.Header.Get("ci-access-key") != h.config.AccessKey {
		h.sendErr(w, 403, "Invalid Access Key")
		return
	}

	imageName := r.URL.Query().Get("image")
	image, ok := h.config.GetContainer(imageName)

	if !ok {
		h.sendErr(w, 400, fmt.Sprintf("Image '%s' not allowed", imageName))
		return
	}

	err := h.deployer.Deploy(image)
	if err != nil {
		h.sendErr(w, 400, "Failed to deploy image: "+err.Error())
		return
	}

	fmt.Fprint(w, "Deployed successfully!")
}

func (h Handler) sendErr(w http.ResponseWriter, errCode int, msg string) {
	fmt.Println(msg)
	w.WriteHeader(errCode)
	fmt.Fprint(w, msg)
}
