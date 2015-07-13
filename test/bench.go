package main

import (
    "encoding/binary"
    "fmt"
    "net"
)

const (
    M    = 4
    N    = 1000000
    Addr = "0.0.0.0:8000"
)

func Bench() {
    var logChan chan uint64 = make(chan uint64)
    var signal chan int = make(chan int)
    // log
    go func() {
        for {
            id := <-logChan
            fmt.Println(id)
        }
    }()

    for i := 0; i < M; i++ {
        go func() {
            for j := 0; j < N/M; j++ {
                var buf [8]byte
                conn, err := net.Dial("tcp", Addr)
                if err != nil {
                    fmt.Println(err)
                    return
                }
                conn.Write([]byte{109, 101})
                conn.Read(buf[:])
                conn.Close()
                id, _ := binary.Uvarint(buf[:])
                logChan <- id
            }
        }()
        signal <- i
    }

    for {
        n := <-signal
        if n == M-1 {
            return
        }
    }
}

func main() {
    Bench()
}
