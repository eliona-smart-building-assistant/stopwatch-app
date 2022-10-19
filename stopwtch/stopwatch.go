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

package stopwtch

import (
	"sync"
	"time"

	"github.com/eliona-smart-building-assistant/go-utils/log"
)

const SMODULE = "Stopwatch"

type SwCallback func(id int32, time time.Duration)

type Stopwatch struct {
	running bool
	ticker  *time.Ticker
	time    time.Duration
	ir      chan bool

	Id int32
}

type StopwatchManager struct {
	stopwatches  []*Stopwatch
	lock         sync.Mutex
	wg           sync.WaitGroup
	callbackFunc SwCallback
}

func NewStopwatchManager(callback SwCallback) *StopwatchManager {
	return &StopwatchManager{
		stopwatches:  []*Stopwatch{},
		lock:         sync.Mutex{},
		callbackFunc: callback,
	}
}

func (swM *StopwatchManager) Start(id int32) {
	swM.lock.Lock()
	defer swM.lock.Unlock()

	sw, _ := swM.getStopwatch(id)
	if sw == nil {
		sw = swM.createNewStopwatch(id)
	}

	if !sw.running {
		swM.wg.Add(1)
		go sw.Start(&swM.wg, swM.callbackFunc)
	}
}

func (swM *StopwatchManager) Stop(id int32) time.Duration {
	var lasttime time.Duration

	swM.lock.Lock()
	defer swM.lock.Unlock()

	sw, index := swM.getStopwatch(id)

	if sw != nil {
		if sw.IsRunning() {
			log.Debug(SMODULE, "timer %d stopped", sw.Id)
			sw.Stop()
		}

		lasttime = sw.GetTime()

		swM.deleteStopwatch(index)
	} else {
		log.Debug(SMODULE, "timer to stop is nil")
	}
	return lasttime
}

func (swM *StopwatchManager) StopAll() {
	swM.lock.Lock()
	defer swM.lock.Unlock()

	for _, sw := range swM.stopwatches {
		log.Debug(SMODULE, "stop all %d", sw.Id)
		sw.Stop()
	}
	swM.wg.Wait()
}

func (swM *StopwatchManager) getStopwatch(id int32) (*Stopwatch, int) {
	var (
		stopwatch *Stopwatch
		index     int
	)

	for index, stopwatch = range swM.stopwatches {
		if stopwatch.Id == id {
			return stopwatch, index
		}
	}

	return nil, -1
}

func (swM *StopwatchManager) createNewStopwatch(id int32) *Stopwatch {
	stopwatch := &Stopwatch{
		Id:      id,
		running: false,
		time:    0,
		ticker:  time.NewTicker(1 * time.Second),
		ir:      make(chan bool),
	}
	swM.stopwatches = append(swM.stopwatches, stopwatch)
	stopwatchPnt, _ := swM.getStopwatch(id)
	return stopwatchPnt
}

func (swM *StopwatchManager) deleteStopwatch(index int) {
	swM.stopwatches = append(swM.stopwatches[:index], swM.stopwatches[index+1:]...)
}

func (sw *Stopwatch) Start(wg *sync.WaitGroup, clbk SwCallback) {
	sw.setRunning()
	go sw.runner(wg, clbk)
}

func (sw *Stopwatch) runner(wg *sync.WaitGroup, clbk SwCallback) {
	defer wg.Done()
	defer sw.setStopped()

	for {
		select {
		case <-sw.ir:
			log.Debug(SMODULE, "ticker interrupted")
			return
		case _, ok := <-sw.ticker.C:
			if ok {
				sw.time += time.Second
				if clbk != nil {
					clbk(sw.Id, sw.time)
				}
			} else {
				log.Debug(SMODULE, "sw ticker not ok")
				return
			}
		}
	}
}

func (sw *Stopwatch) Stop() {
	sw.ticker.Stop()
	sw.ir <- true
}

func (sw *Stopwatch) GetTime() time.Duration {
	return sw.time
}

func (sw *Stopwatch) setRunning() {
	sw.running = true
}

func (sw *Stopwatch) IsRunning() bool {
	return sw.running
}

func (sw *Stopwatch) setStopped() {
	sw.running = false
}
