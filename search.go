package xs

import (
//"regexp"
//"strings"
)

const PAGE_SIZE = 10
const LOG_DB = "log_db"

type XSSearch struct {
	server                                   *XSServer
	defaultOp, limit, offset, count          int
	prefix, fieldSet, resetScheme, lastCount bool
	chat, query                              string
}

func NewSearch(s *XSServer) *XSSearch {
	xs := &XSSearch{
		defaultOp:   XS_CMD_QUERY_OP_AND,
		limit:       0,
		offset:      0,
		count:       0,
		prefix:      false,
		fieldSet:    false,
		resetScheme: false,
		lastCount:   false,
		chat:        "UTF-8",
	}
	xs.server = s
	return xs
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

func (s *XSSearch) Search(query string) {

}

func (s *XSSearch) preQueryString(query string) {
	// query = strings.TrimSpace(query)
	// if s.resetScheme {
	// 	s.clearQuery()
	// }
	// newQuery := ""
	// rgp := regexp.MustCompile("[ \\t\\r\\n]+")
	// queries := rgp.Split(query, -1)

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
