package log

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type errMsg struct {
	t string
	s string
}

var mq chan *errMsg

func InitLog(f string) {
	mq = make(chan *errMsg, 10)

	go func() {
		var file *os.File = nil
		var bufWrite *bufio.Writer = nil
		for {
			select {
			case e := <-mq:
				if file == nil {
					var err error
					file, err = os.OpenFile(f, os.O_APPEND|os.O_CREATE, 0666)
					if err != nil {
						Print(err.Error())
					}
					bufWrite = bufio.NewWriter(file)
				}
				// 写日志到带buffer的writer中
				fmt.Fprintf(bufWrite, "[%s] %s\n", e.t, e.s)
			default:
				// 没数据了，写入文件中
				if bufWrite != nil {
					bufWrite.Flush()
					file.Close()
					file, bufWrite = nil, nil
				}
			}
		}
	}()

}

func Print(msg string) {
	mq <- &errMsg{
		t: time.Now().Format("2006-01-02 15:04:05.000"),
		s: msg,
	}
}
