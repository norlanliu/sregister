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
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
	"github.com/norlanliu/sregister/configuration"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func copyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
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

func parseServiceJson(filePath string, sc *configuration.ServiceConf) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	jsonErr := json.Unmarshal(data, sc)

	if jsonErr != nil {
		return jsonErr
	}

	return nil
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

func startRegisterServer(confDir string) {
	args := []string{"-services_dir=" + confDir}

	sr := NewSRegister(args)
	go sr.Run()

	time.Sleep(2 * time.Second)
}

func testAddNewServiceFile(confDir string, t *testing.T) {
	example := "../example/services/tcp_service.json"
	if _, err := os.Stat(confDir); !os.IsNotExist(err) {
		derr := os.RemoveAll(confDir)
		if derr != nil {
			t.Fatalf("remove dir failed, Error: %v", derr)
		}
	}
	derr := os.Mkdir(confDir, os.ModePerm)
	if derr != nil {
		t.Fatalf("create dir failed, Error: %v", derr)
	}

	//parse service json file
	sc := &configuration.ServiceConf{}

	err := parseServiceJson(example, sc)
	if err != nil {
		t.Fatalf("parse service json file failed, Error: %v", err)
	}
	//start server
	done := make(chan bool)
	defer close(done)
	go tcpServer(sc.Host, sc.Port, done)

	time.Sleep(2 * time.Second)

	num, etcdErr := getNodesNumFromEtcd(sc)
	if etcdErr != nil {
		t.Fatalf("sregister: get the number of nodes from etcd failed, Error: %v", err)
	}
	//copy service file to conf dir
	dstName := confDir + "/" + "your_tcp_service.json"
	copyFile(dstName, example)
	time.Sleep(5 * time.Second)

	err = checkEtcd(num+1, sc)

	if err != nil {
		t.Fatalf("Watcher: report service error %v", err)
	}

	//stop the tcp server
	done <- true

	time.Sleep(5 * time.Second)
	err = checkEtcd(num, sc)

	if err != nil {
		t.Fatalf("Watcher: close reporter error failed, %v", err)
	}
	derr = os.RemoveAll(confDir)
	if derr != nil {
		t.Fatalf("remove dir failed, Error: %v", derr)
	}
}

func testModifyServiceFile(confDir string, t *testing.T) {
	example := "../example/services/another_tcp_service.json"
	if _, err := os.Stat(confDir); !os.IsNotExist(err) {
		derr := os.RemoveAll(confDir)
		if derr != nil {
			t.Fatalf("remove dir failed, Error: %v", derr)
		}
	}
	derr := os.Mkdir(confDir, os.ModePerm)
	if derr != nil {
		t.Fatalf("create dir failed, Error: %v", derr)
	}
	//copy file
	dstName := confDir + "/" + "your_tcp_service.json"
	copyFile(dstName, example)

	//parse service json file
	sc := &configuration.ServiceConf{}

	err := parseServiceJson(example, sc)
	if err != nil {
		t.Fatalf("parse service json file failed, Error: %v", err)
	}

	num, etcdErr := getNodesNumFromEtcd(sc)
	if etcdErr != nil {
		t.Fatalf("sregister: get the number of nodes from etcd failed, Error: %v", err)
	}
	//start server
	done := make(chan bool)
	defer close(done)
	go tcpServer(sc.Host, sc.Port, done)

	time.Sleep(5 * time.Second)

	err = checkEtcd(num+1, sc)

	if err != nil {
		t.Fatalf("Watcher: report service error %v", err)
	}

	//remove service file
	derr = os.RemoveAll(confDir)
	if derr != nil {
		t.Fatalf("remove dir failed, Error: %v", derr)
	}

	time.Sleep(5 * time.Second)
	err = checkEtcd(num, sc)

	if err != nil {
		t.Fatalf("Watcher: close reporter error failed, %v", err)
	}
}
func TestRun(t *testing.T) {
	fmt.Println("Test: sregister run...")

	rand.Seed(time.Now().UnixNano())

	dirName := RandStringRunes(10)
	confDir := "/tmp/" + dirName

	startRegisterServer(confDir)

	//test add a new file
	testAddNewServiceFile(confDir, t)

	//test modify and remove service file
	testModifyServiceFile(confDir, t)

	fmt.Println("... PASS")
}
