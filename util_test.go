package xs

import (
	"fmt"
	"testing"
)

func TestMapCopy(t *testing.T) {
	m1 := map[int]string{1: "one", 2: "two"}
	m2 := map[int]string{3: "three"}
	MapCopy(m2, m1)
	fmt.Println(m2)
}
