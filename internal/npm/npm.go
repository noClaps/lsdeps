package npm

type NpmPackage struct {
	Dependencies         map[string]string `json:"dependencies"`
	OptionalDependencies map[string]string `json:"optionalDependencies"`
	PeerDependencies     map[string]string `json:"peerDependencies"`
}

func GetDeps(pkg *NpmPackage, skipOptional, skipPeer bool) []string {
	deps := []string{}
	for name := range pkg.Dependencies {
		deps = append(deps, name)
	}
	if !skipOptional {
		for name := range pkg.OptionalDependencies {
			deps = append(deps, name)
		}
	}
	if !skipPeer {
		for name := range pkg.PeerDependencies {
			deps = append(deps, name)
		}
	}

	return deps
}
