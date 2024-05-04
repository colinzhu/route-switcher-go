package handler

// ChanWriter is a custom io.Writer implementation that writes data to a channel
type ChanWriter struct {
	Channel chan string
}

func NewChanWriter() *ChanWriter {
	return &ChanWriter{Channel: make(chan string, 100)}
}

func (cw *ChanWriter) Write(data []byte) (int, error) {
	select {
	case cw.Channel <- string(data):
		return len(data), nil
	default:
		return 0, nil
	}
}
