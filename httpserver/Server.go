package httpserver

import (
	"net/http"

	"github.com/distantmagic/paddler/netcfg"
	"github.com/hashicorp/go-hclog"
)

type Server struct {
	Logger                 hclog.Logger
	HttpAddress            *netcfg.HttpAddressConfiguration
	RespondToExtractEntity *RespondToExtractEntity
}

func (self *Server) Serve() error {
	self.Logger.Info("Starting HTTP server", "address", self.HttpAddress.GetHostWithPort())

	mux := http.NewServeMux()
	mux.Handle("/extract/entity", self.RespondToExtractEntity)

	return http.ListenAndServe(self.HttpAddress.GetHostWithPort(), mux)
}
