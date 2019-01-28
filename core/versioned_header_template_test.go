package core_test

import (
	. "github.com/fbiville/headache/core"
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
