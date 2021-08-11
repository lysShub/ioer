package maps_test

import (
	"bufio"
	"crypto/rand"
	"math/big"
	"net"
	"os"
	"testing"

	"github.com/lysShub/ioer/maps"
)

// 测试maps和map的速度
// 数据来源 https://github.com/17mon/china_ip_list

var maplength int = 250 // map(s)中所有的数据量

// goos: windows
// goarch: amd64
// pkg: github.com/lysShub/ioer/maps/maps_test
// cpu: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz
// BenchmarkMWrite-8    	   87490	     13801 ns/op	       0 B/op	       0 allocs/op
// BenchmarkMSWrite-8   	   30934	     47277 ns/op	   22840 B/op	       0 allocs/op
// BenchmarkMRead-8     	  134210	     10072 ns/op	       0 B/op	       0 allocs/op
// BenchmarkMSRead-8    	  156195	      7770 ns/op	       0 B/op	       0 allocs/op
// PASS
// ok  	github.com/lysShub/ioer/maps/maps_test	6.791s

func BenchmarkMWrite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mwrite()
	}
}

func BenchmarkMSWrite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mswrite()
	}
}

func BenchmarkMRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mread()
	}
}

func BenchmarkMSRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		msrade()
	}
}

// -----------------------------------------

var m map[int64]int = make(map[int64]int)
var ms maps.Maps

var addrs []*net.UDPAddr = make([]*net.UDPAddr, 0, maplength)

func mwrite() {
	for i := 0; i < maplength; i++ {
		m[ider(addrs[i])] = i
	}
}
func mswrite() {
	for i := 0; i < maplength; i++ {
		ms.Add(addrs[i], nil)
	}
}

func mread() {
	for i := 0; i < maplength; i++ {
		if _, ok := m[ider(addrs[i])]; ok {

		}
	}
}

func msrade() {
	for i := 0; i < maplength; i++ {
		if _, ok := ms.Read(addrs[i]); ok {

		}
	}
}

func init() {

	var sh *bufio.Scanner
	if fh, err := os.Open(`./ips.txt`); err != nil {
		panic(err)
	} else {
		sh = bufio.NewScanner(fh)
	}

	var randPort = func() string {
		r, err := rand.Int(rand.Reader, big.NewInt(65535))
		if err != nil {
			panic(err)
		}
		return r.String()
	}

	for i := 0; i < maplength && sh.Scan(); i++ {
		if addr, err := net.ResolveUDPAddr("udp", sh.Text()+":"+randPort()); err != nil {
			continue
		} else {
			addrs = append(addrs, addr)
		}
	}
}

func ider(addr *net.UDPAddr) int64 {
	if addr == nil {
		return 0
	} else {
		addr.IP = addr.IP.To16()
		if addr.IP == nil || len(addr.IP) < 16 {
			return int64(addr.Port)
		} else {
			return int64(addr.IP[12])<<+int64(addr.IP[13])<<32 + int64(addr.IP[14])<<24 + int64(addr.IP[15])<<16 + int64(addr.Port)
		}
	}
}
