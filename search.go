package xs

import (
	//"regexp"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const PAGE_SIZE = 10
const LOG_DB = "log_db"

type XSSearch struct {
	xs                                         *XS
	server                                     *XSServer
	defaultOp, limit, offset, count, lastCount int
	fieldSet, resetScheme                      bool
	chat, query                                string
	prefix                                     map[string]bool
}

func NewSearch(s *XSServer, x *XS) *XSSearch {
	search := &XSSearch{
		defaultOp:   XS_CMD_QUERY_OP_AND,
		limit:       0,
		offset:      0,
		count:       0,
		prefix:      make(map[string]bool),
		fieldSet:    false,
		resetScheme: false,
		lastCount:   0,
		chat:        "UTF-8",
	}
	search.server = s
	search.server.Xs = x
	return search
}

/**
 * 开启模糊搜索
 * 默认情况只返回包含所有搜索词的记录, 通过本方法可以获得更多搜索结果
 */
func (s *XSSearch) SetFuzzy(value bool) {
	if value == true {
		s.defaultOp = XS_CMD_QUERY_OP_OR
	} else {
		s.defaultOp = XS_CMD_QUERY_OP_AND
	}
}

func (s *XSSearch) Search(query string, offset, limit int) ([]*XSDocument, error) {
	query = s.preQueryString(query)
	if limit <= 0 {
		limit = PAGE_SIZE
	}
	pageBytes := make([]byte, 8)
	offsetBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(offsetBytes, uint32(offset))
	pageSizeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(pageSizeBytes, uint32(limit))
	copy(pageBytes[0:4], offsetBytes)
	copy(pageBytes[4:8], pageSizeBytes)

	cmd := NewCommand(XS_CMD_SEARCH_GET_RESULT, 0, XS_CMD_QUERY_OP_AND, query, string(pageBytes))
	//fmt.Println(cmd.String())
	res, err := s.server.ExecCommand(cmd, XS_CMD_OK_RESULT_BEGIN, XS_CMD_OK)
	if err != nil {
		return nil, err
	}
	if res != nil {
		data := []byte(res.Buf)
		if len(data) == 4 {
			count := int(binary.LittleEndian.Uint32(data))
			s.lastCount = count
			fmt.Println("count : " + strconv.Itoa(count))
		}
	}

	vnoes := s.xs.GetScheme().GetVnoMap()
	fmt.Println(vnoes)
	docs := []*XSDocument{}
	var doc *XSDocument
	for {
		res = s.server.GetRespond()
		if res.Cmd == XS_CMD_OK && res.GetArg() == XS_CMD_OK_RESULT_END {
			break
		} else if res.Cmd == XS_CMD_SEARCH_RESULT_DOC {
			doc = NewDocument("UTF-8")
			if res.Buf != "" {
				doc.AddMetas([]byte(res.Buf))
			}
			docs = append(docs, doc)
		} else if res.Cmd == XS_CMD_SEARCH_RESULT_FIELD {
			if doc != nil {
				name, ok := vnoes[res.GetArg()]
				if !ok {
					name = strconv.Itoa(res.GetArg())
				}
				doc.SetField(name, res.Buf, false)
			}
		} else if res.Cmd == XS_CMD_SEARCH_RESULT_MATCHED {
			if doc != nil {
				doc.SetField("matched", strings.Split(res.Buf, " "), true)
			}
		} else {
			err = errors.New("Unexpected respond in search {CMD:" + strconv.Itoa(res.Cmd) + ", ARG:" + strconv.Itoa(res.GetArg()) + "}")
		}
	}

	s.count = s.lastCount
	//serv.Close(false)
	return docs, nil
}

func (s *XSSearch) preQueryString(query string) string {
	query = strings.TrimSpace(query)
	if s.resetScheme {
		s.clearQuery()
	}
	// newQuery := ""
	// rgp := regexp.MustCompile("[ \\t\\r\\n]+")
	// queries := rgp.Split(query, -1)
	return query
}

/**
 * 清空默认搜索语句
 */
func (s *XSSearch) clearQuery() {
	cmd := NewCommand1(XS_CMD_QUERY_INIT)
	if s.resetScheme {
		cmd.Arg1 = 1
		s.resetScheme = false
		s.fieldSet = false
	}
	s.server.ExecCommand1(cmd)
	s.query = ""
	s.count = 0
}

/**
 * 登记搜索语句中的字段
 */
func (s *XSSearch) regQueryPrefix(name string) {
	_, ok := s.prefix[name]
	field := s.xs.GetField(name)
	if !ok && (field != nil) && field.Vno != MIXED_VNO {
		cmd := NewCommand(XS_CMD_QUERY_PREFIX, XS_CMD_PREFIX_NORMAL, field.Vno, name, "")
		s.server.ExecCommand1(cmd)
		s.prefix[name] = true
	}
}

/**
 * 设置默认搜索语句
 * 用于不带参数的 {@link count} 或 {@link search} 以及 {@link terms} 调用
 * 可与 {@link addWeight} 组合运用
 * @param string $query 搜索语句, 设为 null 则清空搜索语句, 最大长度为 80 字节
 * @return XSSearch 返回对象本身以支持串接操作
 */
func (s *XSSearch) SetQuery(query string) {
	s.clearQuery()
	if query != "" {
		s.query = query
		s.addQueryString(query, XS_CMD_QUERY_OP_AND, 1)
	}
}

/**
 * 增加默认搜索语句
 * @param string $query 搜索语句
 * @param int $addOp 与旧语句的结合操作符, 如果无旧语句或为空则这此无意义, 支持的操作符有:
 *        XS_CMD_QUERY_OP_AND
 *        XS_CMD_QUERY_OP_OR
 *        XS_CMD_QUERY_OP_AND_NOT
 *        XS_CMD_QUERY_OP_XOR
 *        XS_CMD_QUERY_OP_AND_MAYBE
 *        XS_CMD_QUERY_OP_FILTER
 * @param float $scale 权重计算缩放比例, 默认为 1表示不缩放, 其它值范围 0.xx ~ 655.35
 * @return string 修正后的搜索语句
 */
func (s *XSSearch) addQueryString(query string, addOp int, scale float32) string {
	query = s.preQueryString(query)
	var bScale string
	if scale > 0 && scale != 1 {
		v := uint16(scale * 100)
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, v)
		bScale = string(b)
	} else {
		bScale = ""
	}
	cmd := NewCommand(XS_CMD_QUERY_PARSE, addOp, s.defaultOp, query, bScale)
	s.server.ExecCommand1(cmd)
	return query
}

func (x *XSSearch) Close() {
	x.server.Close(false)
}
