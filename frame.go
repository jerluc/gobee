package gobee

import (
	"fmt"
)

// See https://docs.digi.com/display/WirelessConnectivityKit/Frame+types+in+detail
// for more exact details on each frame type's packet structure

const (
	// Frame start delimiter
	FrameHeader byte = 0x7E
	// TX frame type (64-bit address)
	Tx64FrameType byte = 0x00
	// TX status frame type
	TxStatusFrameType byte = 0x89
	// RX frame type (64-bit address)
	Rx64FrameType byte = 0x80
)

// Broadcast address for use by 64-bit TX frames
var BroadcastAddress = []byte{ 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF }

// Represents a generic XBee frame type
type Frame interface {
	FrameData() []byte
}

// Represents a "generic" (unimplemented) XBee frame
type GenericFrame struct {
	Type byte
	RawData []byte
}

// Constructs a new "generic" frame object from the provided frame data
func BuildGenericFrame(fd []byte) Frame {
	return &GenericFrame{
		Type: fd[0],
		RawData: fd[1:],
	}
}

func (g *GenericFrame) FrameData() []byte {
	return PackBytes(g.Type, g.RawData)
}

func (g *GenericFrame) String() string {
	return fmt.Sprintf("GenericFrame[Type: %X, RawData: %X]", g.Type, g.RawData)
}

// Represents a 64-bit addressed TX frame
type Tx64Frame struct {
	ID          byte
	Destination []byte
	Options     byte
	Data        []byte
}

// Constructs a new TX64 frame object from the provided frame data
func BuildTx64Frame(fd []byte) Frame {
	return &Tx64Frame{
		ID: fd[1],
		Destination: fd[2:10],
		Options: fd[10],
		Data: fd[11:],
	}
}

func (tx *Tx64Frame) FrameData() []byte {
	return PackBytes(
		Tx64FrameType,
		tx.ID,
		tx.Destination,
		tx.Options,
		tx.Data,
	)
}

func (tx *Tx64Frame) String() string {
	return fmt.Sprintf("TX64[ID: %X, Destination: %X, Options: %X, Data: %s]",
		tx.ID, tx.Destination, tx.Options, string(tx.Data))
}

// Represents a TX status frame
type TxStatusFrame struct {
	ID byte
	Status byte
}

// Constructs a new TX status frame object from the provided frame data
func BuildTxStatusFrame(fd []byte) Frame {
	return &TxStatusFrame{
		ID: fd[1],
		Status: fd[2],
	}
}

func (tx *TxStatusFrame) FrameData() []byte {
	return PackBytes(TxStatusFrameType, tx.ID, tx.Status)
}

func (tx *TxStatusFrame) String() string {
	return fmt.Sprintf("TX-Status[ID: %X, Status: %X]", tx.ID, tx.Status)
}

// Represents a 64-bit addressed RX frame
type Rx64Frame struct {
	Source  []byte
	RSSI    byte
	Options byte
	Data    []byte
}

// Constructs a new RX64 frame from the provided frame data
func BuildRx64Frame(fd []byte) Frame {
	return &Rx64Frame{
		Source: fd[1:9],
		RSSI: fd[9],
		Options: fd[10],
		Data: fd[11:],
	}
}

func (rx *Rx64Frame) FrameData() []byte {
	return PackBytes(
		Rx64FrameType,
		rx.Source,
		rx.RSSI,
		rx.Options,
		rx.Data,
	)
}

func (rx *Rx64Frame) String() string {
	return fmt.Sprintf("RX[Source: %X, RSSI: %X, Options: %X, Data: %s]",
		rx.Source, rx.RSSI, rx.Options, string(rx.Data))
}

// Given raw frame data, this function attempts to construct the appropriate
// frame object by inspecting the first byte of the frame data (presumed to
// be the frame type byte). If the frame type cannot be identified by this
// software, a "generic" frame object is constructed preserving the raw data
// for other implementations to handle
func BuildFrame(fd []byte) Frame {
	if len(fd) == 0 {
		return nil
	}
	frameType := fd[0]
	switch frameType {
	case Tx64FrameType:
		return BuildTx64Frame(fd)
	case TxStatusFrameType:
		return BuildTxStatusFrame(fd)
	case Rx64FrameType:
		return BuildRx64Frame(fd)
	default:
		return BuildGenericFrame(fd)
	}
	return nil
}

