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
)

type service struct {
	name       string
	rise       int
	fall       int
	lastResult bool

	resultBuffer []bool
}

type serviceChecker interface {
	newService(name string, host string, port int, check map[string]interface{})
	up() bool
}

func (s *service) check(checkResult bool) bool {
	if len(s.resultBuffer) == 0 {
		s.lastResult = checkResult
		capability := cap(s.resultBuffer)
		for i := 0; i < capability; i++ {
			s.resultBuffer = append(s.resultBuffer, checkResult)
		}
		glog.Infof("SRegister: service check %s initial check returned %b", s.name, checkResult)
	}

	if len(s.resultBuffer) == cap(s.resultBuffer) {
		s.resultBuffer = s.resultBuffer[1:]
	}

	s.resultBuffer = append(s.resultBuffer, checkResult)

	fallSlice := s.resultBuffer[(len(s.resultBuffer) - s.fall):]

	result := false
	for _, v := range fallSlice {
		if v {
			result = true
			break
		}
	}
	if !result {
		if s.lastResult {
			glog.Infof("SRegister: check service %s down after %d failures", s.name, s.fall)
		}
		s.lastResult = false
		return s.lastResult
	}
	result = true
	riseSlice := s.resultBuffer[(len(s.resultBuffer) - s.rise):]
	for _, v := range riseSlice {
		if !v {
			result = false
			break
		}
	}
	if result {
		if !s.lastResult {
			glog.Infof("SRegister: check service %s up after %d successes", s.name, s.rise)
		}
		s.lastResult = true
	}

	return s.lastResult
}
