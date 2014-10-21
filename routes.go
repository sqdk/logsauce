package logsauce

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func init() {

}

func RegisterRoutes(config Configuration) {
	log.Println("Server mode is active")
	r := mux.NewRouter()
	log.Println("Server mode is active")
	port := strconv.Itoa(config.ListenPort)
	log.Println(port)

	if config.ServerMode {
		log.Println("Server mode is active")

		r.HandleFunc("/", serverHandler).Methods("POST")

		//go http.ListenAndServeTLS(":"+string(config.ListenPort), config.ServerConfiguration.ServerCertificate, config.ServerConfiguration.ServerCertificateKey, r)
		go log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))
		log.Println("Server mode is active")
	} else if config.Relaymode {
		r.HandleFunc("/", relayHandler).Methods("POST")
		//go http.ListenAndServeTLS(":"+string(config.ListenPort), config.ServerConfiguration.ServerCertificate, config.ServerConfiguration.ServerCertificateKey, r)
		go log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))
		log.Println("Relay mode is active")
	}
}

func serverHandler(w http.ResponseWriter, r *http.Request) {
	host, err := verifyToken(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	requestBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var newLogline LogLine
	err = json.Unmarshal(requestBytes, &newLogline)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newLogline.HostId = host.Id

	insertLogline(newLogline)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func relayHandler(w http.ResponseWriter, r *http.Request) {

}

func verifyToken(w http.ResponseWriter, r *http.Request) (Host, error) {
	var currentHost Host

	token := r.Header["Token"] //Check if header is present
	log.Println(r.Header)
	if len(token) == 0 {
		w.WriteHeader(http.StatusForbidden)
		return currentHost, errors.New("No token")
	}

	//Verify token
	host, err := getHostWithToken(token[0])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return Host{}, err
	}

	currentHost = host

	return currentHost, nil
}
