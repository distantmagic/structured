package structured

import (
	"bytes"
	"encoding/json"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

type EntityJsonSchemaValidatorBuilder struct {
}

func (self *EntityJsonSchemaValidatorBuilder) BuildEntityJsonSchemaValidator(
	jsonSchema any,
) (*EntityJsonSchemaValidator, error) {
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

	entityValidator := &EntityJsonSchemaValidator{
		CompiledJsonSchema:  schema,
		MarshaledJsonSchema: marshaledJsonSchema,
	}

	return entityValidator, nil
}
