package xs

import (
	"fmt"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	m := make(map[string]string)
	m["name"] = "leon"

	val, ok := m["name1"]
	fmt.Println("val:" + val + " ok:" + strconv.FormatBool(ok))
}

func TestMap2(t *testing.T) {
	m := make(map[string]map[string]int)
	v := make(map[string]int)
	v["name"] = 10

	m["field"] = v

	fmt.Println(m)

}
