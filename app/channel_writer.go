package app

type ChannelWriter struct {
	sendTo chan<- []byte
}

func NewChannelWriter(sendTo chan<- []byte) *ChannelWriter {
	return &ChannelWriter{sendTo: sendTo}
}

func (w *ChannelWriter) Write(p []byte) (n int, err error) {
	w.sendTo <- p
	return len(p), nil
}
