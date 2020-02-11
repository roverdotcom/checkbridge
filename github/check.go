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

import "github.com/roverdotcom/checkbridge/parser"

// CheckStatus represents the status of a check (ongoing, completed)
type CheckStatus string

// CheckConclusion represents the conclusion of a check (success, failure)
type CheckConclusion string

const (
	// CheckStatusCompleted represents a completed check
	CheckStatusCompleted CheckStatus = "completed"

	// CheckConclusionSuccess means the check was successful
	CheckConclusionSuccess CheckConclusion = "success"
	// CheckConclusionFailure means the check failed
	CheckConclusionFailure CheckConclusion = "failure"
	// CheckConclusionNeutral means the check completed with info, but did not fail
	CheckConclusionNeutral CheckConclusion = "neutral"
)

// CheckRun represents the results (intermediate or complete) of a check run
type CheckRun struct {
	Name       string          `json:"name"`
	HeadSHA    string          `json:"head_sha"`
	Status     CheckStatus     `json:"status"`
	Conclusion CheckConclusion `json:"conclusion,omitempty"`
	DetailsURL string          `json:"details_url,omitempty"`
	Output     parser.Result   `json:"output,omitempty"`
}
