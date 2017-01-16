package xs

import (
	"encoding/binary"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestStringSplit(t *testing.T) {
	query := "项目   a b " + "\n" + "祖国" + "\t"
	fmt.Println(query)
	rgp := regexp.MustCompile("[ \\t\\r\\n]+")
	queries := rgp.Split(query, -1)
	fmt.Println(queries)
	newQuery := ""
	for i := 0; i < len(queries); i++ {
		part := strings.TrimSpace(queries[i])
		//fmt.Println(queries[i])
		if part == "" {
			continue
		}
		if newQuery != "" {
			newQuery += " "
		}
		pos := strings.Index(part, ":")
		if pos != -1 {
			// for n := 0; n < pos; n++ {
			// 	if (strings.Index(s, sep))
			// }
			//name := part[0:pos]

			if len(part) > 1 && (part[0:1] == "+" || part[0:1] == "-") && (part[1:2] != "(" && IsChineseString(part)) {
				newQuery += (part[0:1] + "(" + part[1:] + ")")
				continue
			}
		}

		newQuery += part
	}
}

func TestSearch(t *testing.T) {
	app := "demo"
	query := "项目"
	offset := 0
	pageSize := 2

	serv, err := NewServer(XS_SEARCH_HOST, nil)
	if err != nil {
		t.Error(err)
	}
	serv.SetProject(app, "")

	pageBytes := make([]byte, 8)
	offsetBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(offsetBytes, uint32(offset))
	pageSizeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(pageSizeBytes, uint32(pageSize))

	copy(pageBytes[0:4], offsetBytes)
	copy(pageBytes[4:8], pageSizeBytes)

	cmd := NewCommand(XS_CMD_SEARCH_GET_RESULT, 0, XS_CMD_QUERY_OP_AND, query, string(pageBytes))
	res, err := serv.ExecCommand(cmd, XS_CMD_OK_RESULT_BEGIN, XS_CMD_OK)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("res:" + res.String())
	data := []byte(res.Buf)
	if len(data) == 4 {
		count := binary.LittleEndian.Uint32(data)
		fmt.Println("count:" + strconv.Itoa(int(count)))
	}

	//var docs []*XSDocument
	//docs = []*XSDocument{}

	for {
		res = serv.GetRespond()
		if res.Cmd == XS_CMD_OK && res.GetArg() == XS_CMD_OK_RESULT_END {
			break
		} else if res.Cmd == XS_CMD_SEARCH_RESULT_DOC {
			//doc := NewDocument("UTF-8")
			//docs = append(docs, doc)
			fmt.Println("doc:")

			printMeta([]byte(res.Buf))
		} else if res.Cmd == XS_CMD_SEARCH_RESULT_FIELD {

		} else {
			err = errors.New("Unexpected respond in search {CMD:" + strconv.Itoa(res.Cmd) + ", ARG:" + strconv.Itoa(res.GetArg()) + "}")
		}
		fmt.Println("res:" + res.String())
	}

	serv.Close(false)
}

func printMeta(data []byte) {
	fmt.Println(data)
	fmt.Println("data length:" + strconv.Itoa(len(data)))

	metas := make(map[string]interface{})

	docIdBytes := make([]byte, 4)
	copy(docIdBytes, data[0:4])
	fmt.Println(docIdBytes)
	docId := binary.LittleEndian.Uint32(docIdBytes)
	metas["docid"] = docId

	rankBytes := make([]byte, 4)
	copy(rankBytes, data[4:8])
	rank := binary.LittleEndian.Uint32(rankBytes)
	metas["rank"] = rank

	ccountBytes := make([]byte, 4)
	copy(ccountBytes, data[8:12])
	ccount := binary.LittleEndian.Uint32(ccountBytes)
	metas["ccount"] = ccount

	percentBytes := make([]byte, 4)
	copy(percentBytes, data[12:16])
	percent := int32(binary.LittleEndian.Uint32(percentBytes))
	metas["percent"] = percent

	weightBytes := make([]byte, 4)
	copy(weightBytes, data[16:20])
	weight := ByteToFloat32(weightBytes)
	metas["weight"] = weight

	fmt.Println(metas)
}
