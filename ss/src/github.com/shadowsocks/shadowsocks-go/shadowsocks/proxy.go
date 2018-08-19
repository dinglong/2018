package shadowsocks

import (
	"net"
	"time"
)

type ProxyConn struct {
	*Conn
	raddr *ProxyAddr
}

type ProxyAddr struct {
	network string
	address string
}

func (c *ProxyConn) LocalAddr() net.Addr {
	return c.Conn.LocalAddr()
}

func (c *ProxyConn) RemoteAddr() net.Addr {
	return c.raddr
}

func (c *ProxyConn) SetDeadline(t time.Time) error {
	return c.Conn.SetDeadline(t)
}

func (c *ProxyConn) SetReadDeadline(t time.Time) error {
	return c.Conn.SetReadDeadline(t)
}

func (c *ProxyConn) SetWriteDeadline(t time.Time) error {
	return c.Conn.SetWriteDeadline(t)
}

func (a *ProxyAddr) Network() string {
	return a.network
}

func (a *ProxyAddr) String() string {
	return a.address
}
