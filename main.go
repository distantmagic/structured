package main

import (
	"os"

	"github.com/distantmagic/paddler/llamacpp"
	"github.com/distantmagic/paddler/netcfg"
	"github.com/distantmagic/structured/cmd"
	"github.com/hashicorp/go-hclog"
	"github.com/urfave/cli/v2"
)

func main() {
	logger := hclog.Default()

	serve := &cmd.Serve{
		Logger:      logger.Named("server"),
		HttpAddress: &netcfg.HttpAddressConfiguration{},
		LlamaCppConfiguration: &llamacpp.LlamaCppConfiguration{
			HttpAddress: &netcfg.HttpAddressConfiguration{},
		},
	}

	app := &cli.App{
		Name:  "structured",
		Usage: "Extract structured data from unstructured text",
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "start http server for API",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "host",
						Value:       "127.0.0.1",
						Destination: &serve.HttpAddress.Host,
					},
					&cli.UintFlag{
						Name:        "port",
						Value:       8087,
						Destination: &serve.HttpAddress.Port,
					},
					&cli.StringFlag{
						Name:        "scheme",
						Value:       "http",
						Destination: &serve.HttpAddress.Scheme,
					},
					&cli.StringFlag{
						Name:        "llamacpp-host",
						Value:       "127.0.0.1",
						Destination: &serve.LlamaCppConfiguration.HttpAddress.Host,
					},
					&cli.UintFlag{
						Name:        "llamacpp-port",
						Value:       8081,
						Destination: &serve.LlamaCppConfiguration.HttpAddress.Port,
					},
					&cli.StringFlag{
						Name:        "llamacpp-scheme",
						Value:       "http",
						Destination: &serve.LlamaCppConfiguration.HttpAddress.Scheme,
					},
				},
				Action: serve.Action,
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		panic(err)
	}
}
