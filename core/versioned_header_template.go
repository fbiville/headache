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

package core

import "github.com/fbiville/headache/helper"

type VersionedHeaderTemplate struct {
	Current  *HeaderTemplate
	Previous *HeaderTemplate
	Revision string
}

func (t VersionedHeaderTemplate) RequiresFullScan() bool {
	return t.Revision == "" ||
		!helper.SliceEqual(t.Current.Lines, t.Previous.Lines) ||
		!helper.SliceEqual(helper.Keys(t.Current.Data), helper.Keys(t.Previous.Data))
}
