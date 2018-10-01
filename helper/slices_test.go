/*
 * Copyright 2018 Florent Biville (@fbiville)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package helper

import (
	"testing"
	"testing/quick"
)

func TestSlicePrepend(t *testing.T) {
	f := func(head string, tail []string) bool {
		result := PrependString(head, tail)
		return len(result) == 1+len(tail) && result[0] == head && SliceEqual(result[1:], tail)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func SliceEqual(slice1 []string, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}
