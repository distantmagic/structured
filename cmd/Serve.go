package cmd

import (
	"net/http"

	"github.com/distantmagic/paddler/llamacpp"
	"github.com/distantmagic/paddler/netcfg"
	"github.com/distantmagic/structured/httpserver"
	"github.com/distantmagic/structured/structured"
	"github.com/hashicorp/go-hclog"
	"github.com/urfave/cli/v2"
)

type Serve struct {
	Logger                hclog.Logger
	HttpAddress           *netcfg.HttpAddressConfiguration
	LlamaCppConfiguration *llamacpp.LlamaCppConfiguration
}

func (self *Serve) Action(cliContext *cli.Context) error {
	entityExtractor := &structured.EntityExtractor{
		LlamaCppClient: &llamacpp.LlamaCppClient{
			HttpClient:            http.DefaultClient,
			LlamaCppConfiguration: self.LlamaCppConfiguration,
		},
	}

	server := httpserver.Server{
		HttpAddress: self.HttpAddress,
		Logger:      self.Logger,
		RespondToExtractEntity: &httpserver.RespondToExtractEntity{
			EntityExtractor: entityExtractor,
			Logger:          self.Logger.Named("RespondToExtractEntity"),
		},
	}

	return server.Serve()
}
