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

package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/roverdotcom/checkbridge/parser"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var regexCmd = &cobra.Command{
	Use:   "regex",
	Short: "Parse results via regular expression",
	Run: func(cmd *cobra.Command, args []string) {
		vip := viper.GetViper()
		configureLogging(vip)
		regex, err := regexp.Compile(vip.GetString("regex"))
		if err != nil {
			logrus.WithError(err).Error("Unable to compile regular expression")
			os.Exit(2)
		}
		extractor := makeExtractor(vip)
		parse := parser.NewRegexer(regex, extractor, os.Stdin)

		runner := parseRunner{
			environment: environment{
				vip: viper.GetViper(),
				env: os.Getenv,
			},
			name:  vip.GetString("name"),
			parse: parse,
		}
		runner.run()
	},
}

func init() {
	regexCmd.Flags().Bool("warn", false, "treat regex matches as warning (instead of error)")
	regexCmd.Flags().String("name", "", "check name (required)")
	regexCmd.Flags().String("regex", "", "regular expression (required)")
	regexCmd.Flags().Int("line-pos", 2, "position in regex for line (required)")
	regexCmd.Flags().Int("path-pos", 1, "position in regex for path (required)")
	regexCmd.Flags().Int("message-pos", 3, "position in regex for message (required)")
	regexCmd.Flags().Int("column-pos", 0, "position in regex for column")

	regexCmd.MarkFlagRequired("name")
	regexCmd.MarkFlagRequired("regex")
	regexCmd.MarkFlagRequired("line-pos")
	regexCmd.MarkFlagRequired("path-pos")
	regexCmd.MarkFlagRequired("message-pos")
	viper.BindPFlags(regexCmd.Flags())
}

func makeExtractor(vip *viper.Viper) parser.AnnotationExtracter {
	linePos := vip.GetInt("line-pos")
	pathPos := vip.GetInt("path-pos")
	messagePos := vip.GetInt("message-pos")
	columnPos := vip.GetInt("column-pos")
	level := parser.LevelError

	if vip.GetBool("warn") {
		level = parser.LevelWarning
	}

	return func(matches []string) (parser.Annotation, error) {
		column := 0
		var err error
		if columnPos > 0 && len(matches) > columnPos {
			column, err = strconv.Atoi(matches[columnPos])
			if err != nil {
				return parser.Annotation{}, err
			}
		}

		if pathPos >= len(matches) {
			return parser.Annotation{}, fmt.Errorf("path regex position out of bounds: %d", pathPos)
		}
		path := matches[pathPos]

		if linePos >= len(matches) {
			return parser.Annotation{}, fmt.Errorf("line regex position out of bounds: %d", linePos)
		}
		line, err := strconv.Atoi(matches[linePos])
		if err != nil {
			return parser.Annotation{}, fmt.Errorf("error parsing line: %w", err)
		}

		if messagePos >= len(matches) {
			return parser.Annotation{}, fmt.Errorf("message regex position out of bounds: %d", messagePos)
		}
		message := matches[messagePos]

		return parser.Annotation{
			Level:   level,
			Column:  column,
			Path:    path,
			Line:    line,
			EndLine: line,
			Message: message,
		}, nil
	}
}
