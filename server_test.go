package xs

import (
	"fmt"
	"testing"
)

func TestBytes(t *testing.T) {
	data := []byte{0, 1, 2, 3}

	t.Log(data[0:3])

	t.Log(data[0:10])
}

func TestSendTimeout(t *testing.T) {
	serv, err := NewServer("192.168.88.134:8384", nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(serv)

	serv.SetTimeout(10)

	serv.Close(false)

	t.Log("closed")
}

func TestExecCommand(t *testing.T) {
	serv, err := NewServer("192.168.88.134:8383", nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(serv)

	cmd := NewCommand1(XS_CMD_HELLO)
	res, err := serv.ExecCommand(cmd, XS_CMD_NONE, XS_CMD_OK)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)

	serv.Close(false)

	t.Log("closed")
}
