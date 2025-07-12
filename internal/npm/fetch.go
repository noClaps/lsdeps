package npm

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var client = http.Client{}

func Fetch(name string) (*NpmPackage, error) {
	url := "https://registry.npmjs.org/" + name + "/latest"
	res, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%v %s", err, url)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s %s", res.StatusCode, res.Status, url)
	}

	pkg := new(NpmPackage)
	err = json.NewDecoder(res.Body).Decode(pkg)
	if err != nil {
		return nil, fmt.Errorf("%v %s", err, url)
	}

	return pkg, nil
}
