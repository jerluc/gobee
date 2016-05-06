package gobee

import (
	"bytes"
	"encoding/binary"
)

// Computes an 8-bit checksum of XBee frame data
func Checksum(data []byte) byte {
	var sum byte = 0
	for _, b := range data {
		sum += b
	}
	return 0xFF - sum
}

// Verifies the computed checksum of XBee frame data
// with the checksum received over the wire. It is
// encouraged that you discard any XBee frame which
// does not pass this check
func VerifyChecksum(data []byte, receivedChecksum byte) bool {
	return Checksum(data) == receivedChecksum
}

// Converts an unsigned 16-bit integer to a 2-byte
// big Endian array
func Uint16ToBytes(n int) []byte {
	buff := bytes.NewBuffer([]byte{})
	binary.Write(buff, binary.BigEndian, uint16(n))
	return buff.Bytes()
}

// Converts a 2-byte big Endian array into an
// unsigned 16-bit integer
func BytesToUint16(b []byte) int {
	var n uint16
	buff := bytes.NewBuffer(b)
	binary.Read(buff, binary.BigEndian, &n)
	return int(n)
}

// Takes a variable number of single bytes or byte
// arrays and concatenates them into a single byte
// array
func PackBytes(parts ...interface{}) []byte {
	var packet []byte
	for _, part := range parts {
		switch b := part.(type) {
		case byte:
			packet = append(packet, b)
		case []byte:
			packet = append(packet, b...)
		}
	}
	return packet
}
