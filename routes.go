package logsauce

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func init() {

}

func RegisterRoutes(config Configuration) {
	r := mux.NewRouter()
	port := strconv.Itoa(config.ListenPort)

	if config.ServerMode {
		r.HandleFunc("/", serverHandler).Methods("POST")
		//go http.ListenAndServeTLS(":"+string(config.ListenPort), config.ServerConfiguration.ServerCertificate, config.ServerConfiguration.ServerCertificateKey, r)
		go log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))

	} else if config.Relaymode {
		r.HandleFunc("/", relayHandler).Methods("POST")
		//go http.ListenAndServeTLS(":"+string(config.ListenPort), config.ServerConfiguration.ServerCertificate, config.ServerConfiguration.ServerCertificateKey, r)
		go log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))
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
	newLogline.Timestamp = time.Now().Unix()

	err = insertLogline(newLogline)
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
	hosts, err := getAllHosts()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return currentHost, err
	}

	for i := 0; i < len(hosts); i++ {
		if hosts[i].Token == token[0] {
			currentHost = hosts[i]
			break
		}

		if i == len(hosts) {
			log.Println("Unknown host token")
			w.WriteHeader(http.StatusInternalServerError)
			return currentHost, errors.New("Unknown host token")
		}
	}

	return currentHost, nil
}
