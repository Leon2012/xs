package xs

import (
	"errors"
	"fmt"
	_ "strconv"
)

type XSCommand struct {
	XSComponent
	Cmd  int
	Arg1 int
	Arg2 int
	Buf  string
	Buf1 string
}

// func NewCommand(c uint8, a1 uint8, a2 uint8, b1 string, b2 string) *XSCommand {
// 	return &XSCommand{
// 		Cmd:  c,
// 		Arg1: a1,
// 		Arg2: a2,
// 		Buf:  b1,
// 		Buf1: b2,
// 	}
// }

func NewCommand(c int, a1 int, a2 int, b1 string, b2 string) *XSCommand {
	return &XSCommand{
		Cmd:  c,
		Arg1: a1,
		Arg2: a2,
		Buf:  b1,
		Buf1: b2,
	}
}

func NewCommand1(c int) *XSCommand {
	a1 := 0
	a2 := 0
	b1 := ""
	b2 := ""
	return NewCommand(c, a1, a2, b1, b2)
}

func NewCommandWithMap(cm map[string]interface{}) *XSCommand {
	cmd := XS_CMD_NONE
	arg1, arg2 := 0, 0
	buf, buf1 := "", ""

	for key, value := range cm {
		if key == "cmd" {
			cmd = value.(int) //interface{}转int
		} else if key == "arg1" {
			arg1 = value.(int)
		} else if key == "arg2" {
			arg2 = value.(int)
		} else if key == "buf" {
			//interface{}转string
			if str, ok := value.(string); ok {
				buf = str
			} else {
				buf = ""
			}
		} else if key == "buf1" {
			if str, ok := value.(string); ok {
				buf1 = str
			} else {
				buf1 = ""
			}
		}
	}

	return NewCommand(cmd, arg1, arg2, buf, buf1)
}

func (x *XSCommand) ToBytes() ([]byte, error) {
	if len(x.Buf1) > 0xff {
		buf1 := x.Buf1[0:0xff]
		x.Buf1 = buf1
	}
	bufBytes := []byte(x.Buf)
	buf1Bytes := []byte(x.Buf1)

	//pack('CCCCI', $this->cmd, $this->arg1, $this->arg2, strlen($this->buf1), strlen($this->buf)) . $this->buf . $this->buf1;
	//pack("CCCCI") = uchar(1) + uchar(1) + uchar(1) + uchar(1) + uint(4) = 8
	bufLen := 8 + len(bufBytes) + len(buf1Bytes)
	var buf []byte
	buf = make([]byte, bufLen)

	var cmd uint8
	var arg1 uint8
	var arg2 uint8

	if x.Cmd > 0xff || x.Cmd < 0x00 {
		return nil, errors.New("cmd is large")
	}

	if x.Arg1 > 0xff || x.Arg1 < 0x00 {
		return nil, errors.New("Arg1 is large")
	}

	if x.Arg2 > 0xff || x.Arg2 < 0x00 {
		return nil, errors.New("Arg2 is large")
	}

	cmd = uint8(x.Cmd)
	arg1 = uint8(x.Arg1)
	arg2 = uint8(x.Arg2)

	buf[0] = byte(cmd)
	buf[1] = byte(arg1)
	buf[2] = byte(arg2)

	buf1Len := len(x.Buf1)
	buf1LenByte := byte(buf1Len)
	buf[3] = buf1LenByte

	bufLen = len(x.Buf)
	bufLenBytes := make([]byte, 4)
	//copy(buf[4:8], Uint32ToBytes(uint32(bufLen)))
	PutUInt32(uint32(bufLen), bufLenBytes, 60)
	copy(buf[4:8], bufLenBytes)

	copy(buf[8:(8+len(bufBytes))], bufBytes)

	idx := 8 + len(bufBytes)
	copy(buf[idx:(idx+len(buf1Bytes))], buf1Bytes)

	fmt.Println(buf)

	return buf, nil
}

func (x *XSCommand) GetArg() int {
	var i int
	i = x.Arg2 | (x.Arg1 << 8)
	return i
}

func (x *XSCommand) SetArg(arg int) {
	x.Arg1 = (arg >> 8) & 0xff
	x.Arg2 = arg & 0xff
}
