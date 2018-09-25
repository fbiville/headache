package main_test

import (
	. "github.com/fbiville/header"
	. "github.com/onsi/gomega"
	"testing"
)

func TestConfigurationInitWithLineCommentStyle(t *testing.T) {
	I := NewGomegaWithT(t)

	configuration, err := ParseConfiguration(Configuration{
		HeaderFile:   "fixtures/license.txt",
		CommentStyle: "SlashSlash",
		Includes:     []string{"*.txt"},
		Excludes:     []string{},
		TemplateData: map[string]string{
			"Year":  "2018",
			"Owner": "ACME Labs",
		}})

	I.Expect(err).To(BeNil())
	I.Expect(configuration.HeaderContents).To(Equal(`// Copyright 2018 ACME Labs
//
// Some fictional license`))
	I.Expect(configuration.Includes).To(Equal([]string{"*.txt"}))
}

func TestConfigurationInitWithBlockCommentStyle(t *testing.T) {
	I := NewGomegaWithT(t)

	configuration, err := ParseConfiguration(Configuration{
		HeaderFile:   "fixtures/license.txt",
		CommentStyle: "SlashStar",
		Includes:     []string{"*.txt"},
		Excludes:     []string{},
		TemplateData: map[string]string{
			"Year":  "2018",
			"Owner": "ACME Labs",
		}})

	I.Expect(err).To(BeNil())
	I.Expect(configuration.HeaderContents).To(Equal(`/*
 * Copyright 2018 ACME Labs
 *
 * Some fictional license
 */`))
	I.Expect(configuration.Includes).To(Equal([]string{"*.txt"}))
}
