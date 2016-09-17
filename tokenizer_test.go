package xs

import (
	"fmt"
	"strconv"
	"testing"
)

func TestGetResult(t *testing.T) {
	serv, err := NewServer("127.0.0.1:8384", nil)
	if err != nil {
		t.Error(err)
	}
	serv.SetProject("demo", "")

	scws := NewXSTokenizerScws(serv, 3)
	scws.SetIgnore(true)

	//text := "每一个 xunsearch 搜索项目都有一个独立的 INI 配置文件。DEMO 项目的配置文件 位于"
	text := "项目"

	words, err := scws.GetResult(text)

	if err != nil {
		fmt.Println("error:")
		fmt.Println(err)
	} else {
		var word *XSWord
		for i := 0; i < len(words); i++ {
			word = words[i]
			fmt.Println("offset:" + strconv.Itoa(int(word.Off)) + " attr:" + word.Attr + " word:" + word.Word)
		}
	}

	//fmt.Println(words)

	serv.Close(true)
	t.Log("closed")
}

func TestGetTops(t *testing.T) {
	serv, err := NewServer("192.168.88.134:8384", nil)
	if err != nil {
		t.Error(err)
	}
	serv.SetProject("demo", "")

	scws := NewXSTokenizerScws(serv, 3)
	scws.SetIgnore(true)

	text := "每一个 xunsearch 搜索项目都有一个独立的 INI 配置文件。DEMO 项目的配置文件 位于"
	//text := "test"

	words, err := scws.GetTops(text, 10, "")

	if err != nil {
		fmt.Println("error:")
		fmt.Println(err)
	} else {
		var word *XSWord
		for i := 0; i < len(words); i++ {
			word = words[i]
			fmt.Println("offset:" + strconv.Itoa(int(word.Off)) + " attr:" + word.Attr + " word:" + word.Word)
		}
	}

	//fmt.Println(words)

	serv.Close(true)
	t.Log("closed")
}

func TestHasWord(t *testing.T) {
	serv, err := NewServer("192.168.88.134:8384", nil)
	if err != nil {
		t.Error(err)
	}
	serv.SetProject("demo", "")

	scws := NewXSTokenizerScws(serv, 3)
	text := "pack"
	xattr := "test"

	b := scws.HasWord(text, xattr)

	fmt.Println(strconv.FormatBool(b))

	serv.Close(true)
	t.Log("closed")
}

func TestSego(t *testing.T) {
	dict := "/Users/pengleon/go/src/github.com/huichen/sego/data/dictionary.txt"
	text := "每一个 xunsearch 搜索项目都有一个独立的 INI 配置文件。DEMO 项目的配置文件 位于"
	sego := NewXSTokenizerSego(dict)
	results := sego.GetTokens(text, nil)
	t.Log(results)
}
