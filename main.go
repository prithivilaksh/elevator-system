package main

import (
	"fmt"
	"sync"
)

type test struct {
	mu sync.Mutex
	cnt int
}

func (t *test) increment() {
	t.mu.Lock()
	defer t.mu.Unlock()
	fmt.Println("Incrementing count")
	t.cnt++
}

func (t *test) decrement() {
	t.mu.Lock()
	defer t.mu.Unlock()
		fmt.Println("Decrementing count")
	t.cnt--
}

func main() {
	fmt.Println("Hello, World!")
	obj:= &test{
		cnt: 0,
	}

	var wg sync.WaitGroup
	for range 100000{
		wg.Add(1)
		wg.Add(1)

		go func(){
			defer wg.Done()
			obj.increment()	
		} ()

		go func(){
			defer wg.Done()
			obj.decrement()	
		} ()
		
	}

	wg.Wait()
	fmt.Println("Final count:", obj.cnt)

}
