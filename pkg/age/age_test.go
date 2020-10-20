// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package age

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAge(t *testing.T) {
	current_time := time.Now()
	assert.Equal(t, "47s", GetAge(current_time.Add(-time.Second * 47)))
	assert.Equal(t, "4m", GetAge(current_time.Add(-time.Minute * 4)))
	assert.Equal(t, "2h", GetAge(current_time.Add(-time.Hour * 2)))
	assert.Equal(t, "3d", GetAge(current_time.Add(-time.Hour * 76)))
}

