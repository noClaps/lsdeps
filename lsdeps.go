package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"
	"sync"
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
		for dep, version := range packageData.Dependencies {
			deps[dep] = version
		}
	}
	if !skipPeer && len(packageData.PeerDependencies) != 0 {
		for dep, version := range packageData.PeerDependencies {
			deps[dep] = version
		}
	}
	if !skipOptional && len(packageData.OptionalDependencies) != 0 {
		for dep, version := range packageData.OptionalDependencies {
			deps[dep] = version
		}
	}

	return deps, nil
}

func logErrorf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("\033[31m%s\033[0m", format), a...)
}

func findArg(argv []string) (string, []string) {
	for i := range argv {
		if (i == 0 && argv[i][0] != '-') || (i-1 >= 0 && argv[i-1][0] != '-') {
			return argv[i], slices.Concat(argv[:i], argv[i+1:])
		}
	}

	return "", argv
}

func main() {
	var skipOptional bool
	flag.BoolVar(&skipOptional, "skip-optional", false, "Skip counting optional dependencies.")
	flag.BoolVar(&skipOptional, "o", false, "Skip counting optional dependencies.")

	var skipPeer bool
	flag.BoolVar(&skipPeer, "skip-peer", false, "Skip counting peer dependencies.")
	flag.BoolVar(&skipPeer, "p", false, "Skip counting peer dependencies.")

	var version string
	flag.StringVar(&version, "version", "latest", "The version of the package being fetched.")

	var help bool
	flag.BoolVar(&help, "help", false, "Display this help message and exit.")
	flag.BoolVar(&help, "h", false, "Display this help message and exit.")

	packageName, remainingArgs := findArg(os.Args[1:])

	err := flag.CommandLine.Parse(remainingArgs)

	if help || err != nil {
		fmt.Printf(`
USAGE: lsdeps <package> [--skip-optional] [--skip-peer] [--version <version>]

ARGUMENTS:
  <package>              The npm package to count dependencies for.

OPTIONS:
  --skip-optional, -o    Skip counting optional dependencies.
  --skip-peer, -p        Skip counting peer dependencies.
  --version <version>    The version of the package being fetched.
  --help, -h             Display this help message and exit.

`)
		return
	}

	if packageName == "" {
		fmt.Println("USAGE: lsdeps <package> [--skip-optional] [--skip-peer] [--version <version>]")
		os.Exit(1)
	}

	fmt.Printf("Fetching dependencies for %s@%s", packageName, version)

	depSet := map[string]bool{}
	if len(version) >= 4 && version[:4] == "npm:" {
		actualPackage := strings.SplitN(version[4:], "@", 2)
		packageName = actualPackage[0]
		version = actualPackage[1]
	}

	queue, err := getDeps(packageName, skipPeer, skipOptional, version)
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
				deps, err := getDeps(setPackage, skipPeer, skipOptional, setPackageVersion)
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
