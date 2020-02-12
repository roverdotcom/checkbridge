// Copyright (c) 2020 Rover.com
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package github

import (
	"errors"
	"fmt"
	"strconv"
)

var tokenAuthHeaders = map[string]string{
	"Accept": "application/vnd.github.machine-man-preview+json",
}

type tokenClient struct {
	client

	jwt string
}

type installationResponse struct {
	ID          int               `json:"id"`
	Permissions map[string]string `json:"permissions"`
}

func (t tokenClient) installationID(r Repo) (string, error) {
	url := fmt.Sprintf("repos/%s/%s/installation", r.Owner(), r.Name())
	installation := installationResponse{}
	_, err := t.getJSON(url, tokenAuthHeaders, &installation)
	if err != nil {
		return "", err
	}

	if installation.ID == 0 {
		return "", errors.New("no installation ID returned")
	}

	return strconv.Itoa(installation.ID), nil
}

func (t tokenClient) accessTokenURL(installationID string) string {
	return fmt.Sprintf("app/installations/%s/access_tokens", installationID)
}

type accessTokenResponse struct {
	Token       string            `json:"token"`
	Permissions map[string]string `json:"permissions"`
}

type accessTokenRequest struct {
	Permissions map[string]string `json:"permissions"`
}

func (t tokenClient) getAccesssToken(installationID string, perms map[string]string) (string, error) {
	requestData := accessTokenRequest{
		Permissions: perms,
	}

	tokenResp := accessTokenResponse{}

	resp, err := t.postJSON(t.accessTokenURL(installationID), requestData, tokenAuthHeaders, &tokenResp)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 201 {
		return "", fmt.Errorf("non-201 status code %s", resp.Status)
	}

	if tokenResp.Token == "" {
		return "", errors.New("no token in response")
	}

	return tokenResp.Token, nil
}
