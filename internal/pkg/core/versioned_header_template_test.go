/*
 * Copyright 2019 Florent Biville (@fbiville)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core_test

import (
	. "github.com/fbiville/headache/internal/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Versioned header template", func() {

	It("requires a full file scan when no revision is set", func() {
		template := VersionedHeaderTemplate{Revision: ""}

		Expect(template.RequiresFullScan()).To(BeTrue())
	})

	It("requires a full file scan if previous and current contents do not match", func() {
		template := VersionedHeaderTemplate{
			Revision: "some-sha",
			Current:  template("current-contents", map[string]string{}),
			Previous: template("previous-contents", map[string]string{}),
		}

		Expect(template.RequiresFullScan()).To(BeTrue())
	})

	It("requires a full file scan if previous and current data keys do not match", func() {
		template := VersionedHeaderTemplate{
			Revision: "some-sha",
			Current:  template("same-contents", map[string]string{"foo": ""}),
			Previous: template("same-contents", map[string]string{"baz": ""}),
		}

		Expect(template.RequiresFullScan()).To(BeTrue())
	})

	It("does not require a full file scan if revision is set and contents+data keys match", func() {
		template := VersionedHeaderTemplate{
			Revision: "some-sha",
			Current:  template("same-contents", map[string]string{"foo": ""}),
			Previous: template("same-contents", map[string]string{"foo": ""}),
		}

		Expect(template.RequiresFullScan()).To(BeFalse())
	})
})
