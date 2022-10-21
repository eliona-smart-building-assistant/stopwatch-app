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
	"testing"
	"time"

	"stopwatch/stopwtch"

	"github.com/eliona-smart-building-assistant/go-utils/log"
)

var (
	count1    int32
	count5    int32
	count4553 int32
	countTot  int32
)

func TestStopwatchManager(t *testing.T) {
	swMng := stopwtch.NewStopwatchManager(stopwatchCallback)

	swMng.Start(1)
	time.Sleep(2 * time.Second)

	t.Log("stop timer 1")
	timeF := swMng.Stop(1)
	if !InTolerance(2, count1, countTot, int32(timeF.Seconds())) {
		t.Error("final time timer 1 doesn't match")
	}

	count1 = 0
	countTot = 0

	swMng.Start(5)
	swMng.Start(4553)
	swMng.Start(1)

	time.Sleep(3 * time.Second)

	tS := swMng.Stop(5).Seconds()
	if !InTolerance(1, 3, count5, int32(tS)) {
		t.Error("timer 5 endtime not match")
	}

	time.Sleep(3 * time.Second)
	tS = swMng.Stop(1).Seconds()
	if !InTolerance(1, 6, int32(tS), count1) {
		t.Error("timer 3 endtime not match")
	}

	time.Sleep(4 * time.Second)
	tS = swMng.Stop(4553).Seconds()
	if !InTolerance(1, 10, int32(tS), count4553) {
		t.Error("timer 4553 endtime not match")
	}

	if !InTolerance(2, 3*3+2*3+4, countTot) {
		t.Error("tot not in tolerance")
	}

	swMng.Start(1234)
	swMng.Start(54321)
	swMng.Start(4432)
	// check, if stop all works
	t.Log("wait stop all")
	swMng.StopAll()
}

func stopwatchCallback(id int32, time time.Duration) {
	countTot++
	log.Debug("test", "clbk called: id: %d, time: %f", id, time.Seconds())
	if id == 1 {
		count1++
	} else if id == 5 {
		count5++
	} else if id == 4553 {
		count4553++
	}
}

func InTolerance(tolerance int32, should int32, is ...int32) bool {
	for i, c := range is {
		if should-tolerance > c || c > should+tolerance {
			log.Error("test", "not in tolerance[%d]: should %d, is %d", i, should, is)
			return false
		}
	}
	return true
}
