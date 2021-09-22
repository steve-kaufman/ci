package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	config := loadConfig()
	fmt.Printf("%+v\n", config)

	fmt.Println("CI Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", &Handler{
		config:   config,
		deployer: NewDockerDeployer(),
	}))
}
