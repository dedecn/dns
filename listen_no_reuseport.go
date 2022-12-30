//go:build !go1.11 || !(aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd) || android
// +build !go1.11 !aix,!darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd android

package dns

import "net"

const supportsReusePort = false

func SetIfIdx(ifIdx int) {

}
func listenTCP(network, addr string, reuseport bool) (net.Listener, error) {
	if reuseport {
		// TODO(tmthrgd): return an error?
	}

	return net.Listen(network, addr)
}

func listenUDP(network, addr string, reuseport bool) (net.PacketConn, error) {
	if reuseport {
		// TODO(tmthrgd): return an error?
	}

	return net.ListenPacket(network, addr)
}
