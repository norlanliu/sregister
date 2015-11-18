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
package configuration

import (
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"testing"
)

func TestParseConfigure(t *testing.T) {

	args := []string{
		"-services_dir=the_command_dir",
		"-log_dir=hehe",
	}

	servicesDir := "the_command_dir"
	logDir := "hehe"

	cfg := NewConfigure()
	err := cfg.ParseConfigure(args)
	if err != nil {
		fmt.Errorf("%v\n", err)
	}

	if cfg.serviceConfDir != servicesDir {
		t.Fatalf("parse configure error, services_dir: %s, wanted: %s", cfg.serviceConfDir, servicesDir)
	}
	if cfg.logDir != logDir {
		t.Fatalf("parse configure error, log_dir: %s, wanted: %s", cfg.logDir, logDir)
	}

	//Test version
	args = []string{
		"-version",
	}

	cfg = NewConfigure()
	err = cfg.ParseConfigure(args)
	if err != nil {
		fmt.Errorf("%v\n", err)
	}

	// Test Environment

	os.Clearenv()

	servicesDir = "env_services_dir"
	logDir = "env_log_dir"
	os.Setenv("SREG_SERVICES_DIR", servicesDir)
	os.Setenv("SREG_LOG_DIR", logDir)

	cfg = NewConfigure()
	err = cfg.ParseConfigure([]string{})

	if cfg.serviceConfDir != servicesDir {
		t.Fatalf("parse configure error, services_dir: %s, wanted: %s", cfg.serviceConfDir, servicesDir)
	}
	if cfg.logDir != logDir {
		t.Fatalf("parse configure error, log_dir: %s, wanted: %s", cfg.logDir, logDir)
	}

}

func TestParseConfigureFromCE(t *testing.T) {
	//Test command-line and env

	logDir := "hehe3"
	args := []string{
		"--log_dir=hehe3",
	}
	os.Clearenv()
	os.Setenv("SREG_LOG_DIR", "3_log_dir")
	cfg := NewConfigure()
	cfg.ParseConfigure(args)

	if cfg.logDir != logDir {
		t.Fatalf("parse configure error, log_dir: %s, wanted: %s", cfg.logDir, logDir)
	}
}

func TestLogDir(t *testing.T) {
	logDir := "hehe"
	args := []string{
		"--log_dir=hehe",
	}
	cfg := NewConfigure()
	cfg.ParseConfigure(args)

	if _, err := os.Stat(logDir); !os.IsNotExist(err) {
		derr := os.RemoveAll(logDir)
		if derr != nil {
			t.Fatalf("remove dir failed, Error: %v", derr)
		}
	}

	derr := os.Mkdir(logDir, os.ModePerm)
	if derr != nil {
		t.Fatalf("create dir failed, Error: %v", derr)
	}

	glog.Infof("test...")
	glog.Flush()

	if files, err := ioutil.ReadDir(logDir); err != nil || len(files) == 0 {
		t.Fatalf("write log failed")
	}

	if _, err := os.Stat(logDir); !os.IsNotExist(err) {
		derr := os.RemoveAll(logDir)
		if derr != nil {
			t.Fatalf("remove dir failed, Error: %v", derr)
		}
	}
}
