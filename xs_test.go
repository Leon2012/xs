package xs

import (
	"fmt"
	"strconv"
	"testing"
)

const APP_INI_FILE = "/temp/demo.ini"

func TestXSSearch(t *testing.T) {
	query := "项目"
	offset := 0
	limit := 1
	xs, err := NewXS(APP_INI_FILE)
	if err != nil {
		t.Error(err)
	}
	search, err := xs.GetSearch()
	if err != nil {
		t.Error(err)
	}
	docs, err := search.Search(query, offset, limit)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(docs)
	fmt.Println("total count : " + strconv.Itoa(search.count))
	for _, doc := range docs {
		v := doc.F("subject")
		subject, ok := v.(string)
		if ok {
			fmt.Println("subject:" + subject)
		}
	}
}
