package utils

import (
	"fmt"
	"time"
)

func timeStamp() string {
	now := time.Now().Format(time.StampMilli)

	return fmt.Sprintf("[%s]", Yellow(now))
}

func Log(str string) {
	fmt.Printf("%s %s\n", timeStamp(), str)
}

func LogError(error string) {
	str := fmt.Sprintf("[%s] %s", Red("Error"), error)
	Log(str)
}

func LogWarn(warn string) {
	str := fmt.Sprintf("[%s] %s", Yellow("Warn"), warn)
	Log(str)
}

func LogSuccess(h, msg string) {
	str := fmt.Sprintf("[%s] %s", Green(h), msg)
	Log(str)
}
