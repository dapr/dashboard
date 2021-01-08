// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------
package log

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogReader(t *testing.T) {
	scenarios := []struct {
		container   string
		content     string
		expectedLog *Log
		expectError bool
	}{
		{
			"daprd",
			`2021-01-06T17:13:34.850421576Z time="2021-01-06T17:13:34.815684553Z" level=debug msg="certificate signed successfully" app_id=nodeapp instance=nodeapp-X-Y scope=dapr.runtime.grpc.internal type=log ver=edge`,
			&Log{
				Container: "daprd",
				Content:   `time="2021-01-06T17:13:34.815684553Z" level=debug msg="certificate signed successfully" app_id=nodeapp instance=nodeapp-X-Y scope=dapr.runtime.grpc.internal type=log ver=edge`,
				Timestamp: 1609953214850421576,
				Level:     "debug",
			},
			false,
		},
		{
			"myapp",
			`2021-01-06T17:13:34.850421577Z this is my app`,
			&Log{
				Container: "myapp",
				Content:   `this is my app`,
				Timestamp: 1609953214850421577,
				Level:     "info",
			},
			false,
		},
	}

	for _, scenario := range scenarios {
		reader := NewReader(scenario.container, newStringReader(scenario.content))
		logRecord, err := reader.ReadLog()
		if scenario.expectError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}

		assert.Equal(t, scenario.expectedLog, logRecord)
	}
}

func newStringReader(s string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(s))
}
