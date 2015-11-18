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
	"github.com/norlanliu/sregister/configuration"
	"io/ioutil"
	"testing"
)

func TestGetServiceData(t *testing.T) {

	fmt.Printf("Test: reporter get service data...\n")
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

	jsonStr, sgErr := getServiceData(sc)

	if sgErr != nil {
		t.Fatalf("Reporter: Get Service data error. %v", sgErr)
	}

	value := &serviceData{}
	err = json.Unmarshal([]byte(jsonStr), value)
	if value.Port != 9000 {
		t.Fatalf("Reporter: Get Service data error. wanted 9000, got %d", value.Port)
	}

	fmt.Printf("... PASS \n")
}
