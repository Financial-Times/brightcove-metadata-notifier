package main

import (
	"net/http"
	"os"
	"strconv"

	"fmt"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
)

type metadataMapper struct {
	mappings map[string]term
	config   *notifierConfig
	client   *http.Client
}

type notifierConfig struct {
	mappingURL              string
	cmsMetadataNotifierAddr string
	cmsMetadataNotifierHost string
	cmsMetadataNotifierAuth string
	port                    int
}

type healthcheck struct {
	config *notifierConfig
	client *http.Client
}

func main() {
	cliApp := cli.App("brightcove-metadata-notifier", "Gets the video model, maps the video tags to TME metadata from which it builds the raw metadata msg and posts it to cms-metadata-notifier")
	cmsMetadataNotifierAddr := cliApp.String(cli.StringOpt{
		Name:   "cms-metadata-notifier-address",
		Value:  "",
		Desc:   "Address of the cmsMetadataNotifier",
		EnvVar: "CMS_METADATA_NOTIFIER_ADDR",
	})
	cmsMetadataNotifierHost := cliApp.String(cli.StringOpt{
		Name:   "cms-metadata-notifier-host-header",
		Value:  "cms-metadata-notifier",
		Desc:   "Host header to use to connect to cms-metadata-notifier",
		EnvVar: "CMS_METADATA_NOTIFIER_HOST_HEADER",
	})
	cmsMetadataNotifierAuth := cliApp.String(cli.StringOpt{
		Name:   "cms-metadata-notifier-auth",
		Value:  "",
		Desc:   "cms-metadata-notifier authorization header",
		EnvVar: "CMS_METADATA_NOTIFIER_AUTH",
	})
	mappingURL := cliApp.String(cli.StringOpt{
		Name:   "mapping-url",
		Value:  "",
		Desc:   "URL of the metadata mapping spreadsheet in Bertha",
		EnvVar: "MAPPING_URL",
	})
	port := cliApp.Int(cli.IntOpt{
		Name:   "port",
		Value:  8080,
		Desc:   "Listening port of this service",
		EnvVar: "PORT",
	})

	cliApp.Action = func() {
		initLogs(os.Stdout, os.Stdout, os.Stderr)
		if *mappingURL == "" {
			errorLogger.Panicf("Please provide a valid URL")
		}
		nConfig := &notifierConfig{
			mappingURL:              *mappingURL,
			cmsMetadataNotifierAddr: *cmsMetadataNotifierAddr,
			cmsMetadataNotifierHost: *cmsMetadataNotifierHost,
			cmsMetadataNotifierAuth: *cmsMetadataNotifierAuth,
			port: *port,
		}
		infoLogger.Printf("%v", nConfig.prettyPrint())
		httpClient := &http.Client{}

		metadataMapper := buildMetadataMapper(nConfig, httpClient)
		infoLogger.Printf("%v", metadataMapper.prettyPrintMappings())
		hc := healthcheck{nConfig, httpClient}

		listen(metadataMapper, hc)
	}
	err := cliApp.Run(os.Args)
	if err != nil {
		println(err)
	}
}

func buildMetadataMapper(config *notifierConfig, client *http.Client) metadataMapper {
	return metadataMapper{
		mappings: fetchMappings(config.mappingURL),
		config:   config,
		client:   client,
	}
}

func listen(mm metadataMapper, hc healthcheck) {
	r := mux.NewRouter()
	r.HandleFunc("/notify", mm.handleNotification).Methods("POST")
	r.HandleFunc("/__health", hc.health()).Methods("GET")
	r.HandleFunc("/__gtg", hc.gtg).Methods("GET")

	http.Handle("/", r)
	infoLogger.Printf("Starting to listen on port [%d]", mm.config.port)
	err := http.ListenAndServe(":"+strconv.Itoa(mm.config.port), nil)
	if err != nil {
		errorLogger.Panicf("Couldn't set up HTTP listener: [%+v]\n", err)
	}
}

func (nc notifierConfig) prettyPrint() string {
	authSet := "empty"
	if nc.cmsMetadataNotifierAuth != "" {
		authSet = "set, not empty"
	}
	return fmt.Sprintf("\n\t\tmappingURL: [%s]\n\t\tcmsMetadataNotifierAddr: [%s]\n\t\tcmsMetadataNotifierHost: [%s]\n\t\tport: [%s]\n\t\tcmsMetadataNotifierAuth: [%s]\n\t", nc.mappingURL, nc.cmsMetadataNotifierAddr, nc.cmsMetadataNotifierHost, nc.port, authSet)
}
