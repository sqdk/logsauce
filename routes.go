package logsauce

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func init() {

}

func RegisterRoutes(config Configuration) {
	r := mux.NewRouter()
	port := strconv.Itoa(config.ListenPort)

	if config.ServerMode {
		r.HandleFunc("/", serverHandler).Methods("POST")
		r.HandleFunc("/logs/{hostname}/{filepath}/{starttime}/{endtime}", getLogsHandler).Methods("GET")

		//go http.ListenAndServeTLS(":"+string(config.ListenPort), config.ServerConfiguration.ServerCertificate, config.ServerConfiguration.ServerCertificateKey, r)
		go http.ListenAndServe("127.0.0.1:"+port, r)
		log.Println("Server mode is active")
	} else if config.Relaymode {
		r.HandleFunc("/", relayHandler).Methods("POST")
		//go http.ListenAndServeTLS(":"+string(config.ListenPort), config.ServerConfiguration.ServerCertificate, config.ServerConfiguration.ServerCertificateKey, r)
		go http.ListenAndServe("0.0.0.0:"+port, r)
		log.Println("Relay mode is active")
	}
}

/*
Restricted route. needs login or a client token
Slashes in filepath is replaced with +.
*/
func getLogsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["hostname"] == "" || vars["filepath"] == "" || vars["starttime"] == "" || vars["endtime"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	startTime, err := strconv.Atoi(vars["starttime"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	endTime, err := strconv.Atoi(vars["endtime"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	parsedFilepath := strings.Replace(vars["filepath"], "+", "/", -1)
	log.Println(parsedFilepath)
	logs, err := getLoglinesForPeriodForHostnameAndFilepath(vars["hostname"], parsedFilepath, int64(startTime), int64(endTime))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := json.Marshal(&logs)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
	return
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
	if len(token) == 0 {
		w.WriteHeader(http.StatusForbidden)
		return Host{}, errors.New("No token")
	}

	//Verify token
	host, err := getHostWithToken(token[0])
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return Host{}, errors.New("Unknown token")
	}

	currentHost = host

	return currentHost, nil
}
