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
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

const rsaHeader = `-----BEGIN RSA PRIVATE KEY-----`
const jwtExpiry = 60 * time.Second

// AuthProvider handles getting a GitHub token for the given permissions
type AuthProvider interface {
	GetToken(r Repo, perms map[string]string) (string, error)
}

// ConfigProvider is an interface over *viper.Viper
type ConfigProvider interface {
	GetString(string) string
}

type githubAuth struct {
	config  ConfigProvider
	apiBase string
}

func NewAuthProvider(c ConfigProvider) AuthProvider {
	return githubAuth{
		config:  c,
		apiBase: apiBase,
	}
}

func (g githubAuth) readPrivateKey(pathOrKey string) (*rsa.PrivateKey, error) {
	var bytePayload []byte
	if strings.HasPrefix(pathOrKey, rsaHeader) {
		bytePayload = []byte(pathOrKey)
	} else {
		logrus.WithField("path", pathOrKey).Debug("Reading private key from file")

		data, err := ioutil.ReadFile(pathOrKey)
		if err != nil {
			return nil, err
		}
		bytePayload = data
	}
	return jwt.ParseRSAPrivateKeyFromPEM(bytePayload)
}

func (g githubAuth) makeJWT() (string, error) {
	applicationID := g.config.GetString("application-id")
	if applicationID == "0" || applicationID == "" {
		return "", errors.New("no application ID provided")
	}

	privateKey := g.config.GetString("private-key")
	if privateKey == "" {
		return "", errors.New("no private key provided")
	}

	rsaKey, err := g.readPrivateKey(privateKey)
	if err != nil {
		return "", fmt.Errorf("error reading private key: %w", err)
	}

	now := time.Now()
	exp := now.Add(jwtExpiry)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": applicationID,
		"iat": now.Unix(),
		"exp": exp.Unix(),
	})

	return token.SignedString(rsaKey)
}

func (g githubAuth) GetToken(r Repo, perms map[string]string) (string, error) {
	if token := g.config.GetString("github-token"); token != "" {
		logrus.Debug("Using explicit GitHub token, skipping JWT exchange")
		return token, nil
	}

	appJWT, err := g.makeJWT()
	if err != nil {
		return "", err
	}
	logrus.WithField("jwt", appJWT).Debug("Got JWT")

	installationID := g.config.GetString("installation-id")
	tc := tokenClient{
		client: client{
			authToken: appJWT,
			apiBase:   g.apiBase,
		},
	}

	if installationID == "0" || installationID == "" {
		logrus.Debug("No installation ID provided, asking GitHub")
		installationID, err = tc.installationID(r)
		if err != nil {
			return "", err
		}
		logrus.WithField("installationID", installationID).Debug("Got installation ID response from GitHub")
	}

	return tc.getAccesssToken(installationID, perms)
}
