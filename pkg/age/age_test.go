/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

