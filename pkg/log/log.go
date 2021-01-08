// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------
package log

import (
	"bufio"
	"regexp"
	"strings"
	"time"
)

// Log represents a log message with metadata
type Log struct {
	Level     string `json:"level"`
	Timestamp int64  `json:"timestamp"`
	Container string `json:"container"`
	Content   string `json:"content"`
}

// Reader reads logs line by line.
type Reader struct {
	levelExp  *regexp.Regexp
	timeExp   *regexp.Regexp
	container string
	reader    *bufio.Reader
}

// NewReader creates a reader that parses logs.
func NewReader(container string, reader *bufio.Reader) *Reader {
	levelExp, _ := regexp.Compile("(level=)[^ ]*")
	timeExp, _ := regexp.Compile("^[^ ]+")
	return &Reader{
		levelExp:  levelExp,
		timeExp:   timeExp,
		container: container,
		reader:    reader,
	}
}

// ReadLog reads a new log entry
func (r *Reader) ReadLog() (*Log, error) {
	bytes, _, err := r.reader.ReadLine()
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		return nil, nil
	}

	content := string(bytes)

	level := strings.ToLower(strings.Replace(r.levelExp.FindString(content), "level=", "", 1))
	if level == "" {
		level = "info"
	}

	timestampString := r.timeExp.FindString(content)
	timestamp, err := time.Parse(time.RFC3339Nano, timestampString)
	if err != nil {
		return nil, err
	}

	return &Log{
		Level:     level,
		Timestamp: timestamp.UnixNano(),
		Container: r.container,
		Content:   strings.TrimPrefix(strings.TrimPrefix(content, timestampString), " "),
	}, nil
}
