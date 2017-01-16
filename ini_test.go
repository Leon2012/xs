package xs

import (
	"fmt"
	"testing"
)

func TestLoadInI(t *testing.T) {
	iniFile := "/temp/demo.ini"
	ret, err := LoadInIFile(iniFile)
	if err != nil {
		t.Error(err)
	}
	fmt.Print(ret)
}

func TestSubstr(t *testing.T) {
	str := "abcdefghkl"
	pos := 5
	fmt.Println(str[0:pos])
	fmt.Println(str[pos:(len(str))])
}
