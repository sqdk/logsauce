package logsauce

import (
	"github.com/gorilla/mux"
	"net/http"
)

func init() {

}

func RegisterRoutes(config Configuration) {
	r := mux.NewRouter()

	if config.ServerMode {
		r.HandleFunc("/", serverHandler)
		http.ListenAndServeTLS(":"+string(config.ListenPort), config.ServerConfiguration.ServerCertificate, config.ServerConfiguration.ServerCertificateKey, r)
	} else if config.Relaymode {
		r.HandleFunc("/", relayHandler)
		http.ListenAndServeTLS(":"+string(config.ListenPort), config.ServerConfiguration.ServerCertificate, config.ServerConfiguration.ServerCertificateKey, r)
	}

}

func serverHandler(w http.ResponseWriter, r *http.Request) {

}

func relayHandler(w http.ResponseWriter, r *http.Request) {

}
