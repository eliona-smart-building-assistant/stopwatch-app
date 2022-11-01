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
	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/client"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

const (
	MODULE = "eliClient"
	// asset definitons
	ASSET_TYPE_STOPWATCH = "Stopwatch"
)

type StopwatchValueData struct {
	CurrentTime *int32 `json:"current_time,omitempty"`
	LastTime    *int32 `json:"last_time,omitempty"`
}

func GetStopwatches() ([]api.Asset, error) {
	var stopwatches []api.Asset

	assets, resp, err := client.NewClient().AssetsApi.GetAssets(client.AuthenticationContext()).Execute()

	log.Debug(MODULE, "get stopwatches response: %v", resp)

	for _, asset := range assets {
		if asset.AssetType == ASSET_TYPE_STOPWATCH {
			stopwatches = append(stopwatches, asset)
		}
	}
	return stopwatches, err
}
