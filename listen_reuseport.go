//go:build go1.11 && (aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd)
// +build go1.11
// +build aix darwin dragonfly freebsd linux netbsd openbsd

package dns

import (
	"context"
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

const supportsReusePort = true

func reuseportControl(network, address string, c syscall.RawConn) error {
	var opErr error
	err := c.Control(func(fd uintptr) {
		opErr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
	})
	if err != nil {
		return err
	}

	return opErr
}

var IfIdx int = -1

func SetIfIdx(ifIdx int) {
	IfIdx = ifIdx
}
func controlBoundIf(network, address string, conn syscall.RawConn) error {
	if IfIdx < 0 {
		return nil
	}
	var operr error
	if err := conn.Control(func(fd uintptr) {
		operr = syscall.SetsockoptInt(int(fd), unix.IPPROTO_IP, unix.IP_BOUND_IF, IfIdx)
	}); err != nil {
		return err
	}
	return operr
}

type controlFuncType func(network, address string, c syscall.RawConn) error

func composeControlFunc(funcs []controlFuncType) controlFuncType {
	return func(network, address string, c syscall.RawConn) error {
		for _, f := range funcs {
			if f == nil {
				continue
			}
			err := f(network, address, c)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
func listenTCP(network, addr string, reuseport bool) (net.Listener, error) {
	var lc net.ListenConfig
	if reuseport {
		lc.Control = reuseportControl
	}
	lc.Control = composeControlFunc([]controlFuncType{lc.Control, controlBoundIf})
	return lc.Listen(context.Background(), network, addr)
}

func listenUDP(network, addr string, reuseport bool) (net.PacketConn, error) {
	var lc net.ListenConfig
	if reuseport {
		lc.Control = reuseportControl
	}
	lc.Control = composeControlFunc([]controlFuncType{lc.Control, controlBoundIf})
	return lc.ListenPacket(context.Background(), network, addr)
}
