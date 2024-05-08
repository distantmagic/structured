package structured

import (
	"encoding/json"
	"fmt"

	"github.com/distantmagic/paddler/llamacpp"
)

type EntityExtractor struct {
	LlamaCppClient *llamacpp.LlamaCppClient
	MaxRetries     int
}

func (self *EntityExtractor) ExtractFromString(
	responseChannel chan EntityExtractorResult,
	jsonSchema any,
	userInput string,
) {
	defer close(responseChannel)

	entityValidatorBuilder := &EntityValidatorBuilder{}
	entityValidator, err := entityValidatorBuilder.BuildEntityValidator(jsonSchema)

	if err != nil {
		responseChannel <- EntityExtractorResult{
			Error: err,
		}

		return
	}

	llamaCppCompletionResponseChannel := make(chan llamacpp.LlamaCppCompletionToken)

	go self.LlamaCppClient.GenerateCompletion(
		llamaCppCompletionResponseChannel,
		llamacpp.LlamaCppCompletionRequest{
			JsonSchema: jsonSchema,
			NPredict:   100,
			Prompt: fmt.Sprintf(
				`[INST]
				User will provide the phrase. Respond with valid JSON matching
				the schema. Fill the schema with the infromation provided in
				the user phrase. Keep user's input unchanged.

				JSON schema:
				---
				%s
				---
				[/INST]

				User phrase:
				---
				%s
				---`,
				entityValidator.MarshaledJsonSchema,
				userInput,
			),
			Stream: true,
		},
	)

	acc := ""

	for token := range llamaCppCompletionResponseChannel {
		if token.Error != nil {
			responseChannel <- EntityExtractorResult{
				Error: token.Error,
			}

			return
		}

		acc += token.Content
	}

	var unmarshaledLlamaResponse any

	err = json.Unmarshal([]byte(acc), &unmarshaledLlamaResponse)

	if err != nil {
		responseChannel <- EntityExtractorResult{
			Error: err,
		}

		return
	}

	err = entityValidator.CompiledJsonSchema.Validate(unmarshaledLlamaResponse)

	if err != nil {
		responseChannel <- EntityExtractorResult{
			Error: err,
		}

		return
	}

	responseChannel <- EntityExtractorResult{
		Content: acc,
		Entity:  unmarshaledLlamaResponse,
	}
}
