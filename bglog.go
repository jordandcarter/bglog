package bglog

import (
  "bufio"
  "fmt"
  "os"
)

type BgLog struct {
  buffer chan string
  io *bufio.Writer
}

func NewBgLog(f *os.File, size int) *BgLog {
  bl := BgLog{buffer: make(chan string, size), io: bufio.NewWriter(f)}
  go bl.ProcessLog()
  return &bl
}

func (bl *BgLog) ProcessLog() {
	var msg string
	var count int
	controlChan := make(chan int, 2)
	for {
		select {
		case msg = <-bl.buffer:
			count++
			bl.io.Write([]byte(msg))
			if len(controlChan) == 0 {
				controlChan <- 0
			}
		case <-controlChan:
			if count > 1 {
				bl.io.Write([]byte(fmt.Sprintln("Flush: ", count)))
			}
			count = 0
			bl.io.Flush()
		}
	}
}

func (bl *BgLog) Add(msg string) {
  bl.buffer <- msg
}

