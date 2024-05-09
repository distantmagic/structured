package httpserver

type ExtractEntityRequest struct {
	Schema any    `json:"schema"`
	Data   string `json:"data"`
}
