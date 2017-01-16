package xs

import (
	"reflect"
	"unicode"
)

func MapCopy(dst, src interface{}) {
	dv, sv := reflect.ValueOf(dst), reflect.ValueOf(src)
	for _, k := range sv.MapKeys() {
		dv.SetMapIndex(k, sv.MapIndex(k))
	}
}

func IsChineseString(str string) bool {
	ret := true
	for _, s := range str {
		if !unicode.Is(unicode.Scripts["Han"], s) {
			ret = false
			break
		}
	}
	return ret
}
