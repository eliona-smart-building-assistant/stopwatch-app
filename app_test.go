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
	"log"
	"testing"
	"time"

	"stopwatch/stopwtch"
)

var count int

func TestStopwatchManager(t *testing.T) {
	swMng := stopwtch.NewStopwatchManager(stopwatchCallback)
	defer swMng.StopAll()

	swMng.Start(1)
	time.Sleep(2 * time.Second)
	if count > 3 || count < 1 {
		t.Error("timer not at 10 after 10 seconds: ", count)
	}
	t.Log("stop timer 1")
	swMng.Stop(1)
	t.Log("wait stop all")
}

func stopwatchCallback(id int32, time time.Duration) {
	if id == 1 {
		log.Println("clbk called: ", id, time)
		count++
	}
}
