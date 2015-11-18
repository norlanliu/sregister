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
package sregistermain

import (
	"github.com/golang/glog"
	"github.com/norlanliu/sregister/configuration"
	"github.com/norlanliu/sregister/watcher"

	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

const (
	dirCheckInterval = 5
)

type SRegister struct {
	serviceConfDir   string
	services         map[string]time.Time
	serviceDoneChans map[string]chan bool
}

func NewSRegister(arguments []string) *SRegister {
	sr := &SRegister{}
	sr.services = make(map[string]time.Time)
	sr.serviceDoneChans = make(map[string]chan bool)
	cfg := configuration.NewConfigure()

	cfg.ParseConfigure(arguments)

	//remove the slash
	sr.serviceConfDir = strings.TrimSuffix(cfg.GetServiceConfDir(), "/")
	return sr
}

func (sr *SRegister) checkServiceFileExt(name string) bool {
	return filepath.Ext(name) == ".json"
}

func (sr *SRegister) Run() {
	glog.Infof("SRegister: start register")

	//watch the service conf dir
	for {
		launchServices := make(map[string]time.Time)
		relaunchServices := make(map[string]time.Time)
		remainServices := make(map[string]bool)
		size := len(sr.services)

		files, err := ioutil.ReadDir(sr.serviceConfDir)
		if err != nil {
			glog.Infof("SRegister: read service configuration directory failed. Error: %v", err)
		} else {
			for _, file := range files {
				fileName := file.Name()
				if sr.checkServiceFileExt(fileName) {
					if _, exist := sr.services[fileName]; !exist {
						launchServices[fileName] = file.ModTime()
					} else {
						size -= 1
						if !(sr.services[fileName].Equal(file.ModTime())) {
							relaunchServices[fileName] = file.ModTime()
						}
						remainServices[fileName] = true
					}
				}
			}
		}

		if size != 0 {
			for k, _ := range sr.services {
				if _, exist := remainServices[k]; !exist {
					sr.stopServiceWatcher(k)
				}
			}
		}

		for key, mtime := range launchServices {
			sr.launchServiceWatcher(key, mtime)
		}
		for key, mtime := range relaunchServices {
			close(sr.serviceDoneChans[key])
			sr.launchServiceWatcher(key, mtime)
		}

		time.Sleep(dirCheckInterval * time.Second)
	}
}

func (sr *SRegister) stopServiceWatcher(key string) {
	fileName := strings.TrimSuffix(key, filepath.Ext(key))
	glog.Infof("SRegister: stop service watcher %s", fileName)
	close(sr.serviceDoneChans[key])
	delete(sr.services, key)
	delete(sr.serviceDoneChans, key)
}

func (sr *SRegister) launchServiceWatcher(key string, mtime time.Time) {
	fileName := strings.TrimSuffix(key, filepath.Ext(key))
	glog.Infof("SRegister: launch service watcher %s", fileName)

	serviceFilePath := sr.serviceConfDir + "/" + key
	done, err := watcher.LaunchWatcher(serviceFilePath)
	if err != nil {
		glog.Errorf("SRegister: launch service watcher error %v", err)
		glog.Flush()
	} else {
		sr.services[key] = mtime
		sr.serviceDoneChans[key] = done
	}
}
