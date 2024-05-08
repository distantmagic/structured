package structured

import "github.com/distantmagic/structured/unmarshal"

type fixtureBlogPost struct {
	Title   unmarshal.TrimmedString `json:"title"`
	Content unmarshal.TrimmedString `json:"content"`
	Tags    []string                `json:"tags"`
}
