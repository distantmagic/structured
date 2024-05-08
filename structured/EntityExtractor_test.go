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
	MaxRetries: 3,
}

func assertNoErrors(t *testing.T, responseChannel chan EntityExtractorResult) EntityExtractorResult {
	var extractorResult *EntityExtractorResult

	returns := 0

	for result := range responseChannel {
		if result.Error != nil {
			t.Fatal(result.Error)
		}

		extractorResult = &result
		returns += 1
	}

	assert.Equal(t, 1, returns)

	return *extractorResult
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
			"required": []string{"name", "surname", "age"},
		},
		"I am John Doe - living for 40 years and I still like to play chess.",
	)

	extractorResult := assertNoErrors(t, responseChannel)

	var person fixtureNamedPersonWithAge

	err := UnmarshalToStruct(extractorResult, &person)

	assert.Nil(t, err)
	assert.Equal(t, "John", person.Name)
	assert.Equal(t, "Doe", person.Surname)
	assert.Equal(t, 40, person.Age)
}

func TestJsonSchemaNestedConstraints(t *testing.T) {
	responseChannel := make(chan EntityExtractorResult)

	go entityExtractor.ExtractFromString(
		responseChannel,
		map[string]any{
			"description": "Blog post",
			"type":        "object",
			"properties": map[string]any{
				"title": map[string]string{
					"type": "string",
				},
				"content": map[string]string{
					"description": "Blog post content without hashtags and without the title",
					"type":        "string",
				},
				"tags": map[string]any{
					"description": "Hashtags under the blog post without the # symbol",
					"type":        "array",
					"items": map[string]string{
						"type": "string",
					},
				},
			},
			"required": []string{"title", "content", "tags"},
		},
		`
		# How to write a blog post?

		You need to come up with something interesting.

		#blogging #wow #excitement
		`,
	)

	extractorResult := assertNoErrors(t, responseChannel)

	var blogPost fixtureBlogPost

	err := UnmarshalToStruct(extractorResult, &blogPost)

	assert.Nil(t, err)
	assert.Equal(t, "How to write a blog post?", string(blogPost.Title))
	assert.Equal(t, "You need to come up with something interesting.", string(blogPost.Content))
	assert.Equal(t, []string{"blogging", "wow", "excitement"}, blogPost.Tags)
}
