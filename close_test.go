package closer

import (
	"syscall"
	"testing"
	"time"
)

var got string
var localChain = Chain(
	func() { got += "1" },
	func() { got += "2" },
	func() { got += "3" },
)

func TestCloser(t *testing.T) {
	var done = make(chan bool)
	wanted := "123"
	Closer(done, localChain)
	time.Sleep(time.Second)

	select {
	case <-time.After(time.Second):
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}

	// select {
	// case <-time.After(time.Second):
	// 	_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	// }

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	}

	if wanted != got {
		t.Fatalf("closer not executed in order:\nwanted:%s\ngot:%s\n", wanted, got)
	}
}
