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
	"net"
	"strconv"
	"time"
)

type tcpService struct {
	base service

	host    string
	port    int
	timeout time.Duration
}

func (ts *tcpService) newService(name string, host string, port int, check map[string]interface{}) {

	ts.host = host
	ts.port = port
	ts.base.name = name

	ts.base.fall = 1
	ts.base.rise = 1

	ts.timeout = 100 * time.Millisecond

	if _, exist := check["fall"]; exist {
		ts.base.fall = int(check["fall"].(float64))
	}
	if _, exist := check["rise"]; exist {
		ts.base.rise = int(check["rise"].(float64))
	}
	if _, exist := check["timeout"]; exist {
		ts.timeout = time.Duration(int(check["timeout"].(float64))) * time.Millisecond
	}

	bufferSize := 1
	if ts.base.fall > ts.base.rise {
		bufferSize = ts.base.fall
	} else if ts.base.rise > 0 {
		bufferSize = ts.base.rise
	}
	ts.base.resultBuffer = make([]bool, 0, bufferSize)
}

func (ts *tcpService) up() bool {
	address := ts.host + ":" + strconv.Itoa(ts.port)

	conn, err := net.DialTimeout("tcp", address, ts.timeout)

	checkResult := false
	if err != nil {
		glog.Infof("SRegister: connect to tcp service %s error %v", address, err)
	} else {
		conn.Close()
		checkResult = true
	}

	result := ts.base.check(checkResult)

	return result
}
