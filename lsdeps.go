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

func contains(deps []string, packageName string) bool {
	for _, dep := range deps {
		if dep == packageName {
			return true
		}
	}
	return false
}

func getDeps(packageName string) ([]string, error) {
	fmt.Printf("\033[2KFetching dependencies for %s...", packageName)

	var deps []string
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
		deps = append(deps, dep)
	}
	for dep := range packageData.PeerDependencies {
		deps = append(deps, dep)
	}
	for dep := range packageData.OptionalDependencies {
		deps = append(deps, dep)
	}

	return deps, nil
}

func main() {
	packageName := os.Args[1]
	depSet, err := getDeps(packageName)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(depSet); i++ {
		deps, err := getDeps(depSet[i])
		if err != nil {
			log.Fatal(err)
		}
		for _, d := range deps {
			if !contains(depSet, d) {
				depSet = append(depSet, d)
			}
		}
	}

	depsCount := len(depSet)
	plural := "dependencies"
	if depsCount == 1 {
		plural = "dependency"
	}

	fmt.Printf("\033[2KThe \"%s\" package has %d %s.\n", packageName, depsCount, plural)
}
