package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/dapr/dashboard/pkg/instances"
)

const (
	dir = "./web/dist"
)

var epoch = time.Unix(0, 0).Format(time.RFC1123)

var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

var inst = instances.NewInstances()

// RunWebServer starts the web server that serves the Dapr UI dashboard and the API
func RunWebServer() {
	r := mux.NewRouter()
	r.HandleFunc("/api/instances", getInstancesHandler)
	r.HandleFunc("/api/instances/{id}", deleteInstancesHandler).Methods("DELETE")
	r.HandleFunc("/api/instances/{id}/logs", getLogsHandler)
	r.PathPrefix("/").Handler(noCache(http.StripPrefix("/", http.FileServer(http.Dir(dir)))))

	fmt.Println("Dapr Dashboard running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getInstancesHandler(w http.ResponseWriter, r *http.Request) {
	resp := inst.Get()
	respondWithJSON(w, 200, resp)
}

func getLogsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	logs := inst.Logs(id)
	respondWithPlainString(w, 200, logs)
}

func deleteInstancesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := inst.Delete(id)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	respondWithJSON(w, 200, "")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithPlainString(w http.ResponseWriter, code int, payload string) {
	w.WriteHeader(code)
	w.Write([]byte(payload))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func noCache(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Delete any ETag headers that may have been set
		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
