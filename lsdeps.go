package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/noclaps/applause"
)

type Package struct {
	Dependencies         map[string]string `json:"dependencies,omitempty"`
	OptionalDependencies map[string]string `json:"optionalDependencies,omitempty"`
	PeerDependencies     map[string]string `json:"peerDependencies,omitempty"`
}

func fetch(url string) (*Package, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	pkg := new(Package)
	err = json.NewDecoder(res.Body).Decode(pkg)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}

func parseVersion(version string) string {
	re := regexp.MustCompile(`^([0-9]\.[0-9]\.[0-9])(-(alpha|beta|rc)\.[0-9]+)?`)
	if re.MatchString(version) {
		return version
	}

	if version[0] == '~' || version[0] == '^' {
		if re.MatchString(version[1:]) {
			return version[1:]
		}
	}

	if version == "next" {
		return version
	}

	return "latest"
}

func getDeps(name string, skipPeer bool, skipOptional bool, version string) (map[string]string, error) {
	if len(version) >= 4 && version[:4] == "npm:" {
		actualPackage := strings.SplitN(version[4:], "@", 2)

		name = actualPackage[0]
		version = actualPackage[1]
	}

	version = parseVersion(version)

	deps := make(map[string]string)
	packageData, err := fetch(fmt.Sprintf("https://registry.npmjs.com/%s/%s", name, version))
	if err != nil {
		packageData, err = fetch(fmt.Sprintf("https://registry.npmjs.com/%s/latest", name))
		if err != nil {
			return nil, err
		}
	}

	if len(packageData.Dependencies) != 0 {
		maps.Copy(deps, packageData.Dependencies)
	}
	if !skipPeer && len(packageData.PeerDependencies) != 0 {
		maps.Copy(deps, packageData.PeerDependencies)
	}
	if !skipOptional && len(packageData.OptionalDependencies) != 0 {
		maps.Copy(deps, packageData.OptionalDependencies)
	}

	return deps, nil
}

func logErrorf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("\033[31m%s\033[0m", format), a...)
}

type Args struct {
	Package      string `help:"The npm package to count dependencies for."`
	SkipOptional bool   `type:"option" short:"o" help:"Skip counting optional dependencies."`
	SkipPeer     bool   `type:"option" short:"p" help:"Skip counting peer dependencies"`
	Version      string `type:"option" help:"The version of the package being fetched."`
}

func main() {
	args := Args{Version: "latest"}
	err := applause.Parse(&args)

	packageName := args.Package
	version := args.Version
	fmt.Printf("Fetching dependencies for %s@%s", packageName, version)

	depSet := map[string]bool{}
	if len(version) >= 4 && version[:4] == "npm:" {
		actualPackage := strings.SplitN(version[4:], "@", 2)
		packageName = actualPackage[0]
		version = actualPackage[1]
	}

	queue, err := getDeps(packageName, args.SkipPeer, args.SkipOptional, version)
	if err != nil {
		logErrorf("\nERROR: Package %s@%s does not exist\n", packageName, version)
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
				if depSet[setPackage] {
					mu.Unlock()
					return
				}
				depSet[setPackage] = true
				mu.Unlock()

				fmt.Printf("\033[2K\rFetching dependencies for %s@%s", setPackage, setPackageVersion)
				deps, err := getDeps(setPackage, args.SkipPeer, args.SkipOptional, setPackageVersion)
				if err != nil {
					logErrorf("\nERROR: Package %s@%s does not exist\n", setPackage, setPackageVersion)
					return
				}

				mu.Lock()
				for dep, version := range deps {
					if !depSet[dep] {
						queue[dep] = version
					}
				}
				mu.Unlock()
			}()
		}
		wg.Wait()
	}

	fmt.Printf("\033[2K\r")
	fmt.Printf(`
Name: %s
URL: https://npmjs.com/package/%s/v/%s
Dependency count: %d

`, packageName, packageName, version, len(depSet))
}
