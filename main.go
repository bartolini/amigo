package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bartolini/amigo/broadcast"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	broadcaster := broadcast.NewBroadcaster(
		ctx,
	).Listeners(
		makeListener("John"),
		makeListener("Amy"),
		makeListener("Deco"),
		makeListener("Rian"),
		makeListener("Kathy"),
	).Filters(
		keepStringMessagesOnly,
	)
	defer broadcaster.Wait()

	fmt.Println("===> Broadcasting some messages...")

	broadcaster.Broadcast("Hello world!").
		Broadcast("Hello kitty!").
		Broadcast(0xdeadbeef).
		Broadcast(0xcaffebabe).
		Broadcast("Minions vs Oompa Loompas").
		Broadcast("The Last V8").
		Broadcast("Whatever...")

	time.Sleep(time.Millisecond)

	fmt.Println("===> Stopping everything now...")

	cancel()

	fmt.Println("===> Exiting.")
}

func keepStringMessagesOnly(msg broadcast.Message) bool {
	_, ok := msg.(string)
	return !ok
}

func makeListener(name string) broadcast.InitFunc {
	listener := func(ctx context.Context, waitgroup *sync.WaitGroup) chan<- broadcast.Message {
		ch := make(chan broadcast.Message)
		go func() {
			defer waitgroup.Done()
			for {
				select {
				case msg := <-ch:
					fmt.Printf("%s:\t%v\n", name, msg)
				case <-ctx.Done():
					fmt.Printf("%s:\tGoodbye!\n", name)
					return
				}
			}
		}()
		return ch
	}
	return listener
}
