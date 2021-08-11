package maps_test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"os"
	"strconv"
	"testing"

	"github.com/lysShub/ioer/maps"
)

func TestCoverAmap(t *testing.T) {
	var m maps.Maps

	m.Add(&net.UDPAddr{IP: nil, Port: 19945}, nil)
	if v, ok := m.Read(&net.UDPAddr{IP: nil, Port: 19945}); ok {
		fmt.Println("输出", v)
	} else {
		fmt.Println("无数据")
	}
	m.Delete(&net.UDPAddr{IP: nil, Port: 19945})
	if v, ok := m.Read(&net.UDPAddr{IP: nil, Port: 19945}); ok {
		fmt.Println("输出", v)
	} else {
		fmt.Println("无数据")
	}
}

func hash() {
	// 判断hash函数合理性

	var ipd int64 = 315085026

	fh, _ := os.OpenFile(`D:\Desktop\c.txt`, os.O_APPEND|os.O_RDWR, 0666)

	var m map[int]int = make(map[int]int)

	for i := 0; i < 65536; i++ {
		ipd = randid()
		var k1 uint16 = uint16(((ipd>>40)&0x3)<<14 + ((ipd>>24)&0x3)<<12 + ((ipd>>16)&0xff)<<4 + ipd&0xf)

		var k2 uint16 = uint16((ipd>>26)&0xC000 + (ipd>>12)&0x3000 + (ipd>>12)&0xff0 + ipd&0xf)
		if k1 != k2 {
			fmt.Println(k1, k2)
			panic(ipd)
		}

		if n, ok := m[int(k1)]; ok {
			m[int(k1)] = 1 + n
		} else {
			m[int(k1)] = 1
		}
	}

	for k, v := range m {
		fh.Write([]byte(strconv.Itoa(k) + " " + strconv.Itoa(int(v)) + fmt.Sprintln("")))
	}

	fmt.Println("PASS")
}

func randid() int64 {

	n, err := rand.Int(rand.Reader, big.NewInt(0xffffffffffff))
	if err != nil {
		panic(err)
	}
	return n.Int64()
}
