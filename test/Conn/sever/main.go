package main

import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/lysShub/ioer"
)

func main() {

	var s6, s7 int64
	go func() {

		if l, err := ioer.Listen(&net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 19986}); err != nil {
			panic(err)
		} else {
			for {
				conn := l.Accept()
				go func(conn *ioer.Conn) {
					var da []byte = make([]byte, 2000)
					for {
						if n, err := conn.Read(da); err != nil {
							panic(err)
						} else {
							s6 = s6 + int64(n)
						}
					}

				}(conn)
			}
		}
	}()

	go func() {
		if l, err := ioer.Listen(&net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 19987}); err != nil {
			panic(err)
		} else {
			for {
				conn := l.Accept()
				go func(conn *ioer.Conn) {
					var da []byte = make([]byte, 2000)
					for {
						if n, err := conn.Read(da); err != nil {
							panic(err)
						} else {
							s7 = s7 + int64(n)
						}
					}

				}(conn)
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Printf("%d %d \r", s6>>20, s7>>20)
			s6, s7 = 0, 0
		}
	}()

	http.ListenAndServe("0.0.0.0:8792", nil)
}
