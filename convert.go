package xs

import (
	"encoding/binary"
	"math"
	"time"
)

// 均采用大端字节序，事实上这个字节序没有任何用处，
// 只要服务器和客户端约定采用相同的字节序就行
func Uint32ToBytes(v uint32) []byte {
	buf := make([]byte, 4)
	buf[0] = byte(v >> 24)
	buf[1] = byte(v >> 16)
	buf[2] = byte(v >> 8)
	buf[3] = byte(v)
	return buf
}

func Uint16ToBytes(v uint16) []byte {
	buf := make([]byte, 2)
	buf[0] = byte(v >> 8)
	buf[1] = byte(v)
	return buf
}

func BytesToUint32(buf []byte) uint32 {
	v := (uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3]))
	return v
}

func BytesToUint16(buf []byte) uint16 {
	v := (uint16(buf[0])<<8 | uint16(buf[1]))
	return v
}

func TimestampToTimestring(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func BytesToInt32(buf []byte) int32 {
	return int32(binary.BigEndian.Uint32(buf))
}

func ToUInt8(buf []byte, edian byte) uint8 {
	// len(buf) == 1    -->B
	t := uint8(buf[0])
	return t
}

func ToUInt16(buf []byte, edian byte) uint16 {
	// len(buf) == 2    -->H
	t := uint16(buf[0])
	if edian == 62 { // ">"
		t = t<<8 | uint16(buf[1])
	} else if edian == 60 { // "<"
		t = t | uint16(buf[1])<<8
	}

	return t
}

func ToUInt32(buf []byte, edian byte) uint32 {
	// len(buf) == 4    -->I
	t := uint32(buf[0])
	if edian == 62 {
		t = t << 24
		t = t | uint32(buf[1])<<16
		t = t | uint32(buf[2])<<8
		t = t | uint32(buf[3])

	} else if edian == 60 {
		t = t | uint32(buf[1])<<8
		t = t | uint32(buf[2])<<16
		t = t | uint32(buf[3])<<24
	}
	return t
}

func ToUInt64(buf []byte, edian byte) uint64 {
	//len(buf) == 8    -->Q
	t := uint64(buf[0])
	if edian == 62 {
		t = t << 56
		t = t | uint64(buf[1])<<48
		t = t | uint64(buf[2])<<40
		t = t | uint64(buf[3])<<32
		t = t | uint64(buf[4])<<24
		t = t | uint64(buf[5])<<16
		t = t | uint64(buf[6])<<8
		t = t | uint64(buf[7])
	} else if edian == 60 {
		t = t | uint64(buf[1])<<8
		t = t | uint64(buf[2])<<16
		t = t | uint64(buf[3])<<24
		t = t | uint64(buf[4])<<32
		t = t | uint64(buf[5])<<40
		t = t | uint64(buf[6])<<48
		t = t | uint64(buf[7])<<56
	}
	return t
}

func PutUInt8(num uint8, buf []byte, edian byte) {
	// len(buf) == 1
	buf[0] = byte(num)
}

func PutUInt16(num uint16, buf []byte, edian byte) {
	// len(buf) == 2
	buf[0] = byte(num >> 8)
	buf[1] = byte(num)
	if edian == 62 { // ">"

	} else if edian == 60 { // "<"
		buf[0] ^= buf[1]
		buf[1] ^= buf[0]
		buf[0] ^= buf[1]
	}
}

func PutUInt32(num uint32, buf []byte, edian byte) {
	// len(buf) == 4
	buf[0] = byte(num >> 24)
	buf[1] = byte(num >> 16)
	buf[2] = byte(num >> 8)
	buf[3] = byte(num)
	if edian == 62 {

	} else if edian == 60 {
		buf[0] ^= buf[3]
		buf[3] ^= buf[0]
		buf[0] ^= buf[3]

		buf[1] ^= buf[2]
		buf[2] ^= buf[1]
		buf[1] ^= buf[2]
	}
}

func PutUInt64(num uint64, buf []byte, edian byte) {
	// len(buf) == 8
	if edian == 62 {
		buf[0] = byte(num >> 56)
		buf[1] = byte(num >> 48)
		buf[2] = byte(num >> 40)
		buf[3] = byte(num >> 32)
		buf[4] = byte(num >> 24)
		buf[5] = byte(num >> 16)
		buf[6] = byte(num >> 8)
		buf[7] = byte(num)
	} else if edian == 60 {
		buf[0] = byte(num)
		buf[1] = byte(num >> 8)
		buf[2] = byte(num >> 16)
		buf[3] = byte(num >> 24)
		buf[4] = byte(num >> 32)
		buf[5] = byte(num >> 40)
		buf[6] = byte(num >> 48)
		buf[7] = byte(num >> 56)
	}
}

/////////////////////////////小端字节序////////////////////////////////////////////////////////

func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func ByteToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}
