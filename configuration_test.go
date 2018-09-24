package header_test

import (
	. "github.com/fbiville/header"
	. "github.com/onsi/gomega"
	"testing"
)

func TestConfigurationInitWithLineCommentStyle(t *testing.T) {
	I := NewGomegaWithT(t)

	configuration, err := NewConfiguration(
		"fixtures/license.txt",
		SlashSlash,
		[]string{"*.txt"},
		map[string]string{
			"Year":  "2018",
			"Owner": "ACME Labs",
		})

	I.Expect(err).To(BeNil())
	I.Expect(configuration.HeaderContents).To(Equal(`// Copyright 2018 ACME Labs
// 
// Some fictional license`))
	I.Expect(configuration.Includes).To(Equal([]string{"*.txt"}))
}

func TestConfigurationInitWithBlockCommentStyle(t *testing.T) {
	I := NewGomegaWithT(t)

	configuration, err := NewConfiguration(
		"fixtures/license.txt",
		SlashStar,
		[]string{"*.txt"},
		map[string]string{
			"Year":  "2018",
			"Owner": "ACME Labs",
		})

	I.Expect(err).To(BeNil())
	I.Expect(configuration.HeaderContents).To(Equal(`/* 
 * Copyright 2018 ACME Labs
 * 
 * Some fictional license
 */`))
	I.Expect(configuration.Includes).To(Equal([]string{"*.txt"}))
}
