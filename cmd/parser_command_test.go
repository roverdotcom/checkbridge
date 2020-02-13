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

package cmd

import (
	"errors"
	"testing"

	"github.com/roverdotcom/checkbridge/github"
	"github.com/roverdotcom/checkbridge/parser"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type stubClient struct {
	reportedCheck *github.CheckRun
	err           *error
}

func (s *stubClient) CreateCheck(check github.CheckRun) error {
	s.reportedCheck = &check
	if s.err != nil {
		return *s.err
	}
	return nil
}

func TestAPIClient_NoToken(t *testing.T) {
	r := environment{
		vip: viper.New(),
	}

	_, err := r.apiClient(repo{})
	assert.Error(t, err)
}

func TestAPIClient_ValidToken(t *testing.T) {
	vip := viper.New()
	vip.Set("github-token", "token")
	r := environment{
		vip: vip,
	}

	_, err := r.apiClient(repo{})
	assert.NoError(t, err)
}

func TestReportResults_NoViolations(t *testing.T) {
	api := &stubClient{}
	result := parser.Result{}
	p := parseRunner{}
	code := p.reportResults(github.CheckRun{}, result, api)

	assert.Equal(t, code, 0)
	assert.Equal(t, github.CheckConclusionSuccess, api.reportedCheck.Conclusion)
}

func TestReportResults_WithViolations(t *testing.T) {
	api := &stubClient{}
	result := parser.Result{
		Annotations: []parser.Annotation{{
			Path: "main.go",
		}},
	}
	e := environment{vip: viper.GetViper()}
	p := parseRunner{environment: e}
	code := p.reportResults(github.CheckRun{}, result, api)

	assert.Equal(t, code, 1)
	assert.Equal(t, github.CheckConclusionFailure, api.reportedCheck.Conclusion)
}

func TestReportResults_WithViolations_ExitZero(t *testing.T) {
	vip := viper.New()
	vip.Set("exit-zero", true)
	api := &stubClient{}
	result := parser.Result{
		Annotations: []parser.Annotation{{}},
	}
	p := parseRunner{
		environment: environment{vip: vip},
	}
	code := p.reportResults(github.CheckRun{}, result, api)

	assert.Equal(t, code, 0)
	assert.Equal(t, github.CheckConclusionFailure, api.reportedCheck.Conclusion)
}

func TestReportResults_WithViolations_AnnotateOnly(t *testing.T) {
	vip := viper.New()
	vip.Set("annotate-only", true)
	api := &stubClient{}
	result := parser.Result{
		Annotations: []parser.Annotation{{}},
	}
	p := parseRunner{
		environment: environment{vip: vip},
	}
	p.reportResults(github.CheckRun{}, result, api)

	assert.Equal(t, github.CheckConclusionNeutral, api.reportedCheck.Conclusion)
}

func TestReportResults_GitHubError(t *testing.T) {
	err := errors.New("unicorns")
	api := &stubClient{
		err: &err,
	}
	result := parser.Result{}
	p := parseRunner{
		environment: environment{},
	}
	code := p.reportResults(github.CheckRun{}, result, api)

	assert.Equal(t, code, 5)
}
