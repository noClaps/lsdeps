package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Package struct {
	Dependencies map[string]string `json:"dependencies"`
}

func isInArray(arr []string, val string) bool {
	for i := range arr {
		if arr[i] == val {
			return true
		}
	}

	return false
}

func getDeps(packageName string) ([]string, error) {
	var deps []string
	var packageData Package

	r, err0 := http.Get("https://registry.npmjs.com/" + packageName + "/latest")
	if err0 != nil {
		return nil, err0
	}

	err1 := json.NewDecoder(r.Body).Decode(&packageData)
	if err1 != nil {
		return nil, err1
	}

	for dep := range packageData.Dependencies {
		if !isInArray(deps, dep) {
			deps = append(deps, dep)
		}
	}

	return deps, nil
}

func main() {
	packageName := os.Args[1]

	fmt.Print("Counting dependencies...")

	pkgDeps, err0 := getDeps(packageName)
	if err0 != nil {
		log.Println(err0)
	}

	for d := 0; d < len(pkgDeps); d++ {
		deps, err1 := getDeps(pkgDeps[d])
		if err1 != nil {
			log.Println(err1)
		}

		for i := range deps {
			if !isInArray(pkgDeps, deps[i]) {
				pkgDeps = append(pkgDeps, deps[i])
			}
		}
	}

	depsCount := len(pkgDeps)

	fmt.Printf("\033[2K\r")

	plural := "dependencies"
	if depsCount == 1 {
		plural = "dependency"
	}

	fmt.Printf("The \"%s\" package has %d %s.\n", packageName, depsCount, plural)
}
