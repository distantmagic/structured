package structured

import (
	"bytes"
	"encoding/json"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

type EntityValidatorBuilder struct {
}

func (self *EntityValidatorBuilder) BuildEntityValidator(
	jsonSchema any,
) (*EntityValidator, error) {
	marshaledJsonSchema, err := json.Marshal(jsonSchema)

	if err != nil {
		return nil, err
	}

	jsonSchemaCompiler := jsonschema.NewCompiler()

	err = jsonSchemaCompiler.AddResource(
		"schema.json",
		bytes.NewReader(marshaledJsonSchema),
	)

	if err != nil {
		return nil, err
	}

	schema, err := jsonSchemaCompiler.Compile("schema.json")

	if err != nil {
		return nil, err
	}

	entityValidator := &EntityValidator{
		CompiledJsonSchema:  schema,
		MarshaledJsonSchema: marshaledJsonSchema,
	}

	return entityValidator, nil
}
