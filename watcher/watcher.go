//Copyright (c) 2015 Qi Liu AT ICT
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.
package watcher

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/norlanliu/sregister/configuration"
	"github.com/norlanliu/sregister/reporter"
	"io/ioutil"
	"time"
)

//service watcher routine function
func watch(sc serviceChecker, result chan bool) {
	result <- sc.up()
}

//generate a reporter based on the reporter type
func generateReporter(sc *configuration.ServiceConf) reporter.ServiceReporter {
	switch sc.ReporterType {
	case "etcd":
		sreporter := new(reporter.EtcdReporter)
		err := sreporter.NewReporter(sc)
		if err == nil {
			return sreporter
		}
	}
	return nil
}

//Watch and report the status of service
//every checker has a routine
func watchAndReport(checks []serviceChecker, sc *configuration.ServiceConf) chan bool {

	done := make(chan bool)
	go func() {
		checksSize := len(checks)
		resultChans := make(chan bool, checksSize)
		reporter := generateReporter(sc)
		if reporter == nil {
			glog.Errorf("SRegister: new reporter error.")
			glog.Flush()
			return
		}

		defer reporter.Close()

		oldResult := false
		newResult := false
		for {
			for i := 0; i < checksSize; i++ {
				go watch(checks[i], resultChans)
			}
			finalResult := true
			for i := 0; i < checksSize; i++ {
				if !(<-resultChans) {
					finalResult = false
					break
				}
			}
			newResult = finalResult
			if !reporter.Ping() {
				oldResult = !newResult
			}
			if oldResult != newResult {
				if newResult {
					reporter.ReportUp()
					glog.Infof("SRegister: service %s is up now", sc.Name)
				} else {
					reporter.ReportDown()
					glog.Infof("SRegister: service %s is down now", sc.Name)
				}
				oldResult = newResult
			}

			//use time.Sleep() to realise check interval
			checkTicker := sc.CheckInterval
			for checkTicker != 0 {
				select {
				case <-done:
					reporter.Close()
					glog.Infof("SRegister: close the watcher of service %s", sc.Name)
					return
				default:
				}
				time.Sleep(time.Second)
				checkTicker -= 1
			}
		}
	}()

	return done
}

//parse service json file
func parseServiceJson(filePath string, sc *configuration.ServiceConf) ([]serviceChecker, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		glog.Warningf("SRegister: service file %s doesn't exists. Error: %v.", filePath, err)
		return nil, err
	}

	jsonErr := json.Unmarshal(data, sc)

	if jsonErr != nil {
		glog.Errorf("SRegister: parse service file %s failed. Error: %v.", filePath, jsonErr)
		glog.Flush()
		return nil, jsonErr
	}

	checksSize := len(sc.Checks)
	checks := make([]serviceChecker, 0, checksSize)
	for i := 0; i < checksSize; i++ {
		checkData := sc.Checks[i].(map[string]interface{})
		switch checkData["type"].(string) {
		case "tcp":
			ts := new(tcpService)
			ts.newService(sc.Name, sc.Host, sc.Port, checkData)
			checks = append(checks, ts)
		case "http":
			hs := new(httpService)
			hs.newService(sc.Name, sc.Host, sc.Port, checkData)
			checks = append(checks, hs)
		}
	}

	return checks, nil
}

//A Service to A Watcher (namely a goroutine)
func LaunchWatcher(servicePath string) (chan bool, error) {

	var sc configuration.ServiceConf
	checks, err := parseServiceJson(servicePath, &sc)

	if err == nil {
		done := watchAndReport(checks, &sc)
		return done, nil
	} else {
		return nil, err
	}
}
