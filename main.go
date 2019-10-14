//go:generate protoc -I proto/ --go_out=models proto/models.proto
package main

import (
	"os"
	"os/signal"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func main() {

	managerContext := actor.EmptyRootContext

	managerProps := actor.PropsFromProducer(func() actor.Actor {
		return new(AppActor)
	})
	manager := managerContext.Spawn(managerProps)

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)

	<-signals
	managerContext.Stop(manager)
}
