package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {

	myArr := [...]int{1, 2, 3}
	log.Printf("Type %[1]T, values %[1]v\n", myArr)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	generator := func(dataItem interface{}, stream chan interface{}) {
		for {
			select {
			case <-ctx.Done():
				return
			case stream <- dataItem:
			}
		}
	}

	func2 := genericFunc
	func3 := genericFunc

	infiniteAppels := make(chan interface{})
	// infiniteAppels <- myArr
	go generator(myArr, infiniteAppels)

	infiniteIceCreams := make(chan interface{})
	go generator("ice cream", infiniteIceCreams)

	infiniteHappiness := make(chan interface{})
	go generator("happiness", infiniteHappiness)

	wg.Add(1)
	go func1(ctx, &wg, infiniteAppels)
	wg.Add(1)
	go func2(ctx, &wg, infiniteIceCreams)
	wg.Add(1)
	go func3(ctx, &wg, infiniteHappiness)
	wg.Wait()
	// var cnt float64
	// for cnt < time.Second*3 {

	// }
}

func func1(ctx context.Context, parentWG *sync.WaitGroup, stream <-chan interface{}) {
	defer parentWG.Done()
	var wg sync.WaitGroup

	doWork := func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-stream:
				if !ok {
					log.Println("channel is closed...")
					return
				}
				fmt.Println(d)
			}
		}
	}

	newCtx, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go doWork(newCtx)
	}
	wg.Wait()
}

func genericFunc(ctx context.Context, wg *sync.WaitGroup, stream <-chan interface{}) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Println("chan is closed...")
			return
		case d, ok := <-stream:
			if !ok {
				log.Println("channel is closed")
				return
			}
			fmt.Println(d)
		}
	}
}
