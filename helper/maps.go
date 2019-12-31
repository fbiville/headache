/*
 * Copyright 2019 Florent Biville (@fbiville)
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

import "sort"

// extracts string keys in ascending order
func Keys(m map[string]string) []string {
	result := make([]string, len(m))
	i := 0
	for k := range m {
		result[i] = k
		i++
	}
	sort.Strings(result)
	return result
}
