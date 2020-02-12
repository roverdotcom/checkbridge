// Copyright (c) 2020 Rover.com
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

const apiBase = "https://api.github.com"

type client struct {
	authToken string
	apiBase   string
}

func (c client) addAuthHeader(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.authToken))
}

func (c client) decodeResponse(req *http.Request, result interface{}) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(result); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c client) getJSON(url string, headers map[string]string, result interface{}) (*http.Response, error) {
	fullURL := fmt.Sprintf("%s/%s", c.apiBase, url)
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}
	c.addAuthHeader(req)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	return c.decodeResponse(req, result)
}

func (c client) postJSON(url string, body interface{}, headers map[string]string, result interface{}) (*http.Response, error) {
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return nil, err
	}
	fullURL := fmt.Sprintf("%s/%s", c.apiBase, url)
	logrus.WithField("url", fullURL).WithField("body", buf.String()).Debug("Making HTTP request to GitHub API")

	req, err := http.NewRequest(http.MethodPost, fullURL, &buf)
	if err != nil {
		return nil, err
	}

	c.addAuthHeader(req)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	return c.decodeResponse(req, result)
}
