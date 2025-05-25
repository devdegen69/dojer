package downloader

import (
	u "dojer/utils"
	"fmt"
	"io"
	"strings"
	"sync"
)

var (
	pr, pw       = io.Pipe()
	prg          chan string
	lineCounter  = 1
	progressList = map[string]Progress{}
	mu           sync.Mutex
)

func GetPipe() (*io.PipeReader, *io.PipeWriter) {
	return pr, pw
}

func GetProgress() chan string {
	return prg
}

type Progress struct {
	Name     string
	Message  string
	Progress string
	Line     int
}

func (p *Progress) getMessage() string {
	msg := fmt.Sprintf("[%s] %s - progress: %s", u.Yellow(p.Name), p.Message, u.White(p.Progress))
	spaces := strings.Repeat(" ", len(msg)+20)
	return fmt.Sprintf("\r%s\r%s", spaces, msg)
}

func (p *Progress) setMessage(newMsg string) {
	p.Message = newMsg
	progressList[p.Name] = *p
	update()
}

func (p *Progress) setProgress(newProgress string) {
	p.Progress = newProgress
	progressList[p.Name] = *p
	update()
}

func isProgressListEmpty() bool {
	return len(progressList) == 0
}

func createNewProgress(name string, message string) *Progress {
	mu.Lock()
	defer mu.Unlock()
	var p Progress
	p.Name = name
	p.Message = message
	p.Line = lineCounter
	p.Progress = "[0/0]"
	progressList[name] = p
	lineCounter += 1
	return &p
}

func update() {
	for _, p := range progressList {
		fmt.Fprint(pw, "\033[u")
		goDown(p.Line)
		fmt.Fprint(pw, p.getMessage())
	}
}

// func goUp(times int) {
// 	fmt.Fprintf(pw, "\033[%dA\r", times)
// }

func goDown(times int) {
	fmt.Fprintf(pw, "\033[%dB\r", times)
}
