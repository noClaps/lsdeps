package npm

import (
	"fmt"
	"maps"
	"strings"

	"github.com/noclaps/lsdeps/internal/fetch"
)

type npmPackage struct {
	Dependencies         map[string]string `json:"dependencies,omitempty"`
	OptionalDependencies map[string]string `json:"optionalDependencies,omitempty"`
	PeerDependencies     map[string]string `json:"peerDependencies,omitempty"`
}

func GetDeps(name string, skipPeer bool, skipOptional bool, version string) (map[string]string, error) {
	if len(version) >= 4 && version[:4] == "npm:" {
		actualPackage := strings.SplitN(version[4:], "@", 2)

		name = actualPackage[0]
		version = actualPackage[1]
	}

	version = parseVersion(version)

	pkg, err := fetch.Fetch[npmPackage](fmt.Sprintf("https://registry.npmjs.com/%s/%s", name, version))
	if err != nil {
		return nil, err
	}

	totalDepsLen := len(pkg.Dependencies)
	if !skipPeer {
		totalDepsLen += len(pkg.PeerDependencies)
	}
	if !skipOptional {
		totalDepsLen += len(pkg.OptionalDependencies)
	}
	deps := make(map[string]string, totalDepsLen)
	maps.Copy(deps, pkg.Dependencies)
	if !skipPeer {
		maps.Copy(deps, pkg.PeerDependencies)
	}
	if !skipOptional {
		maps.Copy(deps, pkg.OptionalDependencies)
	}

	return deps, nil
}
