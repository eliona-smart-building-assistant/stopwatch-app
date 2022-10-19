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

package conn

import (
	"crypto/tls"
	"flag"
	"net/url"
	"strings"
	"time"

	"github.com/eliona-smart-building-assistant/go-utils/log"
	"github.com/gorilla/websocket"
)

const MODULE = "websock"

type WebsocketClient struct {
	url        url.URL
	conn       *websocket.Conn
	IgnoreCert bool
}

func createUrl(websocketUrl string) url.URL {
	var (
		path  string
		query string
	)
	websocketUrl = strings.ReplaceAll(websocketUrl, "ws://", "")
	websocketUrl = strings.ReplaceAll(websocketUrl, "wss://", "")
	websocketUrl = strings.ReplaceAll(websocketUrl, "http://", "")
	websocketUrl = strings.ReplaceAll(websocketUrl, "https://", "")
	domain := strings.Split(websocketUrl, `/`)[0]
	pathQuery := strings.Split(strings.ReplaceAll(websocketUrl, domain, ""), `?`)
	if len(pathQuery) > 0 {
		path = pathQuery[0]
	}
	if len(pathQuery) > 1 {
		query = pathQuery[1]
	}

	var addr = flag.String("addr", domain, "wss address")

	return url.URL{Scheme: "wss", Host: *addr, Path: path, RawQuery: query}
}

func NewWebsocketClient(websocketUrl string, skipVerifyCertificate bool) *WebsocketClient {
	return &WebsocketClient{
		url:        createUrl(websocketUrl),
		IgnoreCert: skipVerifyCertificate,
	}
}

func (ws *WebsocketClient) ServeForever(rxChannel chan<- []byte, interrupt <-chan bool) {
	var err error

	log.Info(MODULE, "connecting to %s", ws.url.String())

	tlsConfig := tls.Config{InsecureSkipVerify: ws.IgnoreCert}
	websocket.DefaultDialer.TLSClientConfig = &tlsConfig

	ws.conn, _, err = websocket.DefaultDialer.Dial(ws.url.String(), nil)
	if err != nil {
		log.Error(MODULE, "wss dial", err)
	}

	defer ws.conn.Close()

	readerClosed := make(chan bool)
	go ws.readerLoop(rxChannel, readerClosed)

	for {
		select {
		case <-readerClosed:
			log.Error(MODULE, "closed reader")
			return

		case <-interrupt:
			log.Info(MODULE, "interrupted")
			err = ws.conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Info(MODULE, "write close", err)
			}
			// wait for wssReaderLoop
			select {
			case <-readerClosed:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func (ws *WebsocketClient) readerLoop(rxChannel chan<- []byte, closed chan<- bool) {
	defer close(rxChannel)
	defer close(closed)
	for {
		_, message, err := ws.conn.ReadMessage()
		if err != nil {
			log.Error(MODULE, "read message", err)
			return
		}
		rxChannel <- message
	}
}

func (ws *WebsocketClient) Send(message []byte) error {
	err := ws.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Error(MODULE, "send message", err)
	}
	return err
}
