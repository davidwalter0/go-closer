// deterministic signal catching closer when signals will be used for
// shutdown
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
// }
//
// To run with a timeout deadline

// func main() {
// 	closer.Closer(done, closer.SampleChain)
//  go run()
//     select {
//     case <-done:
//     case <-time.After(2 * time.Second):
//       _ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
//     }
//   }
// }

package closer

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Link is the func type for a chain
type Link func()

// SampleChain for application validation test
var SampleChain = Chain(func() { fmt.Print(1) }, func() { fmt.Print(2) }, func() { fmt.Print(3) })

// NoOp terminal empty function in chain or placeholder when needed
// for compilation
var NoOp = (func() {
})

// Chain is a function taking a function list of calls to execute
func Chain(chain ...Link) func() {
	if len(chain) > 1 {
		current := chain[0]
		next := Chain(chain[1:]...)
		return func() {
			current()
			next()
		}
	} else if len(chain) == 1 {
		current := chain[0]
		return func() {
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
		<-c
		chain()
		done <- true
	}()
}
