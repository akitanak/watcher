package app

import "fmt"

type StdoutPrinter struct {
	receiver <-chan []byte
}

func NewStdoutPrinter(receiver <-chan []byte) *StdoutPrinter {
	return &StdoutPrinter{receiver: receiver}
}

func (p *StdoutPrinter) Print() {
	for {
		select {
		case msg := <-p.receiver:
			fmt.Printf("%s", string(msg))
		}
	}
}
