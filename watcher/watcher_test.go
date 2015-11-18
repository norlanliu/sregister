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
	"bufio"
	"errors"
	"fmt"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
	"github.com/norlanliu/sregister/configuration"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestParseServiceJson(t *testing.T) {
	fmt.Println("Test: watcher parse service json...")

	jsonFilePath := "../example/services/tcp_service.json"
	sc := &configuration.ServiceConf{}

	checks, err := parseServiceJson(jsonFilePath, sc)

	if err != nil {
		t.Fatalf("parse json file error: %v.", err)
	}

	if sc.ReporterType != "etcd" {
		t.Fatalf("parse json file failed. wanted: etcd, got: %s", sc.ReporterType)
	}

	if len(checks) == 1 {
		data := checks[0].(*tcpService)
		if data.base.name != "your_tcp_service_name" {
			t.Fatalf("parse json file error, wanted: your_tcp_service_name, got: %s", data.base.name)
		}
		if data.host != "127.0.0.1" {
			t.Fatalf("parse json file error, wanted: 127.0.0.1, got: %s", data.host)
		}
	} else {
		t.Fatalf("parse json file error, no checks")
	}

	fmt.Println("... PASS")
}

func tcpServer(server string, port int, done chan bool) {

	// listen on all interfaces
	ln, _ := net.Listen("tcp", server+":"+strconv.Itoa(port))

	defer ln.Close()
	// run loop forever (or until ctrl-c)
	for {
		select {
		case <-done:
			return
		default:
		}
		// accept connection on port
		conn, _ := ln.Accept()
		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(conn).ReadString('\n')
		// sample process for string received
		newmessage := strings.ToUpper(message)
		// send new string back to client
		conn.Write([]byte(newmessage + "\n"))
	}
}

func getNodesNumFromEtcd(sc *configuration.ServiceConf) (int, error) {
	cfg := client.Config{
		Endpoints: sc.ReporterHosts,
	}
	c, err := client.New(cfg)
	if err != nil {
		return 0, err
	}
	kapi := client.NewKeysAPI(c)

	getOpt := client.GetOptions{Recursive: true}
	response, err := kapi.Get(context.Background(), sc.ReporterPath, &getOpt)

	if err != nil {
		return 0, err
	}

	if !response.Node.Dir {
		return 0, errors.New("get path not a dir")
	}

	return len(response.Node.Nodes), nil
}
func checkEtcd(size int, sc *configuration.ServiceConf) error {
	cfg := client.Config{
		Endpoints: sc.ReporterHosts,
	}
	c, err := client.New(cfg)
	if err != nil {
		return err
	}
	kapi := client.NewKeysAPI(c)

	getOpt := client.GetOptions{Recursive: true}
	response, err := kapi.Get(context.Background(), sc.ReporterPath, &getOpt)

	if err != nil {
		return err
	}

	if !response.Node.Dir {
		return errors.New("get path not a dir")
	}

	if len(response.Node.Nodes) != size {
		return errors.New("there is no child of service path")
	}

	return nil
}

func TestLaunchWatcher(t *testing.T) {
	fmt.Println("Test: launch watcher...")

	jsonFilePath := "../example/services/tcp_service.json"
	sc := &configuration.ServiceConf{}

	_, err := parseServiceJson(jsonFilePath, sc)

	if err != nil {
		t.Fatalf("parse json file error: %v.", err)
	}

	//start server
	done := make(chan bool)
	defer close(done)
	go tcpServer(sc.Host, sc.Port, done)

	time.Sleep(2 * time.Second)

	//get the number of nodes before launching watcher
	num, etcdErr := getNodesNumFromEtcd(sc)
	if etcdErr != nil {
		t.Fatalf("Watcher: get the number of nodes from etcd failed, Error: %v", err)
	}

	//launch watcher
	closeChan, err := LaunchWatcher(jsonFilePath)

	if err != nil {
		t.Fatalf("Watcher: Launch watcher error %v", err)
	}

	//give it some time to report
	time.Sleep(time.Second)
	err = checkEtcd(num+1, sc)

	if err != nil {
		t.Fatalf("Watcher: report service error %v", err)
	}

	//close the tcp server, then check the service's status
	done <- true
	time.Sleep(6 * time.Second)
	err = checkEtcd(num, sc)

	if err != nil {
		t.Fatalf("Watcher: close reporter error failed")
	}

	//restart the tcpserver and check the service's status
	go tcpServer(sc.Host, sc.Port, done)

	time.Sleep(8 * time.Second)
	err = checkEtcd(num+1, sc)

	if err != nil {
		t.Fatalf("Watcher: report service error %v", err)
	}

	//close the watcher
	close(closeChan)

	time.Sleep(2 * time.Second)
	err = checkEtcd(num, sc)

	if err != nil {
		t.Fatalf("Watcher: close reporter error failed. Error: %v", err)
	}
	fmt.Println("... PASS")
}
