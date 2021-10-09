package ioer

import (
	"errors"
	"net"
	"sync"
)

// net.DialUDP会占用一个端口, 导致不能进行端口复用
// ioer的Dial和Accept实际均是使用net.ListenUDP实现
// 因此可以进行端口复用, 只要四元组中有一元不一样都可以进行连接

// Conn 表示一个链接
type Conn struct {
	Lconn        *net.UDPConn // unconnected net.UCPConn
	raddr, laddr *net.UDPAddr

	// buf        chan []byte // 读取数据管道, 不能有容量
	// w          *io.PipeWriter
	// r          *io.PipeReader

	buf    []byte // 数据缓存
	buflen int

	listenerid int64 // 所属的listener
	// ing        uint8 // 表示活跃程度, 每进行一次通信都会进行累加
	done bool // Conn关闭flag
}

var dialLock sync.Mutex

// Dial
func Dial(laddr, raddr *net.UDPAddr) (*Conn, error) {
	// 底层本质还是有 ListenUDP

	dialLock.Lock()
	defer dialLock.Unlock()

	var l *Listener
	var ok bool
	var err error

	if raddr.IP == nil {
		raddr.IP = net.ParseIP("0.0.0.0")
	}

	if l, ok = ListenersList[ider(laddr)]; ok {
		if c, ok := l.connList[ider(raddr)]; ok {
			return c, nil
		} else {
			return l.add(raddr), nil
		}
	} else {
		// Listener 不存在
		if l, err = Listen(laddr); err != nil {
			return nil, err
		} else {
			var c = new(Conn)

			// c.buf = make(chan []byte)
			// c.r, c.w = io.Pipe()
			c.Lconn = l.lconn
			c.raddr, c.laddr = raddr, laddr
			c.listenerid = ider(laddr)

			l.Lock()
			l.connList[ider(raddr)] = c
			l.Unlock()

			return c, nil
		}
	}
}

// Read 读取数据; 确保b的长度足够大(65536), 否则会丢失部分数据
func (c *Conn) Read(b []byte) (int, error) {
	if c.done {
		return 0, errClosed
	} else {

		// return copy(b, <-c.buf), nil
		// return c.r.Read(b)
		return 0, nil
	}
}

// Write 发送数据
func (c *Conn) Write(b []byte) (int, error) {
	if c.done {
		return 0, errClosed
	} else {
		// c.ing = c.ing + 1
		return c.Lconn.WriteToUDP(b, c.raddr)
	}
}

// Close 关闭
func (c *Conn) Close() error {

	if l, ok := ListenersList[c.listenerid]; ok {
		l.Lock()
		delete(l.connList, ider(c.raddr))

		var err error
		if len(l.connList) == 0 { // 如果 Listener.connList 中没有连接了，Listener也要关闭
			err = l.Close()
		}
		l.Unlock()
		c.done = true

		return err
	} else {
		return nil
	}
}

var once sync.Once

// 全局变量
var ListenersList map[int64]*Listener // 记录全局Listener

type Listener struct {
	connList   map[int64]*Conn // key: ider(raddr)    val: *Conn
	sync.Mutex                 // 锁, 凡是connList写的地方需要上锁

	lconn *net.UDPConn // net.ListenUDP
	laddr *net.UDPAddr // lconn的地址
	rConn chan *Conn   // 有新生成的Conn
	done  bool         // 已关闭
}

var listenLock sync.Mutex

// Listen 监听本地地址, 不会阻塞
func Listen(laddr *net.UDPAddr) (*Listener, error) {

	once.Do(func() {
		ListenersList = make(map[int64]*Listener)
	})

	if laddr == nil || laddr.Port == 0 {
		return nil, errors.New("invalid laddr")
	} else if laddr.IP == nil {
		if lip, err := getLanIP(); err != nil {
			return nil, err
		} else {
			laddr.IP = lip
		}
	}

	// 已经存在
	if l, ok := ListenersList[ider(laddr)]; ok {
		return l, nil
	}

	if conn, err := net.ListenUDP("udp", laddr); err != nil {
		return nil, err
	} else {
		var l = new(Listener)

		l.connList = make(map[int64]*Conn)
		l.lconn = conn
		l.rConn = make(chan *Conn, 1)

		listenLock.Lock()
		ListenersList[ider(laddr)] = l
		listenLock.Unlock()

		go l.run()

		return l, nil
	}
}

// Accept 结束请求, 类似于TCP的Accept
// 	 如果Listener.run中没有写权鉴，则会接收所有连接请求
func (l *Listener) Accept() *Conn {
	c := <-l.rConn
	l.Lock()
	l.connList[ider(c.raddr)] = c
	l.Unlock()
	return c
}

// Delete 将Conn中此Listener中删除
// 	可以依据删除
func (l *Listener) Delete(conn *Conn) {
	l.Lock()
	delete(l.connList, ider(conn.raddr))
	l.Unlock()
}

// add 增加一个Conn
func (l *Listener) add(raddr *net.UDPAddr) *Conn {
	if c, ok := l.connList[ider(raddr)]; ok {
		return c
	} else {
		c = new(Conn)
		// c.buf = make(chan []byte)
		// c.r, c.w = io.Pipe()
		c.raddr, c.Lconn = raddr, l.lconn
		c.listenerid = ider(l.laddr)

		l.Lock()
		l.connList[ider(c.raddr)] = c
		l.Unlock()
		return c
	}
}

// Close 会关闭所有链接
func (l *Listener) Close() error {
	l.done = true
	return l.lconn.Close()
}
