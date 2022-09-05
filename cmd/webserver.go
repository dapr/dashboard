/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dapr/dashboard/pkg/version"

	components "github.com/dapr/dashboard/pkg/components"
	configurations "github.com/dapr/dashboard/pkg/configurations"
	instances "github.com/dapr/dashboard/pkg/instances"
	kube "github.com/dapr/dashboard/pkg/kube"
	dashboard_log "github.com/dapr/dashboard/pkg/log"
	"github.com/dapr/dashboard/pkg/platforms"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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

var upgrader = websocket.Upgrader{} // use default options

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

type DaprVersion struct {
	Version        string `json:"version"`
	RuntimeVersion string `json:"runtimeVersion"`
}

var (
	inst    instances.Instances
	comps   components.Components
	configs configurations.Configurations
)

// RunWebServer starts the web server that serves the Dapr UI dashboard and the API
func RunWebServer(port int, isDockerCompose bool, componentsPath string, configPath string, dockerComposePath string) {
	platform := platforms.Standalone
	kubeClient, daprClient, _ := kube.Clients()
	if kubeClient != nil {
		platform = platforms.Kubernetes
	} else if isDockerCompose {
		platform = platforms.DockerCompose
	}

	inst = instances.NewInstances(platform, kubeClient, dockerComposePath)
	comps = components.NewComponents(platform, daprClient, componentsPath)
	configs = configurations.NewConfigurations(platform, daprClient, configPath)

	r := mux.NewRouter()
	api := r.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/features", getFeaturesHandler).Methods("GET")
	api.HandleFunc("/instances/{scope}", getInstancesHandler).Methods("GET")
	api.HandleFunc("/instances/{scope}/{id}", deleteInstancesHandler).Methods("DELETE")
	api.HandleFunc("/instances/{scope}/{id}", getInstanceHandler).Methods("GET")
	api.HandleFunc("/instances/{scope}/{id}/containers", getContainersHandler).Methods("GET")
	api.HandleFunc("/instances/{scope}/{id}/logstreams/{container}", getLogStreamsHandler)
	api.HandleFunc("/components/{scope}", getComponentsHandler).Methods("GET")
	api.HandleFunc("/components/{scope}/{name}", getComponentHandler).Methods("GET")
	api.HandleFunc("/deploymentconfiguration/{scope}/{id}", getDeploymentConfigurationHandler).Methods("GET")
	api.HandleFunc("/configurations/{scope}", getConfigurationsHandler).Methods("GET")
	api.HandleFunc("/configurations/{scope}/{name}", getConfigurationHandler).Methods("GET")
	api.HandleFunc("/controlplanestatus", getControlPlaneHandler).Methods("GET")
	api.HandleFunc("/metadata/{scope}/{id}", getMetadataHandler).Methods("GET")
	api.HandleFunc("/platform", getPlatformHandler).Methods("GET")
	api.HandleFunc("/scopes", getScopesHandler).Methods("GET")
	api.HandleFunc("/features", getFeaturesHandler).Methods("GET")
	api.HandleFunc("/version", getVersionHandler).Methods("GET")

	spa := spaHandler{staticPath: "web/dist", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("0.0.0.0:%v", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Printf("Dapr Dashboard running on http://localhost:%v\n", port)
	log.Fatal(srv.ListenAndServe())
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the volume of the absolute path and remove it
	volume := filepath.VolumeName(path)
	resourcePath := strings.Replace(path, volume, "", 1)

	// prepend the path with the path to the static directory
	resourcePath = filepath.Join(h.staticPath, resourcePath)

	baseHref, hasCustomBaseHref := os.LookupEnv("SERVER_BASE_HREF")

	// check whether a file exists at the given path
	_, err = os.Stat(resourcePath)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html

		if hasCustomBaseHref {
			generateIndexFile(w, r, baseHref)
		} else {
			http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		}

		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if hasCustomBaseHref && (strings.HasSuffix(path, "index.html") || strings.HasSuffix(path, "/")) {
		generateIndexFile(w, r, baseHref)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	noCache(http.StripPrefix("/", http.FileServer(http.Dir(h.staticPath)))).ServeHTTP(w, r)
}

func getInstancesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scope := vars["scope"]
	if scope == "All" {
		scope = ""
	}
	resp := inst.GetInstances(scope)
	respondWithJSON(w, 200, resp)
}

func getComponentsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scope := vars["scope"]
	if scope == "All" {
		scope = ""
	}
	resp := comps.GetComponents(scope)
	respondWithJSON(w, 200, resp)
}

func getComponentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scope := vars["scope"]
	if scope == "All" {
		scope = ""
	}
	name := vars["name"]
	resp := comps.GetComponent(scope, name)
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
	if inst.CheckPlatform() == platforms.Kubernetes {
		features = append(features, "status")
	}
	respondWithJSON(w, 200, features)
}

func getPlatformHandler(w http.ResponseWriter, r *http.Request) {
	resp := inst.CheckPlatform()
	respondWithPlainString(w, 200, string(resp))
}

func getContainersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scope := vars["scope"]
	if scope == "All" {
		scope = ""
	}
	id := vars["id"]
	containers := inst.GetContainers(scope, id)
	respondWithJSON(w, 200, containers)
}

func getLogStreamsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scope := vars["scope"]
	if scope == "All" {
		scope = ""
	}
	id := vars["id"]
	container := vars["container"]
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer c.Close()
	streams, err := inst.GetLogStream(scope, id, container)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var readStreams []io.Reader
	for _, stream := range streams {
		defer stream.Close()
		readStreams = append(readStreams, stream)
	}
	reader := io.MultiReader(readStreams...)

	lineReader := bufio.NewReader(reader)
	logReader := dashboard_log.NewReader(container, lineReader)
	for {
		logRecord, err := logReader.ReadLog()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if logRecord == nil {
			// Now wait some time before reading more.
			time.Sleep(time.Second)
			continue
		}

		bytes, err := json.Marshal(logRecord)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = c.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			log.Println("fail to write log stream, aborting:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func getDeploymentConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scope := vars["scope"]
	if scope == "All" {
		scope = ""
	}
	id := vars["id"]
	details := inst.GetDeploymentConfiguration(scope, id)
	respondWithPlainString(w, 200, details)
}

func getConfigurationsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scope := vars["scope"]
	if scope == "All" {
		scope = ""
	}
	resp := configs.GetConfigurations(scope)
	respondWithJSON(w, 200, resp)
}

func getConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scope := vars["scope"]
	if scope == "All" {
		scope = ""
	}
	name := vars["name"]
	resp := configs.GetConfiguration(scope, name)
	respondWithJSON(w, 200, resp)
}

func getInstanceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scope := vars["scope"]
	if scope == "All" {
		scope = ""
	}
	id := vars["id"]
	resp := inst.GetInstance(scope, id)
	respondWithJSON(w, 200, resp)
}

func getControlPlaneHandler(w http.ResponseWriter, r *http.Request) {
	resp := inst.GetControlPlaneStatus()
	respondWithJSON(w, 200, resp)
}

func getMetadataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scope := vars["scope"]
	if scope == "All" {
		scope = ""
	}
	id := vars["id"]
	md := inst.GetMetadata(scope, id)
	respondWithJSON(w, 200, md)
}

func deleteInstancesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scope := vars["scope"]
	if scope == "All" {
		scope = ""
	}
	id := vars["id"]
	err := inst.DeleteInstance(scope, id)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	respondWithJSON(w, 200, "")
}

func getScopesHandler(w http.ResponseWriter, r *http.Request) {
	resp := inst.GetScopes()
	respondWithJSON(w, 200, resp)
}

func getVersionHandler(w http.ResponseWriter, r *http.Request) {
	runtimeVersion, err := version.GetRuntimeVersion()
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	resp := DaprVersion{version.GetVersion(), runtimeVersion}
	respondWithJSON(w, 200, resp)
}

func generateIndexFile(w http.ResponseWriter, r *http.Request, baseHref string) {
	path, _ := os.Getwd()
	buf, err := ioutil.ReadFile(filepath.Join(path, "/web/dist/index.html"))
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	file := string(buf)
	file = strings.Replace(file, `<base href="/">`, fmt.Sprintf(`<base href="%s">`, baseHref), 1)
	respondWithHtml(w, 200, file)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithPlainString(w http.ResponseWriter, code int, payload string) {
	w.WriteHeader(code)
	_, err := w.Write([]byte(payload))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func respondWithHtml(w http.ResponseWriter, code int, payload string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
	_, err := w.Write(([]byte(payload)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
