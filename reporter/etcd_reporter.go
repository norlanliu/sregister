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
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
	"github.com/golang/glog"
	"github.com/norlanliu/sregister/configuration"
)

type EtcdReporter struct {
	etcdClient client.KeysAPI
	path       string
	key        string
	value      string
}

func (erepoter *EtcdReporter) NewReporter(sc *configuration.ServiceConf) error {
	config := client.Config{
		Endpoints: sc.ReporterHosts,
		Transport: client.DefaultTransport,
	}
	c, err := client.New(config)
	if err != nil {
		glog.Errorf("SRegister: create etcd client failed. Error: %v", err)
		glog.Flush()
		return err
	} else {
		erepoter.etcdClient = client.NewKeysAPI(c)
	}

	erepoter.path = "/"
	if sc.ReporterPath != "" {
		erepoter.path = sc.ReporterPath
	}

	str, derr := getServiceData(sc)
	if derr != nil {
		glog.Errorf("SRegister: geenrate service value failed. Error: %v", derr)
		glog.Flush()
		return derr
	}
	erepoter.value = str
	erepoter.key = ""

	return nil
}

func (erepoter *EtcdReporter) ReportUp() {
	if erepoter.key == "" {
		cioOptions := client.CreateInOrderOptions{}
		response, err := erepoter.etcdClient.CreateInOrder(context.Background(), erepoter.path, erepoter.value, &cioOptions)
		if err != nil {
			glog.Errorf("SRegister: reporter create key %s in order failed. Error: %v", erepoter.path, err)
			glog.Flush()
		} else {
			erepoter.key = response.Node.Key
			glog.Infof("SRegister: reporter create key %s with value %s succeeded.", erepoter.key, erepoter.value)
		}
	} else {
		setOpt := client.SetOptions{}
		_, err := erepoter.etcdClient.Set(context.Background(), erepoter.key, erepoter.value, &setOpt)
		if err != nil {
			glog.Errorf("SRegister: reporter set key %s failed. Error: %v", erepoter.key, err)
			glog.Flush()
		} else {
			glog.Infof("SRegister: reporter set key %s with value %s succeeded.", erepoter.key, erepoter.value)
		}
	}
}

func (erepoter *EtcdReporter) ReportDown() {
	if erepoter.key != "" {
		opt := client.DeleteOptions{}
		_, err := erepoter.etcdClient.Delete(context.Background(), erepoter.key, &opt)

		if err != nil {
			realErr, ok := err.(client.Error)
			if !ok || realErr.Code != client.ErrorCodeKeyNotFound {
				glog.Errorf("SRegister: reporter delete key %s failed. Error: %v", erepoter.key, err)
				glog.Flush()
				return
			}
		}
		erepoter.key = ""
		glog.Infof("SRegister: reporter delete key %s succeeded.", erepoter.key)
	}
}
func (erepoter *EtcdReporter) Ping() bool {
	_, err := erepoter.etcdClient.Get(context.Background(), erepoter.key, &client.GetOptions{})

	return err == nil
}

func (erepoter *EtcdReporter) Close() {
	glog.Infof("SRegister: close reporter")
	erepoter.ReportDown()
}
