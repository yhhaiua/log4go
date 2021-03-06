package main

import (
	l4g "github.com/yhhaiua/log4go"
	"github.com/yhhaiua/log4go/examples/log"
)

func main() {
	// Load the configuration (isn't this easy?)
	l4g.LoadConfiguration("log4j.xml")

	// And now we're ready!
	l4g.Finest("This will only go to those of you really cool UDP kids!  If you change enabled=true.")
	l4g.Debug("Oh no!  %d + %d = %d!", 2, 2, 2+2)
	l4g.Info("About that time, eh chaps?")

	log.LogRegister()
	select {}
}
