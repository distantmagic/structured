package structured

import (
	"net/http"
	"testing"

	"github.com/distantmagic/paddler/llamacpp"
	"github.com/distantmagic/paddler/netcfg"
	"github.com/stretchr/testify/assert"
)

var entityExtractor *EntityExtractor = &EntityExtractor{
	LlamaCppClient: &llamacpp.LlamaCppClient{
		HttpClient: http.DefaultClient,
		LlamaCppConfiguration: &llamacpp.LlamaCppConfiguration{
			HttpAddress: &netcfg.HttpAddressConfiguration{
				Host:   "127.0.0.1",
				Port:   8081,
				Scheme: "http",
			},
		},
	},
}

func TestJsonSchemaConstrainedCompletionsAreGenerated(t *testing.T) {
	responseChannel := make(chan EntityExtractorResult)

	go entityExtractor.ExtractFromString(
		responseChannel,
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]string{
					"type": "string",
				},
				"surname": map[string]string{
					"type": "string",
				},
				"age": map[string]string{
					"description": "Age in years.",
					"type":        "integer",
				},
			},
		},
		"I am John Doe - living for 40 years and I still like to play chess.",
	)

	returns := 0

	for result := range responseChannel {
		if result.Error != nil {
			t.Fatal(result.Error)
		}

		returns += 1
	}

	assert.Equal(t, 1, returns)
}
