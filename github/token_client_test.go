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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dummyRepo struct {
	owner, name string
}

func (d dummyRepo) Owner() string { return d.owner }
func (d dummyRepo) Name() string  { return d.name }

func TestInstallationID_NoData(t *testing.T) {
	handle := func(w http.ResponseWriter) {
		w.WriteHeader(201)
	}
	withResponse(t, handle, func(c client) {
		tc := tokenClient{
			client: c,
			jwt:    "jwt",
		}

		_, err := tc.installationID(dummyRepo{})
		assert.Error(t, err)
	})
}

func TestInstallationID_Empty(t *testing.T) {
	handle := func(w http.ResponseWriter) {
		w.WriteHeader(201)
		w.Write([]byte(`{}`))

	}
	withResponse(t, handle, func(c client) {
		tc := tokenClient{
			client: c,
			jwt:    "jwt",
		}

		_, err := tc.installationID(dummyRepo{})
		assert.Error(t, err)
	})
}

func TestInstallationID_Valid(t *testing.T) {
	handle := func(w http.ResponseWriter) {
		w.WriteHeader(201)
		w.Write([]byte(`{"id": 12345}`))
	}
	withResponse(t, handle, func(c client) {
		tc := tokenClient{
			client: c,
			jwt:    "jwt",
		}

		id, err := tc.installationID(dummyRepo{})
		require.NoError(t, err)
		assert.Equal(t, "12345", id)
	})
}

func TestGetAccessToken_OK(t *testing.T) {
	token := "v1.1234"
	handle := func(w http.ResponseWriter) {
		w.WriteHeader(201)
		w.Write([]byte(`{"token": "` + token + `"}`))
	}
	withResponse(t, handle, func(c client) {
		tc := tokenClient{
			client: c,
			jwt:    "jwt",
		}

		token, err := tc.getAccesssToken("fake-installation-id", nil)
		require.NoError(t, err)
		assert.Equal(t, token, token)
	})
}
