package core

import (
	"encoding/json"
	"github.com/fbiville/headache/fs"
	jsonsch "github.com/xeipuuv/gojsonschema"
	"log"
)

type ConfigurationLoader struct {
	Reader fs.FileReader
}

func (cl *ConfigurationLoader) ReadConfiguration(configFile *string) (*Configuration, error) {
	err := cl.validateConfiguration(configFile)
	if err != nil {
		return nil, err
	}

	payload, err := cl.Reader.Read(*configFile)
	if err != nil {
		return nil, err
	}
	configuration, err := cl.UnmarshallConfiguration(payload)
	if err != nil {
		return nil, err
	}
	configuration.Path = configFile
	return configuration, err
}

func (cl *ConfigurationLoader) UnmarshallConfiguration(configurationPayload []byte) (*Configuration, error) {
	result := Configuration{}
	err := json.Unmarshal(configurationPayload, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (cl *ConfigurationLoader) validateConfiguration(configFile *string) error {
	schema := loadSchema()
	if schema == nil {
		return nil
	}
	jsonSchemaValidator := JsonSchemaValidator{
		Schema:     schema,
		FileReader: cl.Reader,
	}
	return jsonSchemaValidator.Validate("file://" + *configFile)
}

func loadSchema() *jsonsch.Schema {
	schema, err := jsonsch.NewSchema(jsonsch.NewReferenceLoader("https://fbiville.github.io/headache/schema.json"))
	if err != nil {
		log.Printf("headache configuration warning: cannot load schema, skipping configuration validation. See reason below:\n\t%v\n", err)
		return nil
	}
	return schema
}
