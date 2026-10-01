package servers

import (
	"net"
	"testing"
)

func TestListenLoopback(t *testing.T) {
	t.Run("skips_busy_port", func(t *testing.T) {
		busy, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = busy.Close() })
		start := busy.Addr().(*net.TCPAddr).Port

		ln, port, err := listenLoopback("127.0.0.1", start, 20)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = ln.Close() })
		if port == start {
			t.Fatalf("bound busy port %d", port)
		}
		if port < start || port >= start+20 {
			t.Fatalf("port %d outside %d-%d", port, start, start+19)
		}
	})

	t.Run("fails_after_attempts", func(t *testing.T) {
		busy, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = busy.Close() })
		start := busy.Addr().(*net.TCPAddr).Port

		_, _, err = listenLoopback("127.0.0.1", start, 1)
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
