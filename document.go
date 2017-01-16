package xs

import (
	_ "fmt"
	"strings"
)

type XSDocument struct {
	data    map[string]interface{}
	terms   map[string]map[string]int
	texts   map[string]string
	charset string
	meta    map[string]interface{}
	resSize int
}

func NewDocument(c string) *XSDocument {
	return &XSDocument{
		data:    make(map[string]interface{}),
		terms:   nil,
		texts:   nil,
		charset: c,
		meta:    nil,
		resSize: 20,
	}
}

func (x *XSDocument) SetMeta(data []byte) {
	//$_resFormat = 'Idocid/Irank/Iccount/ipercent/fweight';
	// metas := make(map[string]interface{})

	// docIdBytes := make([]byte, 4)
	// copy(data[0:4], docIdBytes)
	// docId := BytesToUint32(docIdBytes)
	// metas["docid"] = docId

}

func (x *XSDocument) Get(name string) interface{} {
	val, ok := x.data[name]
	if !ok {
		return nil
	}

	return val
}

func (x *XSDocument) Set(name string, value interface{}) {
	x.SetField(name, value, false)
}

func (x *XSDocument) SetField(name string, value interface{}, isMeta bool) {
	if value == nil {
		if isMeta {
			delete(x.meta, name)
		} else {
			delete(x.data, name)
		}
	} else {
		if isMeta {
			x.meta[name] = value
		} else {
			x.data[name] = value
		}
	}
}

func (x *XSDocument) GetFields() map[string]interface{} {
	return x.data
}

func (x *XSDocument) SetFields(m map[string]interface{}) {
	if m == nil {
		x.data = make(map[string]interface{})
		x.meta, x.terms, x.texts = nil, nil, nil
	} else {
		MapCopy(x.data, m)
	}
}

func (x *XSDocument) F(name string) interface{} {
	return x.Get(name)
}

func (x *XSDocument) GetAddTerms(field string) map[string]int {
	if x.terms == nil {
		return nil
	}

	val, ok := x.terms[field]
	if !ok {
		return nil
	}

	terms := make(map[string]int)
	for term, weight := range val {
		terms[term] = weight
	}
	return terms
}

func (x *XSDocument) AddTerm(field string, term string, weight int) {
	if x.terms == nil {
		x.terms = make(map[string]map[string]int)
	}

	val, ok := x.terms[field]
	if !ok {
		x.terms[field] = map[string]int{term: weight}
	} else {
		val1, ok := val[term]
		if !ok {
			x.terms[field] = map[string]int{term: weight}
		} else {
			x.terms[field] = map[string]int{term: (weight + val1)}
		}
	}
}

func (x *XSDocument) GetAddIndex(field string) string {
	if x.texts == nil {
		return ""
	}

	val, ok := x.texts[field]
	if !ok {
		return ""
	}
	return val
}

func (x *XSDocument) AddIndex(field, text string) {
	if x.texts == nil {
		x.texts = make(map[string]string)
	}
	val, ok := x.texts[field]
	if !ok {
		x.texts[field] = text
	} else {
		val += ("\n" + text)
		x.texts[field] = val
	}
}

func (x *XSDocument) GetCharset() string {
	return x.charset
}

func (x *XSDocument) SetCharset(c string) {
	x.charset = strings.ToUpper(c)
	if x.charset == "UTF8" {
		x.charset = "UTF-8"
	}
}
