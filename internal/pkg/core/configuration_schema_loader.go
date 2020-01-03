package core

import (
	json_schema "github.com/xeipuuv/gojsonschema"
	"log"
)

type JsonSchemaLoader interface {
	// loads specified JSON schema or nil if an error occurs
	Load(schemaLocation string) *json_schema.Schema
}

type JsonSchemaFileLoader struct{}

func (*JsonSchemaFileLoader) Load(schemaLocation string) *json_schema.Schema {
	schema, err := json_schema.NewSchema(json_schema.NewReferenceLoader(schemaLocation))
	if err != nil {
		log.Printf("headache configuration warning: cannot load schema, skipping configuration validation. See reason below:\n\t%v\n", err)
		return nil
	}
	return schema
}
