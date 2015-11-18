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
	"fmt"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
	"github.com/norlanliu/sregister/configuration"
	"io/ioutil"
	"testing"
)

func TestReporter(t *testing.T) {
	fmt.Printf("Test: etcd reporter...\n")
	filePath := "../example/services/tcp_service.json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("SRegister: service file %s doesn't exists. Error: %v.", filePath, err)
	}

	sc := &configuration.ServiceConf{}
	jsonErr := json.Unmarshal(data, sc)

	if jsonErr != nil {
		t.Fatalf("SRegister: parse service file %s failed. Error: %v.", filePath, jsonErr)
	}

	etcdRep := &EtcdReporter{}

	err = etcdRep.NewReporter(sc)

	if err != nil {
		t.Fatalf("reporter: new reporter failed. %v", err)
	}

	if etcdRep.path == "/" {
		t.Fatalf("reporter: new reporter failed")
	}

	etcdRep.ReportUp()

	response, err := etcdRep.etcdClient.Get(context.Background(), etcdRep.key, &client.GetOptions{})

	if err != nil || response.Node == nil || response.Node.Value == "" {
		t.Fatalf("reporter: reporter up failed, get null value")
	}

	if !etcdRep.Ping() {
		t.Fatalf("reporter: the key doesn't exist after reporting up")
	}

	key := etcdRep.key
	etcdRep.ReportDown()

	_, err = etcdRep.etcdClient.Get(context.Background(), key, &client.GetOptions{})

	realErr, ok := err.(client.Error)
	if !ok || realErr.Code != client.ErrorCodeKeyNotFound {
		t.Fatalf("reporter: reporter down failed")
	}

	fmt.Printf("... PASS\n")

}
