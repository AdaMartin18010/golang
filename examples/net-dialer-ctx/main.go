// Package main demonstrates Go 1.26 net.Dialer Context-aware dialing
//
// Go 1.26 adds DialIP, DialTCP, DialUDP, and DialUnix methods to net.Dialer,
// all accepting context.Context for cancellation and timeout control.
package main

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"time"
)

func main() {
	fmt.Println("Go 1.26 net.Dialer Context-aware Demo")
	fmt.Println("======================================")

	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	// Example 1: DialTCP with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	fmt.Println("\n1. DialTCP with context timeout:")
	fmt.Println("   dialer.DialTCP(ctx, \"tcp\", laddr, raddr)")

	// Convert net.TCPAddr to netip.AddrPort for Go 1.26 API
	remoteAddr, _ := netip.AddrPortFrom(netip.MustParseAddr("93.184.216.34"), 80).MarshalBinary()
	_ = remoteAddr
	// In Go 1.26, DialTCP takes netip.AddrPort instead of *net.TCPAddr
	conn, err := dialer.DialTCP(ctx, "tcp", netip.AddrPort{}, netip.AddrPort{})
	if err != nil {
		fmt.Printf("   Result: connection failed (expected in demo): %v\n", err)
	} else {
		defer conn.Close()
		fmt.Printf("   Result: connected to %s\n", conn.RemoteAddr())
	}

	// Example 2: Context cancellation
	fmt.Println("\n2. Context cancellation:")
	ctx2, cancel2 := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("   Cancelling context...")
		cancel2()
	}()

	conn2, err := dialer.DialTCP(ctx2, "tcp", netip.AddrPort{}, netip.AddrPort{})
	if err != nil {
		fmt.Printf("   Result: %v\n", err)
	} else {
		conn2.Close()
	}

	fmt.Println("\nGo 1.26 Dialer new methods:")
	fmt.Println("  DialIP(ctx, network, laddr, raddr)")
	fmt.Println("  DialTCP(ctx, network, laddr, raddr)")
	fmt.Println("  DialUDP(ctx, network, laddr, raddr)")
	fmt.Println("  DialUnix(ctx, network, laddr, raddr)")
	fmt.Println("\nAll accept context.Context for cancellation and deadlines.")
	fmt.Println("Note: laddr and raddr are netip.AddrPort in Go 1.26.")
}
