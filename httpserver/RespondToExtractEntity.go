package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/distantmagic/structured/structured"
	"github.com/hashicorp/go-hclog"
)

type RespondToExtractEntity struct {
	EntityExtractor *structured.EntityExtractor
	Logger          hclog.Logger
}

func (self *RespondToExtractEntity) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	var extractEntityRequest ExtractEntityRequest

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&extractEntityRequest)

	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)

		return
	}

	self.Logger.Debug(
		"received",
		"schema", extractEntityRequest.Schema,
		"data", extractEntityRequest.Data,
	)

	responseChannel := make(chan structured.EntityExtractorResult)

	go self.EntityExtractor.ExtractFromString(
		responseChannel,
		extractEntityRequest.Schema,
		extractEntityRequest.Data,
	)

	extractorResult := <-responseChannel

	if extractorResult.Error != nil {
		http.Error(response, extractorResult.Error.Error(), http.StatusBadRequest)

		return
	}

	output := ExtractEnityResponse{
		Entity: extractorResult.Entity,
	}

	response.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(response).Encode(output)

	if err != nil {
		self.Logger.Error("error", err)
	}
}
