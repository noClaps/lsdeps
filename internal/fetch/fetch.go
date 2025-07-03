package fetch

import (
	"encoding/json"
	"net/http"
)

func Fetch[T any](url string) (*T, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	pkg := new(T)
	err = json.NewDecoder(res.Body).Decode(pkg)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}
