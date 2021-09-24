package main

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerDeployer struct {
	cli *client.Client
}

func NewDockerDeployer() *DockerDeployer {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err.Error())
	}
	return &DockerDeployer{
		cli: dockerClient,
	}
}

func (d DockerDeployer) Deploy(container Container) error {
	fmt.Println("Deploying:", container.Name)
	err := d.pullImage(container.Config.Image)
	if err != nil {
		return err
	}

	return d.restartContainer(container)
}

func (d DockerDeployer) pullImage(imageName string) error {
	fmt.Println("Pulling new image")
	rc, err := d.cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(&bytes.Buffer{}, rc)
	fmt.Println("Pulled new image")
	return nil
}

func (d DockerDeployer) restartContainer(container Container) error {
	d.stopContainer(container.Name)
	d.removeContainer(container.Name)

	return d.startNewContainer(container)
}

func (d DockerDeployer) stopContainer(containerName string) {
	fmt.Println("Stopping old container")
	err := d.cli.ContainerStop(context.Background(), containerName, nil)
	if err != nil {
		fmt.Println("No container already running")
	}
	fmt.Println("Stopped old container")
}

func (d DockerDeployer) removeContainer(containerName string) {
	fmt.Println("Removing old container just in case")
	err := d.cli.ContainerRemove(context.Background(), containerName, types.ContainerRemoveOptions{})
	if err != nil {
		fmt.Println("Container did not need removing")
		return
	}
	fmt.Println("Container removed")
}

func (d DockerDeployer) startNewContainer(container Container) error {
	err := d.createContainer(container)
	if err != nil {
		return err
	}
	err = d.cli.ContainerStart(context.Background(), container.Name, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	fmt.Println("Started new container")

	return nil
}

func (d DockerDeployer) createContainer(container Container) error {
	fmt.Println("Starting new container")
	_, err := d.cli.ContainerCreate(
		context.Background(),
		container.Config, container.HostConfig,
		container.NetworkingConfig, container.Platform,
		container.Name,
	)
	return err
}
