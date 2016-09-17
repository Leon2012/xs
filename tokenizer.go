package xs

import (
	"fmt"
	"strings"

	"github.com/huichen/sego"
)

const DFL = 0
const MULTI_MASK = 15

const (
	SCWS_MULTI_NONE    = 0
	SCWS_MULTI_SHORT   = 1
	SCWS_MULTI_DUALITY = 2
	SCWS_MULTI_ZMAIN   = 4
	SCWS_MULTI_ZALL    = 8
	SCWS_XDICT_XDB     = 1
	SCWS_XDICT_MEM     = 2
	SCWS_XDICT_TXT     = 4
)

type XSTokenizer interface {

	/**
	 * 执行分词并返回词列表
	 */
	GetTokens(value string, doc *XSDocument) []string
}

/**
 * 内置空分词器
 */
type XSTokenizerNone struct {
}

func (x *XSTokenizerNone) GetTokens(value string, doc *XSDocument) []string {
	return []string{}
}

/**
 * 内置整值分词器
 */
type XSTokenizerFull struct {
}

func (x *XSTokenizerFull) GetTokens(value string, doc *XSDocument) []string {
	return []string{value}
}

/**
 * 内置的分割分词器
 */
type XSTokenizerSplit struct {
	arg string
}

func NewXSTokenizerSplit(a string) *XSTokenizerSplit {
	return &XSTokenizerSplit{
		arg: a,
	}
}

func (x *XSTokenizerSplit) GetTokens(value string, doc *XSDocument) []string {
	return strings.Split(value, x.arg)
}

/**
 * 内置的定长分词器
 *
 * @author hightman <hightman@twomice.net>
 * @version 1.0.0
 * @package XS.tokenizer
 */
// type XSTokenizerXlen struct {
// }

/**
 * SCWS - 分词器(与搜索服务端通讯)
 */
type XSTokenizerScws struct {
	charset string
	server  *XSServer
	setting map[string]*XSCommand
	dict    []*XSCommand
}

type XSWord struct {
	Off  uint32
	Attr string
	Word string
}

func NewXSTokenizerScws(serv *XSServer, mode int) *XSTokenizerScws {
	ts := &XSTokenizerScws{
		charset: "UTF-8",
		server:  serv,
		setting: make(map[string]*XSCommand),
		dict:    []*XSCommand{},
	}
	ts.SetMulti(mode)
	return ts
}

func (x *XSTokenizerScws) GetTokens(value string, doc *XSDocument) []string {
	tokens := []string{}
	x.SetIgnore(true)
	words, err := x.GetResult(value)
	if err != nil {
		return nil
	}
	var word *XSWord
	for i := 0; i < len(words); i++ {
		word = words[i]
		tokens = append(tokens, word.Word)
	}

	return tokens
}

/**
 * 设置忽略标点符号
 */
func (x *XSTokenizerScws) SetIgnore(yes bool) {
	var arg2 int
	if yes {
		arg2 = 1
	} else {
		arg2 = 0
	}
	x.setting["ignore"] = NewCommand(XS_CMD_SEARCH_SCWS_SET, XS_CMD_SCWS_SET_IGNORE, arg2, "", "")
}

func (x *XSTokenizerScws) SetMulti(mode int) {
	v := mode & MULTI_MASK
	x.setting["multi"] = NewCommand(XS_CMD_SEARCH_SCWS_SET, XS_CMD_SCWS_SET_MULTI, v, "", "")
}

func (x *XSTokenizerScws) SetDict(fpath string, mode int) {
	x.setting["set_dict"] = NewCommand(XS_CMD_SEARCH_SCWS_SET, XS_CMD_SCWS_SET_DICT, mode, fpath, "")
	x.dict = []*XSCommand{}
}

func (x *XSTokenizerScws) AddDict(fpath string, mode int) {
	cmd := NewCommand(XS_CMD_SEARCH_SCWS_SET, XS_CMD_SCWS_ADD_DICT, mode, fpath, "")
	x.dict = append(x.dict, cmd)
}

/**
 * 设置散字二元组合
 */
func (x *XSTokenizerScws) SetDuality(yes bool) {
	var arg2 int
	if yes {
		arg2 = 1
	} else {
		arg2 = 0
	}
	x.setting["duality"] = NewCommand(XS_CMD_SEARCH_SCWS_SET, XS_CMD_SCWS_SET_DUALITY, arg2, "", "")
}

func (x *XSTokenizerScws) GetVersion() (string, error) {
	cmd := NewCommand(XS_CMD_SEARCH_SCWS_GET, XS_CMD_SCWS_GET_VERSION, 0, "", "")
	res, err := x.server.ExecCommand(cmd, XS_CMD_OK_INFO, XS_CMD_OK)
	if err != nil {
		return "", err
	}
	return res.Buf, nil
}

func (x *XSTokenizerScws) GetResult(text string) ([]*XSWord, error) {
	//fmt.Println("call GetResult function")
	words := []*XSWord{}
	//x.applySetting()
	cmd := NewCommand(XS_CMD_SEARCH_SCWS_GET, XS_CMD_SCWS_GET_RESULT, 0, text, "")
	//fmt.Println("send GetResult command")
	//fmt.Println(cmd.ToBytes())
	res, err := x.server.ExecCommand(cmd, XS_CMD_OK_SCWS_RESULT, XS_CMD_OK)
	if err != nil {
		return nil, err
	}

	for {
		if res.Buf == "" {
			break
		}
		bytes := []byte(res.Buf)
		//fmt.Println(bytes)
		bytesLen := len(bytes)
		if bytesLen >= 8 {
			off := ToUInt32(bytes[0:4], 60)
			attr := string(bytes[4:8])
			var word string
			wordLen := bytesLen - 8
			if wordLen > 0 {
				word = string(bytes[8:bytesLen])
			} else {
				word = ""
			}
			words = append(words, &XSWord{
				Off:  off,
				Attr: attr,
				Word: word,
			})
		}
		res = x.server.GetRespond()
	}
	return words, nil
}

func (x *XSTokenizerScws) GetTops(text string, limit int, xattr string) ([]*XSWord, error) {
	words := []*XSWord{}
	cmd := NewCommand(XS_CMD_SEARCH_SCWS_GET, XS_CMD_SCWS_GET_TOPS, limit, text, xattr)
	res, err := x.server.ExecCommand(cmd, XS_CMD_OK_SCWS_TOPS, XS_CMD_OK)
	if err != nil {
		return nil, err
	}
	for {
		if res.Buf == "" {
			break
		}
		bytes := []byte(res.Buf)
		bytesLen := len(bytes)
		if bytesLen >= 8 {
			off := ToUInt32(bytes[0:4], 60)
			attr := string(bytes[4:8])
			var word string
			wordLen := bytesLen - 8
			if wordLen > 0 {
				word = string(bytes[8:bytesLen])
			} else {
				word = ""
			}
			words = append(words, &XSWord{
				Off:  off,
				Attr: attr,
				Word: word,
			})
		}
		res = x.server.GetRespond()
	}
	return words, nil
}

func (x *XSTokenizerScws) HasWord(text, xattr string) bool {
	cmd := NewCommand(XS_CMD_SEARCH_SCWS_GET, XS_CMD_SCWS_HAS_WORD, 0, text, xattr)
	res, err := x.server.ExecCommand(cmd, XS_CMD_OK_SCWS_RESULT, XS_CMD_OK)
	if err != nil {
		return false
	}
	if res.Buf == "OK" {
		return true
	} else {
		return false
	}
}

func (x *XSTokenizerScws) applySetting() {
	//x.server.Reopen(false)
	for key, cmd := range x.setting {
		fmt.Println("cmd key:" + key)
		x.server.ExecCommand1(cmd)
	}
	cmds := x.dict
	//for key, cmds := range x.dict {
	cmdsLen := len(cmds)
	for i := 0; i < cmdsLen; i++ {
		fmt.Println("aaaaaa")
		cmd := cmds[i]
		x.server.ExecCommand1(cmd)
	}
	//}
}

/**
 * Sego分词器
 * https://github.com/huichen/sego
 */
type XSTokenizerSego struct {
	segmenter sego.Segmenter
	dictFile  string
}

func NewXSTokenizerSego(dictFile string) *XSTokenizerSego {
	ts := &XSTokenizerSego{
		dictFile: dictFile,
	}
	ts.segmenter.LoadDictionary(dictFile)
	return ts
}

func (x *XSTokenizerSego) GetTokens(value string, doc *XSDocument) []string {
	segments := x.segmenter.Segment([]byte(value))
	return sego.SegmentsToSlice(segments, false)
}
