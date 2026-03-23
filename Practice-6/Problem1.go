// package main
// import (
// 	"fmt"
// )
// func main() {
// 	unsafeMap := make(map[string]int)
// 	for i := 0; i < 100; i++ {
// 		go func(key int) {
// 			unsafeMap["key"] = key
// 		}(i)
// 	}
// 	value := unsafeMap["key"]
// 	fmt.Printf("Value: %d\n", value)
// }


// 1) we can use syncmap

// package main

// import (
// 	"fmt"
// 	"sync"
// )

// func main() {
// 	var safeMap sync.Map
// 	var wg sync.WaitGroup

// 	for i := 0; i < 100; i++ {
// 		wg.Add(1)
// 		go func(key int) {
// 			defer wg.Done()
// 			safeMap.Store("key", key)
// 		}(i)
// 	}

// 	wg.Wait()

// 	value, ok := safeMap.Load("key")
// 	if ok {
// 		fmt.Printf("Value: %v\n", value)
// 	}
// }


// 2) We can do it with sync.RWMutex

// package main

// import (
// 	"fmt"
// 	"sync"
// )

// func main() {
// 	unsafeMap := make(map[string]int)
// 	var mu sync.RWMutex
// 	var wg sync.WaitGroup

// 	for i := 0; i < 100; i++ {
// 		wg.Add(1)
// 		go func(key int) {
// 			defer wg.Done()
// 			mu.Lock()        
// 			unsafeMap["key"] = key
// 			mu.Unlock()       
// 		}(i)
// 	}

// 	wg.Wait()

// 	mu.RLock()              
// 	value := unsafeMap["key"]
// 	mu.RUnlock()

// 	fmt.Printf("Value: %d\n", value)
// }