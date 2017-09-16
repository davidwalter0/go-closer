// deterministic signal catching closer
//
// Catch one of 3 signals and enable a chain of methods that must be
// called. If there are blocking methods, perhaps consider calls with
// time outs
//
// func main() {
// 	closer.Closer(done, closer.SampleChain)
// 	closer.Closer(done, closer.NoOp) // do nothing
// 	go run()
// 	<-done
// 	fmt.Println("main after Close")
// }

package closer

import (
	"fmt"
	trace "github.com/davidwalter0/tracer"
	"os"
	"os/signal"
	"syscall"
)

var tracer = trace.New()

var enable = true
var detail = false

// Link is the func type for a chain
type Link func()

// var chain = NoOp
var SampleChain = Chain(func() { fmt.Print(1) }, func() { fmt.Print(2) }, func() { fmt.Print(3) })

var NoOp = (func() {
	defer tracer.Detailed(detail).Enable(enable).ScopedTrace()()
})

// Chain is a function taking a function list of calls to execute
func Chain(chain ...Link) func() {
	if len(chain) > 1 {
		current := chain[0]
		next := Chain(chain[1:]...)
		return func() {
			defer tracer.Detailed(detail).Enable(enable).ScopedTrace()()
			current()
			next()
		}
	} else if len(chain) == 1 {
		defer tracer.Detailed(detail).Enable(enable).ScopedTrace()()
		current := chain[0]
		return func() {
			defer tracer.Detailed(detail).Enable(enable).ScopedTrace()()
			current()
		}
	}
	return NoOp
}

// Closer handles signals and closes
func Closer(done chan bool, chain Link) {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		fmt.Println("Closer before signal")
		// <-c
		s := <-c
		fmt.Println("Closer after signal", s)
		fmt.Println("Closer after Close")
		chain()
		done <- true
	}()
}
