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

func RegisterRoutes(listenPort int, relayMode, serverMode bool) {
	r := mux.NewRouter()
	port := strconv.Itoa(listenPort)

	if serverMode {
		r.HandleFunc("/logs", serverHandler).Methods("POST")
		r.HandleFunc("/logs/{hostname}/{filepath}/{starttime}/{endtime}", getLogsHandler).Methods("GET")
		r.HandleFunc("/compute", addDefaultHeaders(computeHandler)).Methods("POST")
		//go http.ListenAndServeTLS(":"+string(config.ListenPort), config.ServerConfiguration.ServerCertificate, config.ServerConfiguration.ServerCertificateKey, r)
		go http.ListenAndServe("0.0.0.0:"+port, r)
		log.Println("Server mode is active")
	} else if relayMode {
		r.HandleFunc("/", relayHandler).Methods("POST")
		//go http.ListenAndServeTLS(":"+string(config.ListenPort), config.ServerConfiguration.ServerCertificate, config.ServerConfiguration.ServerCertificateKey, r)
		go http.ListenAndServe("0.0.0.0:"+port, r)
		log.Println("Relay mode is active")
	}
}

/*
Restricted route. needs login or a client token
*/
func computeHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var computeRequest ComputeRequest
	err = json.Unmarshal(body, &computeRequest)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("%#v", computeRequest)

	if computeRequest.TimeStart <= 0 && computeRequest.TimeEnd <= 0 {
		log.Println("Bad timestart or end")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch computeRequest.Operation {
	case "dist":
		{
			resolution, err := strconv.Atoi(computeRequest.Parameter1)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			fieldNames := strings.Split(computeRequest.Parameter2, ",")
			log.Println(fieldNames, resolution)

			var responses []ComputeResponse
			for i := 0; i < len(fieldNames); i++ {
				newResponses := calcDistributionOverTime(computeRequest.LogtypeName, computeRequest.Host, computeRequest.Filename, resolution, fieldNames[i], computeRequest.TimeStart, computeRequest.TimeEnd)
				for i := 0; i < len(newResponses); i++ {
					responses = append(responses, newResponses[i])
				}
			}

			b, err := json.Marshal(&responses)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Write(b)
			return
		}

	case "uniq":
		{
			resolution, err := strconv.Atoi(computeRequest.Parameter1)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			fieldNames := strings.Split(computeRequest.Parameter2, ",")
			log.Println(fieldNames, resolution)

			var responses []ComputeResponse
			for i := 0; i < len(fieldNames); i++ {
				newResponses := countUniqueOverTime(computeRequest.LogtypeName, computeRequest.Host, computeRequest.Filename, resolution, fieldNames[i], computeRequest.TimeStart, computeRequest.TimeEnd)
				for i := 0; i < len(newResponses); i++ {
					responses = append(responses, newResponses[i])
				}
			}

			b, err := json.Marshal(&responses)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Write(b)
			return
		}
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

func addDefaultHeaders(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		fn(w, r)
	}
}
