package npm

import "regexp"

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
