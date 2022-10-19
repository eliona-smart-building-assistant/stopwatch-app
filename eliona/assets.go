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

package eliona

import (
	"context"
	"encoding/json"

	"stopwatch/conn"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/client"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

const (
	MODULE             = "eliClient"
	RX_BUFFER          = 20000
	WEBSOCKET_ENDPOINT = "/data-listener"
)

func GetStopwatches() ([]api.Asset, error) {
	var stopwatches []api.Asset

	assets, resp, err := client.NewClient().AssetsApi.GetAssets(context.Background()).Execute()

	log.Debug(MODULE, "get stopwatches response: %v", resp)

	for _, asset := range assets {
		if asset.AssetType == "Stopwatch" {
			stopwatches = append(stopwatches, asset)
		}
	}
	return stopwatches, err
}

func ListenHeapEvents(ir <-chan bool, rxApiData chan<- api.Data) {
	var (
		rx chan []byte
	)

	rx = make(chan []byte, RX_BUFFER)

	ws := conn.NewWebsocketClient(client.ApiEndpointString()+WEBSOCKET_ENDPOINT, false)
	go ws.ServeForever(rx, ir)
	for data := range rx {
		var apiData api.Data
		err := json.Unmarshal(data, &apiData)
		if err == nil {
			rxApiData <- apiData
		} else {
			log.Warn(MODULE, "cannot unmarshal ws data", err)
		}
	}
	log.Warn(MODULE, "rx channel from ws closed")
}
