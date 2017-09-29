package cache

import "testing"
import "fmt"
import "time"

func Test_Cache(t *testing.T) {
	recvChan := make(chan []byte)
	c := NewCache(4096, 500, recvChan)
	defer c.Close()
	go c.Start()

	go func() {
		for {
			select {
			case data := <-recvChan:
				fmt.Printf("receive: %v\n", string(data))
			}
		}
	}()

	for index := 0; index < 20; index++ {
		t := "no_" + fmt.Sprintf("%v", index)
		c.Add([]byte(t))
		time.Sleep(500 * time.Millisecond)
	}
}
