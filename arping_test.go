package arping

import (
	"net"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestPingWithInvalidIP(t *testing.T) {
	ip := net.ParseIP("invalid")

	_, _, err := Ping(ip)
	if err == nil {
		t.Error("error expected")
	}

	validateInvalidV4AddrErr(t, err)
}

func TestPingWithV6IP(t *testing.T) {
	ip := net.ParseIP("fe80::e2cb:4eff:fed5:ca4e")

	_, _, err := Ping(ip)
	if err == nil {
		t.Error("error expected")
	}

	validateInvalidV4AddrErr(t, err)
}

func TestGratuitousArpWithInvalidIP(t *testing.T) {
	ip := net.ParseIP("invalid")

	err := GratuitousArp(ip)
	if err == nil {
		t.Error("error expected")
	}

	validateInvalidV4AddrErr(t, err)
}

func TestGratuitousArpWithV6IP(t *testing.T) {
	ip := net.ParseIP("fe80::e2cb:4eff:fed5:ca4e")

	err := GratuitousArp(ip)
	if err == nil {
		t.Error("error expected")
	}

	validateInvalidV4AddrErr(t, err)
}

/*
	func TestRedundantIPCheckWithGratuitousArpOverIface(t *testing.T) {
		EnableVerboseLog()

		srcIP := net.ParseIP("10.10.10.133")
		ifaceName := "ens34"
		iface, err := net.InterfaceByName(ifaceName)
		if err != nil {
			t.Error(err)
		}
		redundantIPExists, err := RedundantIPCheckWithGratuitousArpOverIface(srcIP, *iface)
		if err != nil {
			t.Error(err)
		}
		t.Logf("redundant ip exists: %t", redundantIPExists)
	}
*/

func TestGoroutinesDoesNotLeak(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")
	SetTimeout(time.Duration(10 * time.Millisecond))

	spawnNumGoroutines := 5
	for i := 0; i < spawnNumGoroutines; i++ {
		_, _, err := Ping(ip)
		if err != ErrTimeout {
			t.Fatalf("timeout error expected, but not received - received err: %v", err)
		}
	}

	ok := make(chan bool, 1)
	go func() {
		for {
			if runtime.NumGoroutine() < spawnNumGoroutines {
				ok <- true
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	select {
	case <-ok:
		// ok
	case <-time.After(30 * time.Second):
		t.Fatalf("timeout waiting for goroutine cleanup - num goroutines: %d",
			runtime.NumGoroutine())
	}
}

func validateInvalidV4AddrErr(t *testing.T, err error) {
	if !strings.Contains(err.Error(), "not a valid v4 Address") {
		t.Errorf("unexpected error: %s", err)
	}
}
