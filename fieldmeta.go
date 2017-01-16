package xs

import (
	_ "strconv"
	"strings"
)

const (

	/**
	 * 词条权重最大值
	 */
	MAX_WDF = 0x3f

	/**
	 * 字段类型常量定义
	 */
	TYPE_STRING  = 0
	TYPE_NUMERIC = 1
	TYPE_DATE    = 2
	TYPE_ID      = 10
	TYPE_TITLE   = 11
	TYPE_BODY    = 12

	/**
	 * 索引标志常量定义
	 */
	FLAG_INDEX_SELF    = 0x01
	FLAG_INDEX_MIXED   = 0x02
	FLAG_INDEX_BOTH    = 0x03
	FLAG_WITH_POSITION = 0x10
	FLAG_NON_BOOL      = 0x80 // 强制让该字段参与权重计算 (非布尔)
)

type XSFieldMeta struct {
	Name       string
	Cutlen     int
	Weight     int
	Type       int
	Vno        int
	Tokenizer  string
	Flag       int
	Tokenizers []XSTokenizer
}

func NewFieldMeta(name string, config map[string]string) *XSFieldMeta {
	fm := &XSFieldMeta{}
	fm.Name = name
	if config != nil {
		fm.FromConfig(config)
	}
	return fm
}

func (x *XSFieldMeta) String() string {
	return x.Name
}

func (x *XSFieldMeta) Val(value string) string {
	return value
}

func (x *XSFieldMeta) IsNumeric() {
	return x.Type == TYPE_NUMERIC
}

func (x *XSFieldMeta) IsSpeical() {
	return x.Type == TYPE_ID || x.Type == TYPE_TITLE || x.Type == TYPE_BODY
}

func (x *XSFieldMeta) FromConfig(config map[string]string) {
	var predef string

	typeVal, ok := config["type"]
	if ok {
		predef = strings.ToUpper(typeVal)
		if predef == "ID" {
			x.Type = TYPE_ID
			x.Flag = FLAG_INDEX_SELF
			x.Tokenizer = "full"
		} else if predef == "TITLE" {
			x.Type = TYPE_TITLE
			x.Flag = FLAG_INDEX_BOTH | FLAG_WITH_POSITION
			x.Weight = 5
		} else if predef == "BODY" {
			x.Vno = MIXED_VNO
			x.Flag = FLAG_INDEX_SELF | FLAG_WITH_POSITION
			x.Cutlen = 300
		}
	}

	// indexVal, ok := config["index"]
	// if ok && x.Type != TYPE_BODY {
	// 	predef = strings.ToUpper(indexVal)
	// 	if predef == "SELF" {

	// 	}
	// }

	// cutLenVal, ok := config["cutlen"]
	// if ok {
	// 	x.Cutlen = strconv.Itoa(cutLenVal)
	// }

	// weightVal, ok := config["weight"]
	// if ok && x.Type != TYPE_BODY {
	// 	x.Weight = strconv.Itoa(weightVal) & MAX_WDF
	// }

	// phraseVal, ok := config["phrase"]
	// if (ok) {
	// 	if (phraseVal != "yes") {
	// 		x.Flag |= FLAG_WITH_POSITION
	// 	}else if (phraseVal != "no") {
	// 		x.Flag &= ~ FLAG_WITH_POSITION;
	// 	}
	// }
}
