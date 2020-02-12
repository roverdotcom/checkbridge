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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGithubAuth_GetTokenFromEnv(t *testing.T) {
	vip := viper.New()
	mytoken := "mytoken"
	vip.Set("github_token", mytoken)
	auth := NewAuthProvider(vip)

	token, err := auth.GetToken(dummyRepo{}, map[string]string{})
	require.NoError(t, err, "error getting token")
	assert.Equal(t, mytoken, token)
}

func TestGithubAuth_GetTokenNoEnv(t *testing.T) {
	auth := NewAuthProvider(viper.New())

	_, err := auth.GetToken(dummyRepo{}, nil)
	assert.Error(t, err)
}

func TestMakeJWT_NoAppId(t *testing.T) {
	conf := viper.New()
	gh := githubAuth{
		config: conf,
	}

	_, err := gh.makeJWT()
	assert.Error(t, err)
}

func TestMakeJWT_NoPrivateKey(t *testing.T) {
	conf := viper.New()
	conf.Set("application_id", "42")
	gh := githubAuth{
		config: conf,
	}

	_, err := gh.makeJWT()
	assert.Error(t, err)
}
func TestMakeJWT_InvalidPrivateKey(t *testing.T) {
	conf := viper.New()
	conf.Set("application_id", "42")
	conf.Set("private_key", "bad/path/to/pem")
	gh := githubAuth{
		config: conf,
	}

	_, err := gh.makeJWT()
	assert.Error(t, err)
}

const testPrivateKey = "testdata/key.pem"

func TestMakeJWT_OK(t *testing.T) {
	appID := "42"
	conf := viper.New()
	conf.Set("application_id", appID)
	conf.Set("private_key", testPrivateKey)
	gh := githubAuth{
		config: conf,
	}

	jwtVal, err := gh.makeJWT()
	require.NoError(t, err)
	assert.NotEmpty(t, jwtVal)

	parser := &jwt.Parser{}
	claims := jwt.MapClaims{}
	_, _, err = parser.ParseUnverified(jwtVal, claims)
	require.NoError(t, err)
}

func TestGetToken_EndToEnd(t *testing.T) {
	stubToken := "stub.token"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ensure auth for all calls
		bearer := r.Header.Get("Authorization")
		assert.Contains(t, bearer, "Bearer ")
		assert.NotEmpty(t, r.Header.Get("Accept"), "missing Accept header")

		if strings.HasSuffix(r.URL.Path, "/installation") {
			w.Write([]byte(`{"id":42}`))
		} else if strings.HasSuffix(r.URL.Path, "/access_tokens") {
			assert.Contains(t, r.URL.Path, "installations/42/", "missing installation ID")
			w.WriteHeader(201)
			w.Write([]byte(`{"token":"` + stubToken + `"}`))
		} else {
			assert.Fail(t, "unexpected URL hit:"+r.URL.Path)
		}
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	vip := viper.New()
	vip.Set("private_key", testPrivateKey)
	vip.Set("application_id", "my-id")
	gh := githubAuth{
		config:  vip,
		apiBase: server.URL,
	}

	repo := dummyRepo{"org", "repo"}
	perms := map[string]string{
		"foo": "bar",
	}
	token, err := gh.GetToken(repo, perms)
	require.NoError(t, err)
	assert.Equal(t, stubToken, token)
}
