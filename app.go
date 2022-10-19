//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"net/http"
	"stopwatch/apiserver"
	"stopwatch/apiservices"
	"stopwatch/eliona"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

const (
	MODULE          = "app"
	API_DATA_BUFFER = 20000
)

var (
	stopwatches []api.Asset
	ir          chan bool
)

// doAnything is the main app function which is called periodically
func actualizeStopwatches() {
	var err error

	stopwatches, err = eliona.GetStopwatches()
	log.Debug(MODULE, "got stopwatches %v with err %v", stopwatches, err)
}

func setupApp() {
	log.Info(MODULE, "setup application")

	var (
		apiData chan api.Data
	)

	ir = make(chan bool)
	apiData = make(chan api.Data, API_DATA_BUFFER)

	actualizeStopwatches()

	go eliona.ListenHeapEvents(ir, apiData)
	go stopwatchEventCatcher(apiData)
}

func stopwatchEventCatcher(apiData <-chan api.Data) {
	for data := range apiData {
		log.Debug(MODULE, "data from datalistener: %v", data)
	}
	log.Warn(MODULE, "eventcatcher exited")
}

// listenApiRequests starts an API server and listen for API requests
// The API endpoints are defined in the openapi.yaml file
func listenApiRequests() {
	err := http.ListenAndServe(":"+common.Getenv("API_SERVER_PORT", "3000"), apiserver.NewRouter(
		apiserver.NewConfigurationApiController(apiservices.NewConfigurationApiService()),
	))
	log.Fatal("Hailo", "Error in API Server: %v", err)
}
