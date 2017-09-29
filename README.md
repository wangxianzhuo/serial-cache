## Install

```shell
go get -u github.com/wangxianzhuo/serial-cache
```

## Example

```golang
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
```