package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Package struct {
	Dependencies         map[string]string `json:"dependencies"`
	PeerDependencies     map[string]string `json:"peerDependencies"`
	OptionalDependencies map[string]string `json:"optionalDependencies"`
}

func getDeps(packageName string) (map[string]bool, error) {
	deps := make(map[string]bool)
	var packageData Package
	r, err := http.Get("https://registry.npmjs.com/" + packageName + "/latest")
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(r.Body).Decode(&packageData)
	if err != nil {
		return nil, err
	}

	for dep := range packageData.Dependencies {
		deps[dep] = true
	}
	for dep := range packageData.PeerDependencies {
		deps[dep] = true
	}
	for dep := range packageData.OptionalDependencies {
		deps[dep] = true
	}

	return deps, nil
}

func main() {
	packageName := os.Args[1]
	depSet, err := getDeps(packageName)
	if err != nil {
		log.Fatal(err)
	}

	for dep := range depSet {
		deps, err := getDeps(dep)
		if err != nil {
			log.Fatal(err)
		}
		for d := range deps {
			depSet[d] = true
		}
	}

	depsCount := len(depSet)
	plural := "dependencies"
	if depsCount == 1 {
		plural = "dependency"
	}

	fmt.Printf("The \"%s\" package has %d %s.\n", packageName, depsCount, plural)
}
