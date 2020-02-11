// MIT License
//
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
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockHandler struct {
	called bool
	status int
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.called = true
	w.WriteHeader(m.status)
}

func createHandler(status int) mockHandler {
	return mockHandler{
		status: status,
	}
}

func createCheck(handler *mockHandler) error {
	server := httptest.NewServer(handler)
	defer server.Close()

	client := client{
		apiBase: server.URL,
		token:   "fake-token",
		repo:    "repo",
		owner:   "owner",
	}

	return client.CreateCheck(CheckRun{
		Name: "my-name",
	})
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
