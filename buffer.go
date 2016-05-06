package gobee

import (
	"io"
)

// A generic full-duplex frame buffer that allows for reading/
// writing XBee frames to the underlying serial device (or any
// other kind of io.ReadWriter
type FrameBuffer struct{
	rw io.ReadWriter
}

// Creates a new FrameBuffer object
func NewFrameBuffer(rw io.ReadWriter) *FrameBuffer {
	fb := &FrameBuffer{rw}
	return fb
}

// Writes a Frame object to the underlying serial device. Note that
// this does *not* buffer writes nor does it do packet fragmentation
func (fb *FrameBuffer) WriteFrame(frame Frame) (n int, err error) {
	frameBytes := PackBytes(
		FrameHeader,
		Uint16ToBytes(uint(len(frame.FrameData()))),
		frame.FrameData(),
		Checksum(frame.FrameData()),
	)
	return fb.rw.Write(frameBytes)
}

// Reads a single byte from the underlying serial device.
func (fb *FrameBuffer) readByte() byte {
	return fb.readBytes(1)[0]
}

// Reads n bytes from the underlying serial device. If less than n
// bytes were available to read, this method will block until the
// full n are available to read
func (fb *FrameBuffer) readBytes(n int) []byte {
	b := make([]byte, n)
	readLen, err := fb.rw.Read(b)
	if err != nil {
		panic(err)
	}
	if readLen < n {
		slop := fb.readBytes(n - readLen)
		return PackBytes(b[:readLen], slop)
	}
	return b
}

// Reads a single XBee frame from the underlying serial device.
// Note that this method will block on reads until the full frame
// has been consumed. Additionally, this method will skip any bytes
// received between the presumed end of the previous frame and the
// beginning (0x7E) of the next
func (fb *FrameBuffer) ReadFrame() Frame {
	for {
		if fb.readByte() == FrameHeader {
			lengthBytes := fb.readBytes(2)
			length := BytesToUint16(lengthBytes)
			frameData := fb.readBytes(int(length))
			checksum := fb.readByte()
			if VerifyChecksum(frameData, checksum) {
				return BuildFrame(frameData)
			}
		}
	}
}

