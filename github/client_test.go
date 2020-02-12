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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func withResponse(t *testing.T, handle func(w http.ResponseWriter), action func(client)) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handle(w)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := client{
		apiBase:   server.URL,
		authToken: "token",
	}

	action(client)
}
func TestGetJSON_NoData(t *testing.T) {
	handle := func(w http.ResponseWriter) {
		w.WriteHeader(200)
	}
	withResponse(t, handle, func(c client) {
		_, err := c.getJSON("/foo", nil, nil)
		assert.Error(t, err)
	})
}

func TestGetJSON_OK(t *testing.T) {
	data := map[string]int{}
	handle := func(w http.ResponseWriter) {
		w.Write([]byte(`{"answer": 42}`))
	}

	withResponse(t, handle, func(c client) {
		_, err := c.getJSON("/foo", nil, &data)
		require.NoError(t, err)
		assert.Equal(t, 42, data["answer"])
	})
}

func TestPostJSON_OK(t *testing.T) {
	headers := map[string]string{
		"key": "value",
	}
	data := map[string]string{}
	handle := func(w http.ResponseWriter) {
		w.Write([]byte(`{"hello": "world"}`))
	}

	withResponse(t, handle, func(c client) {
		_, err := c.postJSON("/url", nil, headers, &data)

		require.NoError(t, err)
		assert.Equal(t, "world", data["hello"])
	})
}
