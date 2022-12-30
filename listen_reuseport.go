//go:build go1.11 && (aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd) && !android
// +build go1.11
// +build aix darwin dragonfly freebsd linux netbsd openbsd
// +build !android

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
		if opErr == nil {
			if IfIdx >= 0 {
				opErr = syscall.SetsockoptInt(int(fd), unix.IPPROTO_IP, unix.IP_BOUND_IF, IfIdx)
			}
		}
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

func listenTCP(network, addr string, reuseport bool) (net.Listener, error) {
	var lc net.ListenConfig
	if reuseport {
		lc.Control = reuseportControl
	}
	return lc.Listen(context.Background(), network, addr)
}

func listenUDP(network, addr string, reuseport bool) (net.PacketConn, error) {
	var lc net.ListenConfig
	if reuseport {
		lc.Control = reuseportControl
	}
	return lc.ListenPacket(context.Background(), network, addr)
}
