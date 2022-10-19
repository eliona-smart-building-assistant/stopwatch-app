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
	"stopwatch/stopwtch"
	"time"

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
	swMng       stopwtch.StopwatchManager
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

	swMng = *stopwtch.NewStopwatchManager(stopwatchCallbackHandler)

	go eliona.ListenHeapEvents(ir, apiData)
	go stopwatchEventCatcher(apiData)
}

// start / stop stopwatches
func stopwatchEventCatcher(apiData <-chan api.Data) {
	for data := range apiData {
		log.Debug(MODULE, "data from datalistener: %v", data)

		asset := getStopwatchByAssetId(data.AssetId)
		if asset != nil {
			start := data.Data[eliona.ATTRIBUTE_START]
			stop := data.Data[eliona.ATTRIBUTE_STOP]
			if stop != nil && stop.(float64) >= 1 {
				lastTime := swMng.Stop(asset.GetId())
				if lastTime.Seconds() > 0 {
					log.Info(MODULE, "timer stoppend %d @ %d s", asset.GetId(), lastTime.Seconds())
					eliona.UpdateLastTime(data.AssetId, lastTime.Seconds())
				}
			} else if start != nil && start.(float64) >= 1 {
				log.Info(MODULE, "timer started %d", asset.GetId())
				swMng.Start(asset.GetId())
			}
		}
	}
	log.Warn(MODULE, "eventcatcher exited")
}

// stopwatch ticks
func stopwatchCallbackHandler(assetId int32, t time.Duration) {
	log.Debug(MODULE, "stopwatch callback called %d,%v", assetId, t)

	eliona.UpdateTime(assetId, t.Seconds())
}

func getStopwatchByAssetId(assetId int32) *api.Asset {
	for _, sw := range stopwatches {
		if sw.GetId() == assetId {
			return &sw
		}
	}
	return nil
}

// listenApiRequests starts an API server and listen for API requests
// The API endpoints are defined in the openapi.yaml file
func listenApiRequests() {
	err := http.ListenAndServe(":"+common.Getenv("API_SERVER_PORT", "3000"), apiserver.NewRouter(
		apiserver.NewConfigurationApiController(apiservices.NewConfigurationApiService()),
	))
	log.Fatal("Hailo", "Error in API Server: %v", err)
}
