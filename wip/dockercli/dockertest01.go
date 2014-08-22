package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"log"
)

func main() {
	fmt.Println("Hello, World")
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		log.Fatal("Error create client: ", err)
	}
	imgs, err := client.ListImages(true)
	if err != nil {
		log.Fatal("Error list images: ", err)
	}

	for _, img := range imgs {
		fmt.Println("ID: ", img.ID)
		fmt.Println("RepoTags: ", img.RepoTags)
		fmt.Println("Created: ", img.Created)
		fmt.Println("Size: ", img.Size)
		fmt.Println("Virtual Size: ", img.VirtualSize)
		fmt.Println("ParentId: ", img.ParentId)
		fmt.Println("Repository: ", img.Repository)
	}
}
