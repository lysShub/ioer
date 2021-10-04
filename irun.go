// get ioer.Conn

package ioer

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"
)

var errClosed error = errors.New("connection closed")

// run 运行：1.路由接收到的UDP数据包; 2.接收新请求, 生成新的Conn
func (l *Listener) run() {

	var (
		id    int64
		c     *Conn
		ok    bool
		n     int
		raddr *net.UDPAddr
		err   error
		tmp   []byte = make([]byte, 65536)
	)

	defer func() {
		if recover() != nil {
			// l.Close() 会关闭l.lconn，会导致ReadFromUDP发生panic
		}
	}()

	for !l.done {
		if n, raddr, err = l.lconn.ReadFromUDP(tmp); !l.done && err != nil {
			errlog(err)

		} else if n > 0 {
			id = ider(raddr)

			if c, ok = l.connList[id]; !ok { // 新增连接Conn 瓶颈
				// 权鉴

				// 新增连接
				c = new(Conn)
				// c.buf= make(chan []byte)
				c.r, c.w = io.Pipe()
				c.raddr, c.Lconn = raddr, l.lconn
				c.listenerid = ider(l.laddr)
				select {
				case l.rConn <- c:
				default:
				}
			}

			// if len(c.io) < cap(c.io) {
			// 	c.io <- l.tmp[:n]
			// }

			c.ing = c.ing + 1

			c.Write(tmp[:n])

			// select {
			// // case c.buf <- tmp[:n]: // 写入数据
			// case :
			// default:
			// }

		}
	}
}

func ider(addr *net.UDPAddr) int64 {
	if addr == nil {
		return 0
	} else {
		if len(addr.IP) < 16 {
			if addr.IP = addr.IP.To16(); len(addr.IP) < 16 {
				return int64(addr.Port)
			}
		}
		return int64(addr.IP[12])<<+int64(addr.IP[13])<<32 + int64(addr.IP[14])<<24 + int64(addr.IP[15])<<16 + int64(addr.Port)
	}
}

func getLanIP() (net.IP, error) {
	conn, err := net.DialTimeout("ip4:1", "8.8.8.8", time.Second*2)
	if err != nil {
		return nil, err
	}
	return net.ParseIP(conn.LocalAddr().String()), nil
}

// errlog 输出信息至os.Stderr
func errlog(err ...error) bool {
	var haveErr bool = false
	for i, e := range err {
		if e != nil {
			haveErr = true
			_, fp, ln, _ := runtime.Caller(1)
			if len(err) == 1 {
				fmt.Fprintln(os.Stderr, fp+":"+strconv.Itoa(ln)+" :\n    "+e.Error())
			} else {
				fmt.Fprintln(os.Stderr, "["+strconv.Itoa(i+1)+"]. "+fp+":"+strconv.Itoa(ln)+" \n    "+e.Error())
			}
		}
	}
	return haveErr
}
