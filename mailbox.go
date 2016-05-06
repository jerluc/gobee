package gobee

import (
	"io"
)

type Mailbox struct {
	fb *FrameBuffer
}

func NewMailbox(rw io.ReadWriter) *Mailbox {
	return &Mailbox{NewFrameBuffer(rw)}
}

func (m *Mailbox) Inbox() <-chan Frame {
	inbox := make(chan Frame)
	go func() {
		for {
			incoming := m.fb.ReadFrame()
			inbox <- incoming
		}
	}()
	return inbox
}

func (m *Mailbox) Outbox() chan<- Frame {
	outbox := make(chan Frame)
	go func() {
		for {
			f := <-outbox
			m.fb.WriteFrame(f)
		}
	}()
	return outbox
}
