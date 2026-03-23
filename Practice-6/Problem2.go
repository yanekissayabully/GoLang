// Почему counter++ не даёт 1000??

// counter++ это не одна инструкция а три: read, increment и write. 
// Две горутины могут прочитать одинаковое значение, обе добавят 1, обе запишут наше новое число.По итогу вместо +2 получается +1.
// Это как раз таки race condition.

// package main
// import (
// 	"fmt"
// 	"sync"
// )
// func main() {
// 	var counter int
// 	var wg sync.WaitGroup
// 	for i := 0; i < 1000; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			counter++
// 		}()
// 	}
// 	wg.Wait()
// 	fmt.Println(counter)
// }



// так же через mutex

// package main

// import (
// 	"fmt"
// 	"sync"
// )

// func main() {
// 	var counter int
// 	var wg sync.WaitGroup
// 	var mu sync.Mutex

// 	for i := 0; i < 1000; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			mu.Lock()
// 			counter++
// 			mu.Unlock()
// 		}()
// 	}

// 	wg.Wait()
// 	fmt.Println(counter)
// }



// через sync/atomic

// package main

// import (
// 	"fmt"
// 	"sync"
// 	"sync/atomic"
// )

// func main() {
// 	var counter int64
// 	var wg sync.WaitGroup

// 	for i := 0; i < 1000; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			atomic.AddInt64(&counter, 1) 
// 		}()
// 	}

// 	wg.Wait()
// 	fmt.Println(atomic.LoadInt64(&counter)) 
// }