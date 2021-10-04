package main

import (
	"net"
	"net/http"
	_ "net/http/pprof"

	"github.com/lysShub/ioer"
)

func main() {

	A()
	go A()

	http.ListenAndServe(":8792", nil)
}

var severIP net.IP = net.ParseIP("172.21.70.32")

func A() {

	// var path1, path2 string = `E:\浏览器下载\goland-2021.2.exe`, `E:\浏览器下载\Docker Desktop Installer.exe`
	// var h1, h2 *os.File

	var err error
	var conn *ioer.Conn
	conn, err = ioer.Dial(&net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 8080}, &net.UDPAddr{IP: severIP, Port: 19986})
	if err != nil {
		panic(err)
	}

	// h1, err = os.Open(path1)
	// if err != nil {
	// 	panic(err)

	var da []byte = make([]byte, 1024)
	var n int64
	for i := int64(0); ; { //
		// if n, err = h1.ReadAt(da, i); err == nil {
		if _, err = conn.Write(da); err != nil {
			panic(err)
		}
		i = i + int64(n)
		// } else {
		// 	fmt.Println("path1", time.Since(a))
		// 	break
		// }
	}

}

func B() {

	var err error
	var conn *ioer.Conn
	conn, err = ioer.Dial(&net.UDPAddr{IP: nil, Port: 8080}, &net.UDPAddr{IP: severIP, Port: 19987})
	if err != nil {
		panic(err)
	}

	// h2, err = os.Open(path2)
	// if err != nil {
	// 	panic(err)
	// }
	var da []byte = make([]byte, 1500)

	var n int
	for i := int64(0); ; {
		// if n, err = h2.ReadAt(da, i); err == nil {
		if _, err = conn.Write(da); err != nil {
			panic(err)
		}
		i = i + int64(n)
		// } else {
		// 	fmt.Println("path2", time.Since(a))
		// 	break
		// }
	}

}
