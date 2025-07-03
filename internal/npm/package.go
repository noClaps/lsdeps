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

	deps := make(map[string]string)
	packageData, err := fetch.Fetch[npmPackage](fmt.Sprintf("https://registry.npmjs.com/%s/%s", name, version))
	if err != nil {
		return nil, err
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
