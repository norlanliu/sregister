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
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestHttpUp(t *testing.T) {

	fmt.Println("Test: test http service up...")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	}))

	u, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("Watcher: http service watcher parse url failed")
	}

	qi := strings.Index(u.Host, ":")

	name := "test_service"
	host := u.Host[:qi]
	port, _ := strconv.Atoi(u.Host[qi+1:])

	check := make(map[string]interface{})
	check["rise"] = 3.0
	check["fall"] = 2.0
	check["uri"] = ""
	check["expect"] = "ok"

	ts := &httpService{}
	ts.newService(name, host, port, check)

	for i := 0; i < 100; i++ {
		ans := ts.up()

		if !ans {
			t.Fatalf("Watcher: tcp service watcher check up wrong")
		}
	}

	//close http server
	server.Close()
	time.Sleep(2 * time.Second)

	ans := true
	for i := 0; i < 10; i++ {
		ans = ts.up()
	}
	if ans {
		t.Fatalf("Watcher: tcp service watcher check down wrong")
	}
	fmt.Println("... PASS")

}
