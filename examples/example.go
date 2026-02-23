package main

import (
	"fmt"
	"time"
)

func main() {
	// Greet the world
	message := "Hello, World!"
	
	for i := 0; i < 3; i++ {
		fmt.Printf("%d: %s\n", i+1, message)
		time.Sleep(100 * time.Millisecond)
	}
	
	fmt.Println("Done! 🎉")
}
