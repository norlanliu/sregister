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
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type httpService struct {
	base service

	url    string
	expect string
	client http.Client
}

func (hs *httpService) newService(name string, host string, port int, check map[string]interface{}) {

	url := "http://"
	url += host
	url += ":" + strconv.Itoa(port)
	hs.base.name = name

	hs.base.fall = 1
	hs.base.rise = 1

	if _, exist := check["fall"]; exist {
		hs.base.fall = int(check["fall"].(float64))
	}
	if _, exist := check["rise"]; exist {
		hs.base.rise = int(check["rise"].(float64))
	}

	timeout := 100 * time.Millisecond
	if _, exist := check["timeout"]; exist {
		timeout = time.Duration(int(check["timeout"].(float64))) * time.Millisecond
	}

	if _, exist := check["uri"]; exist {
		url += check["uri"].(string)
	}

	hs.url = url
	if _, exist := check["expect"]; exist {
		hs.expect = check["expect"].(string)
	}

	hs.client = http.Client{
		Timeout: timeout,
	}

	bufferSize := 1
	if hs.base.fall > hs.base.rise {
		bufferSize = hs.base.fall
	} else if hs.base.rise > 0 {
		bufferSize = hs.base.rise
	}
	hs.base.resultBuffer = make([]bool, 0, bufferSize)
}

func (hs *httpService) up() bool {
	resp, err := hs.client.Get(hs.url)

	checkResult := false
	if err != nil {
		glog.Infof("SRegister: get http service %s error %v", hs.url, err)
		checkResult = false
	} else {
		body, berr := ioutil.ReadAll(resp.Body)
		if berr != nil || string(body) != hs.expect {
			glog.Infof("SRegister: get http service %s wrong, wanted %s, got %s", hs.url, hs.expect, string(body))
			checkResult = false
		} else {
			checkResult = true
		}
	}

	result := hs.base.check(checkResult)

	return result
}
