package xs

import (
	"testing"
)

func TestConvert(t *testing.T) {
	var i int
	i = 200

	t.Log(uint8(i))

	//t.Log(uint8(280))

	t.Log(byte(254))
}

func TestNewCommand(t *testing.T) {
	cmd := XS_CMD_USE
	arg1 := 1
	arg2 := 2
	buf := "buf"
	buf1 := "buf1"

	command := NewCommand(cmd, arg1, arg2, buf, buf1)
	bytes, _ := command.ToBytes()
	t.Log(bytes)
}

func TestArg(t *testing.T) {
	cmd := XS_CMD_NONE
	arg1 := 0
	arg2 := 0
	buf := ""
	buf1 := ""

	command := NewCommand(cmd, arg1, arg2, buf, buf1)
	t.Log(command.GetArg())

	command.SetArg(10)
	t.Log(command)

	// var i int
	// i = arg1 | (arg2 << 8)

	// t.Log(i)
}

func TestSwcsCommand(t *testing.T) {
	text := "中华人民共和国"
	cmd := NewCommand(XS_CMD_SEARCH_SCWS_GET, XS_CMD_SCWS_GET_RESULT, 0, text, "")
	t.Log(cmd.ToBytes())
}
