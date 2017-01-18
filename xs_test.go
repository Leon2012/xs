package xs

import (
	"fmt"
	"strconv"
	"testing"
)

const APP_INI_FILE = "/temp/demo.ini"

func TestXSSearch(t *testing.T) {
	query := "习近平"
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

func TestXSIndex(t *testing.T) {
	xs, err := NewXS(APP_INI_FILE)
	if err != nil {
		t.Error(err)
	}
	index, err := xs.GetIndex()
	if err != nil {
		t.Error(err)
	}

	doc := NewDocument("UTF-8")
	doc.SetField("pid", 4, false)
	doc.SetField("subject", "习近平在世界经济论坛年会开幕式上的演讲", false)
	doc.SetField("message", "我想说的是，困扰世界的很多问题，并不是经济全球化造成的。比如，过去几年来，源自中东、北非的难民潮牵动全球，数以百万计的民众颠沛流离，甚至不少年幼的孩子在路途中葬身大海，让我们痛心疾首。导致这一问题的原因，是战乱、冲突、地区动荡。解决这一问题的出路，是谋求和平、推动和解、恢复稳定。再比如，国际金融危机也不是经济全球化发展的必然产物，而是金融资本过度逐利、金融监管严重缺失的结果。把困扰世界的问题简单归咎于经济全球化，既不符合事实，也无助于问题解决。", false)
	doc.SetField("chrono", 1484727451, false)

	index.Update(doc, true)

}
