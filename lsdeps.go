package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
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

var args struct {
	Package      string `arg:"positional,required" help:"The npm package to count dependencies for."`
	SkipOptional bool   `arg:"-o,--skip-optional" help:"Skip counting optional dependencies."`
	SkipPeer     bool   `arg:"-p,--skip-peer" help:"Skip counting peer dependencies."`
	Version      string `help:"The version of the package being fetched."`
	Help         bool   `arg:"-h,--help" help:"Display this help message and exit"`
}

func parseArgs(argv []string) {
	for i := range argv {
		if argv[i][0] == '-' {
			// Flag or option
			switch argv[i] {
			case "-o", "--skip-optional":
				args.SkipOptional = true
			case "-p", "--skip-peer":
				args.SkipPeer = true
			case "-h", "--help":
				args.Help = true
			case "--version":
				args.Version = argv[i+1]
				i += 2
			}
		} else {
			if i == 0 || argv[i-1] != "--version" {
				args.Package = argv[i]
			}
		}
	}
}

func main() {
	parseArgs(os.Args[1:])

	if args.Help {
		fmt.Printf(`
USAGE: lsdeps <package> [--skip-optional] [--skip-peer] [--silent] [--version <version>]

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

	if args.Version == "" {
		args.Version = "latest"
	}

	fmt.Printf("Fetching dependencies for %s@%s", args.Package, args.Version)

	depSet := map[string]bool{}
	if len(args.Version) >= 4 && args.Version[:4] == "npm:" {
		actualPackage := strings.SplitN(args.Version[4:], "@", 2)
		args.Package = actualPackage[0]
		args.Version = actualPackage[1]
	}

	queue, err := getDeps(args.Package, args.SkipPeer, args.SkipOptional, args.Version)
	if err != nil {
		logErrorf("\nERROR: Package %s@%s does not exist\n", args.Package, args.Version)
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

`, args.Package, args.Package, args.Version, len(depSet))
}
