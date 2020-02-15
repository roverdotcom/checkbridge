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

type stubEnv struct {
	concreteEnv

	sc    *stubClient
	token *string
}

func (s stubEnv) apiClient(_ repo) (github.CheckClient, error) {
	if s.sc != nil {
		return s.sc, nil
	}
	return nil, errors.New("stub client missing")
}

func (s stubEnv) githubToken(_ repo) (string, error) {
	if s.token != nil {
		return *s.token, nil
	}
	return "", errors.New("no token")
}

func TestAPIClient_NoToken(t *testing.T) {
	r := concreteEnv{
		v: viper.New(),
	}

	_, err := r.apiClient(repo{})
	assert.Error(t, err)
}

func TestAPIClient_ValidToken(t *testing.T) {
	vip := viper.New()
	vip.Set("github-token", "token")
	r := concreteEnv{
		v: vip,
	}

	_, err := r.apiClient(repo{})
	assert.NoError(t, err)
}

func TestParseRunnerRun_NoRepo(t *testing.T) {
	p := parseRunner{
		environment: concreteEnv{
			v: viper.New(),
			e: envStub{}.get,
		},
	}
	assert.Equal(t, 3, p.run())
}

func TestParseRunnerRun_BadPrivateKey(t *testing.T) {
	vip := viper.New()
	vip.Set("commit-sha", "fake-sha")
	vip.Set("github-repo", "ghost/example")
	vip.Set("private-key", "non/existent/path.pem")

	p := parseRunner{
		environment: concreteEnv{
			v: vip,
			e: envStub{}.get,
		},
	}

	assert.Equal(t, 4, p.run())
}

type stubParser struct {
	r   parser.Result
	err error
}

func (s stubParser) Run() (parser.Result, error) {
	return s.r, s.err
}

func fakeRepoConfig() *viper.Viper {
	vip := viper.New()
	vip.Set("commit-sha", "fake-sha")
	vip.Set("github-repo", "ghost/example")
	vip.Set("github-token", "fake-token")
	vip.Set("mark-in-progress", true)

	return vip
}

func TestParseRunnerRun_ErrorSendingCheck(t *testing.T) {
	vip := fakeRepoConfig()

	err := errors.New("tubes are clogged")
	sc := stubClient{
		err: &err,
	}

	p := parseRunner{
		environment: stubEnv{
			concreteEnv: concreteEnv{
				v: vip,
				e: envStub{}.get,
			},
			sc: &sc,
		},
		parse: stubParser{},
	}

	assert.Equal(t, 5, p.run())
}

func TestParseRunnerRun_ErrorParsingResult(t *testing.T) {
	vip := fakeRepoConfig()

	sc := stubClient{}

	p := parseRunner{
		environment: stubEnv{
			concreteEnv: concreteEnv{
				v: vip,
				e: envStub{}.get,
			},
			sc: &sc,
		},
		parse: stubParser{
			err: errors.New("can't parse this"),
		},
	}

	assert.Equal(t, 3, p.run())
	assert.Equal(t, github.CheckConclusionFailure, sc.reportedCheck.Conclusion)
	assert.Equal(t, "can't parse this", sc.reportedCheck.Output.Summary)
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
	e := concreteEnv{v: viper.GetViper()}
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
		environment: concreteEnv{v: vip},
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
		environment: concreteEnv{v: vip},
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
		environment: concreteEnv{},
	}
	code := p.reportResults(github.CheckRun{}, result, api)

	assert.Equal(t, code, 5)
}

func TestSummaryResult_NoIssues(t *testing.T) {
	result := parser.Result{}
	assert.Equal(t, "no issues", summarizeResult(result))
}

func TestSummaryResult_TwoErrors(t *testing.T) {
	result := parser.Result{
		Annotations: []parser.Annotation{
			{Level: parser.LevelError},
			{Level: parser.LevelError},
		},
	}
	assert.Equal(t, "2 errors", summarizeResult(result))
}

func TestSummaryResult_OneOfEach(t *testing.T) {
	result := parser.Result{
		Annotations: []parser.Annotation{
			{Level: parser.LevelError},
			{Level: parser.LevelWarning},
		},
	}
	assert.Equal(t, "1 error and 1 warning", summarizeResult(result))
}

func TestSummaryResult_TwoWarnings(t *testing.T) {
	result := parser.Result{
		Annotations: []parser.Annotation{
			{Level: parser.LevelWarning},
			{Level: parser.LevelWarning},
		},
	}
	assert.Equal(t, "2 warnings", summarizeResult(result))
}

func TestCapitalizeFirstChar_TwoWords(t *testing.T) {
	assert.Equal(t, "No issues", capitalizeFirstChar("no issues"))
}
