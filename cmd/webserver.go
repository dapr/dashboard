package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	components "github.com/dapr/dashboard/pkg/components"
	configurations "github.com/dapr/dashboard/pkg/configurations"
	instances "github.com/dapr/dashboard/pkg/instances"
	kube "github.com/dapr/dashboard/pkg/kube"
	status "github.com/dapr/dashboard/pkg/status"
	"github.com/gorilla/mux"
)

const (
	dir  = "./web/dist"
	port = 8080
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

var inst instances.Instances
var comps components.Components
var configs configurations.Configurations
var stats status.Status

// RunWebServer starts the web server that serves the Dapr UI dashboard and the API
func RunWebServer() {
	kubeClient, daprClient, _ := kube.Clients()
	inst = instances.NewInstances(kubeClient)
	comps = components.NewComponents(daprClient)
	configs = configurations.NewConfigurations(daprClient)
	stats = status.NewStatus(kubeClient)

	r := mux.NewRouter()
	r.HandleFunc("/api/features", getFeaturesHandler)
	r.HandleFunc("/api/instances", getInstancesHandler)
	r.HandleFunc("/api/instances/{id}", deleteInstancesHandler).Methods("DELETE")
	r.HandleFunc("/api/instances/{id}", getInstanceHandler).Methods("GET")
	r.HandleFunc("/api/instances/{id}/logs", getLogsHandler)
	r.HandleFunc("/api/components", getComponentsHandler)
	r.HandleFunc("/api/components/status", getComponentsStatusHandler)
	r.HandleFunc("/api/configuration/{id}", getConfigurationHandler)
	r.HandleFunc("/api/daprconfig", getDaprConfigHandler)
	r.HandleFunc("/api/environments", getEnvironmentsHandler)
	r.HandleFunc("/api/controlplanestatus", getControlPlaneHandler)

	fs := http.FileServer(http.Dir(dir))
	r.HandleFunc("/{rest:.*}", angularHandler)
	r.PathPrefix("/").Handler(noCache(http.StripPrefix("/", fs)))

	fmt.Println(fmt.Sprintf("Dapr Dashboard running on http://localhost:%v", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}

func angularHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/dist/index.html")
}

func getInstancesHandler(w http.ResponseWriter, r *http.Request) {
	resp := inst.Get()
	respondWithJSON(w, 200, resp)
}

func getComponentsHandler(w http.ResponseWriter, r *http.Request) {
	resp := comps.Get()
	respondWithJSON(w, 200, resp)
}

func getComponentsStatusHandler(w http.ResponseWriter, r *http.Request) {
	resp := comps.GetStatus()
	respondWithJSON(w, 200, resp)
}

func getFeaturesHandler(w http.ResponseWriter, r *http.Request) {
	features := []string{}
	if comps.Supported() {
		features = append(features, "components")
	}
	if configs.Supported() {
		features = append(features, "configurations")
	}
	if stats.Supported() {
		features = append(features, "status")
	}
	respondWithJSON(w, 200, features)
}

func getEnvironmentsHandler(w http.ResponseWriter, r *http.Request) {
	resp := inst.CheckSupportedEnvironments()
	respondWithJSON(w, 200, resp)
}

func getLogsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	logs := inst.Logs(id)
	respondWithPlainString(w, 200, logs)
}

func getConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	details := inst.Configuration(id)
	respondWithPlainString(w, 200, details)
}

func getDaprConfigHandler(w http.ResponseWriter, r *http.Request) {
	resp := configs.Get()
	respondWithJSON(w, 200, resp)
}

func getInstanceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	resp := inst.GetInstance(id)
	respondWithJSON(w, 200, resp)
}

func getControlPlaneHandler(w http.ResponseWriter, r *http.Request) {
	resp := stats.Get()
	respondWithJSON(w, 200, resp)
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
