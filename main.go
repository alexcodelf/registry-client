package main

import (
	"fmt"

	"github.com/alexcodelf/registry-client/registry"
)

func main() {
	service, err := registry.New("https://registry-1.docker.io", "", "")
	if err != nil {
		panic(err)
	}

	manifest, err := service.ManifestDigest("library/alpine", "latest")
	if err != nil {
		panic(err)
	}

	fmt.Println(manifest)
}
