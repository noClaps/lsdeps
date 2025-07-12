package main

import (
	"fmt"
	"sync"

	"github.com/noclaps/applause"
	"github.com/noclaps/lsdeps/internal/logger"
	"github.com/noclaps/lsdeps/internal/npm"
)

type args struct {
	Package      string `help:"The npm package to count dependencies for."`
	SkipOptional bool   `type:"option" short:"o" help:"Skip counting optional dependencies."`
	SkipPeer     bool   `type:"option" short:"p" help:"Skip counting peer dependencies"`
}

func main() {
	args := args{}
	err := applause.Parse(&args)
	if err != nil {
		logger.Fatalln(err)
	}

	name := args.Package

	var mu sync.Mutex

	toFetch := map[string]struct{}{name: {}}
	toProcess := map[*npm.NpmPackage]struct{}{}
	deps := map[string]struct{}{}

	fmt.Print("Fetching dependencies...")

	for len(toFetch) > 0 {
		var fetchWg sync.WaitGroup
		for name := range toFetch {
			fetchWg.Add(1)
			go func() {
				defer fetchWg.Done()

				pkg, err := npm.Fetch(name)
				if err != nil {
					logger.Errorln(err)
					return
				}

				mu.Lock()
				if _, ok := toProcess[pkg]; !ok {
					toProcess[pkg] = struct{}{}
				}
				mu.Unlock()
			}()
		}
		fetchWg.Wait()
		clear(toFetch)

		var procWg sync.WaitGroup
		for pkg := range toProcess {
			procWg.Add(1)
			go func() {
				defer procWg.Done()

				pkgsDeps := npm.GetDeps(pkg, args.SkipOptional, args.SkipPeer)

				mu.Lock()
				for _, dep := range pkgsDeps {
					if _, ok := deps[dep]; !ok {
						deps[dep] = struct{}{}
					}
					if _, ok := toFetch[dep]; !ok {
						toFetch[dep] = struct{}{}
					}
				}
				mu.Unlock()
			}()
		}
		procWg.Wait()
		clear(toProcess)
	}

	fmt.Printf("\033[2K\r")
	fmt.Printf(`Name: %s
URL: https://npmjs.com/package/%s
Dependency count: %d
`, name, name, len(deps))
}
