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
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"
)

func server(port int, done chan bool) {

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":"+strconv.Itoa(port))

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
func TestUp(t *testing.T) {

	fmt.Println("Test: test tcp service up...")

	name := "test_service"
	host := "localhost"
	port := 9000

	check := make(map[string]interface{})
	check["rise"] = 3.0
	check["fall"] = 2.0

	done := make(chan bool)
	go server(port, done)

	ts := &tcpService{}
	ts.newService(name, host, port, check)

	for i := 0; i < 100; i++ {
		ans := ts.up()

		if !ans {
			t.Fatalf("Watcher: tcp service watcher check up wrong")
		}
	}

	close(done)
	ans := true
	for i := 0; i < 10; i++ {
		ans = ts.up()
	}
	if ans {
		t.Fatalf("Watcher: tcp service watcher check down wrong")
	}
	fmt.Println("... PASS")

}
