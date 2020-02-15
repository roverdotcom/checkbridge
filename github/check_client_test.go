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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/roverdotcom/checkbridge/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockHandler struct {
	called bool
	status int
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.called = true
	w.WriteHeader(m.status)
	w.Write([]byte(`{}`))
}

func createHandler(status int) mockHandler {
	return mockHandler{
		status: status,
	}
}

func createCheckWithRun(handler http.Handler, run CheckRun) error {
	server := httptest.NewServer(handler)
	defer server.Close()

	client := checkClient{
		client: client{
			apiBase:   server.URL,
			authToken: "fake-token",
		},
		repo:  "repo",
		owner: "owner",
	}

	return client.CreateCheck(run)
}

func createCheck(handler http.Handler) error {
	return createCheckWithRun(handler, CheckRun{
		Name: "my-name",
	})
}

func TestNewCheck_client(t *testing.T) {
	r := dummyRepo{"owner", "repo"}
	c := NewCheckClient("token", r)
	assert.Equal(t, c.(checkClient).owner, "owner")
}

func TestCreateCheck_OK(t *testing.T) {
	handler := createHandler(201)
	err := createCheck(&handler)
	assert.NoError(t, err)
	assert.True(t, handler.called)
}

func TestCreateCheck_NotFound(t *testing.T) {
	handler := createHandler(404)
	err := createCheck(&handler)
	assert.Error(t, err)
}
func TestCreateCheck_ManyAnnotations(t *testing.T) {
	sentRun := CheckRun{}
	annotations := make([]parser.Annotation, 100)
	run := CheckRun{
		Output: parser.Result{
			Annotations: annotations,
		},
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte(`{}`))
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		require.NoError(t, decoder.Decode(&sentRun))
	})

	err := createCheckWithRun(handler, run)

	require.NoError(t, err, "error sending check")
	assert.Equal(t, 50, len(sentRun.Output.Annotations), "expected large annotation list truncated to 50")
}

func TestCreateCheck_BadURL(t *testing.T) {
	c := checkClient{
		client: client{
			apiBase: "gopher://",
		},
	}
	assert.Error(t, c.CreateCheck(CheckRun{}))
}
