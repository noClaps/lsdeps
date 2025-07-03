package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/noclaps/applause"
	"github.com/noclaps/lsdeps/internal/logger"
	"github.com/noclaps/lsdeps/internal/npm"
)

type args struct {
	Package      string `help:"The npm package to count dependencies for."`
	SkipOptional bool   `type:"option" short:"o" help:"Skip counting optional dependencies."`
	SkipPeer     bool   `type:"option" short:"p" help:"Skip counting peer dependencies"`
	Version      string `type:"option" help:"The version of the package being fetched."`
}

func main() {
	args := args{Version: "latest"}
	err := applause.Parse(&args)
	if err != nil {
		logger.Errorf("ERROR: %v\n", err)
		os.Exit(1)
	}

	packageName := args.Package
	version := args.Version
	fmt.Print("Fetching dependencies...")

	depSet := map[string]struct{}{}

	queue, err := npm.GetDeps(packageName, args.SkipPeer, args.SkipOptional, version)
	if err != nil {
		logger.Errorf("\033[2K\rERROR: Package %s@%s does not exist", packageName, version)
		return
	}

	var wg sync.WaitGroup // To stop code from continuing before all async tasks are finished
	var mu sync.Mutex     // To lock reading and writing shared values

	for len(queue) > 0 {
		mu.Lock()
		currentQueue := queue
		queue = make(map[string]string)
		mu.Unlock()

		for setPackage, setPackageVersion := range currentQueue {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mu.Lock()
				if _, exists := depSet[setPackage]; exists {
					mu.Unlock()
					return
				}
				depSet[setPackage] = struct{}{}
				mu.Unlock()

				deps, err := npm.GetDeps(setPackage, args.SkipPeer, args.SkipOptional, setPackageVersion)
				if err != nil {
					logger.Errorf("\033[2K\rERROR: Package %s@%s does not exist", setPackage, setPackageVersion)
					return
				}

				mu.Lock()
				for dep, version := range deps {
					if _, exists := depSet[dep]; !exists {
						queue[dep] = version
					}
				}
				mu.Unlock()
			}()
		}
		wg.Wait()
	}

	fmt.Printf("\033[2K\r")
	fmt.Printf(`Name: %s
URL: https://npmjs.com/package/%s/v/%s
Dependency count: %d
`, packageName, packageName, version, len(depSet))
}
