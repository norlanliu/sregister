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
package reporter

import (
	"encoding/json"
	"errors"
	"github.com/golang/glog"
	"github.com/norlanliu/sregister/configuration"
)

type serviceData struct {
	Host   string
	Port   int
	Weight int
}

type ServiceReporter interface {
	NewReporter(sc *configuration.ServiceConf) error
	ReportUp()
	ReportDown()
	Ping() bool
	Close()
}

func getServiceData(sc *configuration.ServiceConf) (string, error) {
	var sd serviceData

	if sc.Host == "" || sc.Port == 0 {
		err := errors.New("missing required feilds of Host and Port for Reporter")
		glog.Errorf("SRegister: initial reporter failed. Error: %v", err)
		glog.Flush()
		return "", err
	}

	sd.Host = sc.Host
	sd.Port = sc.Port
	sd.Weight = sc.Weight

	ret, err := json.Marshal(&sd)
	if err != nil {
		glog.Errorf("SRegister: json marshal service data failed. Error: %v", err)
		glog.Flush()
		return "", err
	}

	return string(ret), nil
}
