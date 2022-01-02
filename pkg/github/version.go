// SPDX-License-Identifier: Apache-2.0

package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	releasesURL = "https://api.github.com/repos/%s/%s/releases?per_page=1"
)

func GetLatestVersion(org, repo string) (string, error) {
	client := http.Client{Timeout: time.Second * 3}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(releasesURL, org, repo), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "lfedge.org/eve-cli")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var data []map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}
	if len(data) < 1 {
		return "", errors.New("no releases found")
	}
	tagNameValue, ok := data[0]["tag_name"]
	if !ok {
		return "", errors.New("no tag name provided")
	}
	tagName, ok := tagNameValue.(string)
	if !ok {
		return "", errors.New("tag name unrecognizable")
	}
	return tagName, nil
}
