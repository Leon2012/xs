// server
package xs

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"
)

const (
	FILE   = 0x01
	BROKEN = 0x02
)

type XSServer struct {
	XSComponent
	Xs         *XS
	address    string
	conn       net.Conn
	flag       int
	project    string
	sendBuffer []byte
}

func NewServer(addr string, xs *XS) (*XSServer, error) {
	serv := &XSServer{}
	serv.Xs = xs
	serv.conn = nil
	err := serv.Open(addr)
	if err != nil {
		return nil, err
	}
	return serv, nil
}

func (x *XSServer) GetConnString() string {
	return x.address
}

func (x *XSServer) GetSocket() net.Conn {
	return x.conn
}

func (x *XSServer) GetProject() string {
	return x.project
}

func (x *XSServer) SetProject(name string, home string) {
	if name != x.project {
		cmd := NewCommand(XS_CMD_USE, 0, 0, name, home)
		x.ExecCommand(cmd, XS_CMD_OK_PROJECT, XS_CMD_OK)
		x.project = name
	}
}

func (x *XSServer) Open(addr string) error {
	x.Close(false)

	x.address = addr
	x.flag = BROKEN
	x.project = ""
	x.sendBuffer = []byte{}
	err := x.connect()
	x.flag ^= BROKEN
	return err
}

func (x *XSServer) Reopen(force bool) {
	if (x.flag&BROKEN) == 1 || force == true {
		x.Open(x.address)
	}
}

func (x *XSServer) Close(ioerr bool) {
	if x.conn != nil && (x.flag&BROKEN) == 0 {
		if !ioerr && len(x.sendBuffer) > 0 {
			buf := x.sendBuffer[0:]
			x.write(buf)
			x.sendBuffer = []byte{}
		}
		if !ioerr && (x.flag&FILE) == 0 {
			cmd := NewCommand1(XS_CMD_QUIT)
			fmt.Println("quit command")
			fmt.Println(cmd)
			bytes, _ := cmd.ToBytes()
			x.conn.Write(bytes)
		}

		x.conn.Close()
		x.flag |= BROKEN
	}
}

func (x *XSServer) SetTimeout(sec int) {
	cmd := NewCommand1(XS_CMD_TIMEOUT)
	cmd.SetArg(sec)
	x.ExecCommand(cmd, XS_CMD_OK_TIMEOUT_SET, XS_CMD_OK)
}

func (x *XSServer) ExecCommand1(cmd *XSCommand) (*XSCommand, error) {
	return x.ExecCommand(cmd, XS_CMD_NONE, XS_CMD_OK)
}

func (x *XSServer) ExecCommand(cmd *XSCommand, resArg int, resCmd int) (*XSCommand, error) {
	if (cmd.Cmd & 0x80) == 1 {
		bytes, _ := cmd.ToBytes()
		x.appendSendBuffer(bytes)
		return nil, nil
	}

	bytes, _ := cmd.ToBytes()
	x.appendSendBuffer(bytes)
	buf := x.sendBuffer[0:]
	x.sendBuffer = []byte{}

	err := x.write(buf)
	if err != nil {
		return nil, err
	}

	if (x.flag & FILE) == 1 {
		return nil, nil
	}

	res := x.GetRespond()
	if res == nil {
		return nil, errors.New("exception : get respond error ")
	}
	//fmt.Println(res)

	if res.Cmd == XS_CMD_ERR && resCmd != XS_CMD_ERR {
		return nil, errors.New("exception : " + res.Buf)
	}

	if res.Cmd != resCmd || (resArg != XS_CMD_NONE && res.GetArg() != resArg) {
		return nil, errors.New("Unexpected respond {CMD:" + strconv.Itoa(res.Cmd) + ", ARG:" + strconv.Itoa(res.GetArg()) + "}")
	}

	return res, nil

}

func (x *XSServer) SendCommand(cmd *XSCommand) error {
	bytes, _ := cmd.ToBytes()
	return x.write(bytes)
}

func (x *XSServer) GetRespond() *XSCommand {
	data, err := x.read(8)
	if err != nil {
		//fmt.Println("get respond error")
		//fmt.Println(err)
		return NewCommand1(XS_CMD_NONE)
	}
	//fmt.Println("call GetRespond function")
	//fmt.Println(data)

	//hdr := make(map[string]interface{})
	cmd := int(data[0])
	arg1 := int(data[1])
	arg2 := int(data[2])

	buf1Len := int(data[3])
	bufLenBytes := data[4:8]
	bufLen := int(ToUInt32(bufLenBytes, 60))

	b1, _ := x.read(buf1Len)
	buf1 := string(b1)

	b, _ := x.read(bufLen)
	buf := string(b)

	return NewCommand(cmd, arg1, arg2, buf, buf1)
}

func (x *XSServer) HasRespond() bool {
	if x.conn == nil || (x.flag&(BROKEN|FILE)) == 1 {
		return false
	}

	return true
}

func (x *XSServer) appendSendBuffer(bytes []byte) {
	for _, b := range bytes {
		x.sendBuffer = append(x.sendBuffer, b)
	}
}

func (x *XSServer) connect() error {
	addr := x.address
	seconds := 30
	timeout := time.Duration(seconds) * time.Second
	sock, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return err
	}
	x.conn = sock
	return nil
}

func (x *XSServer) read(dataLen int) ([]byte, error) {
	if dataLen == 0 {
		return nil, errors.New("data len is empty")
	}

	err := x.check()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, dataLen)
	_, err = x.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (x *XSServer) write(data []byte) error {

	//x.conn.SetWriteDeadline(time.Now().Add(time.Duration(30) * time.Second))

	dataLen := len(data)
	if dataLen == 0 {
		return nil
	}

	err := x.check()
	if err != nil {
		return err
	}
	n, err := x.conn.Write(data)
	if err != nil {
		return err
	}

	fmt.Println("write data length : " + strconv.Itoa(n))

	// var n, idx int
	// var buf []byte
	// var bufLen int
	// bufLen = 1024
	// n, idx = 0, 0
	// for {
	// 	buf = data[idx:(bufLen + idx)]
	// 	n, err = x.conn.Write(buf)
	// 	if n == 0 || n == dataLen {
	// 		break //发送完成
	// 	}
	// 	dataLen -= n
	// 	idx += n
	// }

	return nil
}

func (x *XSServer) check() error {
	if x.conn == nil {
		return errors.New("No server connection")
	}
	if (x.flag & BROKEN) == 1 {
		return errors.New("Broken server connection")
	}
	return nil
}
