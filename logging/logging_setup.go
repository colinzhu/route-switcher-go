package logging

import (
	"io"
	"log"
	"os"
)

func InitLogging() (*os.File, chan string) {
	chanWriter := NewChanWriter()

	file, err := os.OpenFile("route-switcher.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	multiWriter := io.MultiWriter(os.Stdout, file, chanWriter)

	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("log initialized")
	return file, chanWriter.Channel
}
