package cache

import (
	"fmt"
	"sync"
	"time"
)

// Cache use ReceiveChan to output data
// FlushInterval default 4096
type Cache struct {
	Data          []byte
	Size          int
	FlushInterval int
	Capacity      int
	ReceiveChan   chan []byte
	Quit          chan bool
	Mux           sync.Mutex
}

// NewCache create a Cache
func NewCache(capacity, interval int, recvChan chan []byte) *Cache {
	if capacity < 1 {
		capacity = 4096
	}
	if interval < 0 {
		interval = 0
	}
	return &Cache{
		Data:          make([]byte, capacity, capacity),
		Size:          0,
		FlushInterval: interval,
		Capacity:      capacity,
		ReceiveChan:   recvChan,
		Quit:          make(chan bool),
	}
}

func (c Cache) String() string {
	return fmt.Sprintf("Data: [%X]\tSize: %v\tFlushInterval: %v\tCapacity: %v", c.Data[:c.Size], c.Size, c.FlushInterval, c.Capacity)
}

// Start use a time.Ticker to cycle flushing data.
// Can maual to do this
func (c *Cache) Start() {
	tick := time.NewTicker(time.Duration(c.FlushInterval) * time.Millisecond)
	for {
		select {
		case <-c.Quit:
			close(c.ReceiveChan)
			return
		default:
		}
		select {
		case <-tick.C:
			res, n := c.Flush()
			if n == 0 {
				continue
			}
			c.ReceiveChan <- res
		}
	}
}

// Flush data in Cache
func (c *Cache) Flush() ([]byte, int) {
	c.Mux.Lock()
	defer c.Mux.Unlock()

	result := make([]byte, 0)
	for index := 0; index < c.Size; index++ {
		result = append(result, c.Data[index])
	}

	c.Size = 0
	return result, len(result)
}

// Add data to the Cache
func (c *Cache) Add(b []byte) error {
	addtionSize := len(b)
	c.Mux.Lock()
	defer c.Mux.Unlock()

	addNumber := addtionSize
	if c.Capacity < c.Size+addtionSize {
		addNumber = c.Capacity - c.Size
	}

	for index := 0; index < addNumber; index++ {
		c.Data[c.Size+index] = b[index]
	}
	c.Size += addNumber
	return nil
}

// Close ReceiveChan
func (c *Cache) Close() {
	close(c.Quit)
}
