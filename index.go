package xs

import (
	"errors"
)

type XSIndex struct {
	xs      *XS
	server  *XSServer
	buf     string
	bufSize int
	rebuild bool
	servers []*XSServer
}

func NewIndex(s *XSServer, x *XS) *XSIndex {
	return &XSIndex{
		buf:     "",
		bufSize: 0,
		rebuild: false,
		servers: []*XSServer{},
		server:  s,
		xs:      x,
	}
}

func (x *XSIndex) AddServer(addr string) (*XSServer, error) {
	serv, err := NewServer(addr, x.xs)
	if err != nil {
		return nil, err
	}
	x.servers = append(x.servers, serv)
	return serv, nil
}

func (x *XSIndex) execCommand(cmd *XSCommand, resArg, resCmd int) (*XSCommand, error) {
	if resArg == 0 {
		resArg = XS_CMD_NONE
	}
	if resCmd == 0 {
		resCmd = XS_CMD_OK
	}
	res, err := x.server.ExecCommand(cmd, resArg, resCmd)
	if err != nil {
		return nil, err
	}
	for _, serv := range x.servers {
		serv.ExecCommand(cmd, resArg, resCmd)
	}
	return res, nil
}

func (x *XSIndex) Clean() {
	cmd := NewCommand(XS_CMD_INDEX_CLEAN_DB, 0, 0, "", "")
	x.execCommand(cmd, XS_CMD_OK_DB_CLEAN, 0)
}

func (x *XSIndex) Add(doc *XSDocument) {
	x.Update(doc, true)
}

/**
 * 更新索引文档
 * 该方法相当于先根据主键删除已存在的旧文档, 然后添加该文档
 * 如果你能明确认定是新文档, 则建议使用 {@link add}
 */
func (x *XSIndex) Update(doc *XSDocument, add bool) error {
	fid := x.xs.GetIdField()
	key := doc.F1(fid.String())
	if key == "" {
		return errors.New("Missing value of primary key")
	}
	cmd := NewCommand(XS_CMD_INDEX_REQUEST, XS_CMD_INDEX_REQUEST_ADD, 0, "", "")
	if add {
		cmd.Arg1 = XS_CMD_INDEX_REQUEST_UPDATE
		cmd.Arg2 = fid.Vno
		cmd.Buf = key
	}

	cmds := []*XSCommand{}
	fields := x.xs.GetAllFields()
	for _, field := range fields {
		value := doc.F1(field.String())
		if value != "" {
			var varg int
			if field.IsNumeric() {
				varg = XS_CMD_VALUE_FLAG_NUMERIC
			} else {
				varg = 0
			}
			value = field.Val(value)
			cmds = append(cmds, NewCommand(XS_CMD_DOC_VALUE, varg, field.Vno, value, ""))
		}
	}

	for _, command := range cmds {
		x.execCommand(command, 0, 0)
	}
	x.execCommand(NewCommand1(XS_CMD_INDEX_SUBMIT), XS_CMD_OK_RQST_FINISHED, 0)
	return nil
}

func (x *XSIndex) SetDb(name string) error {
	_, err := x.execCommand(NewCommand(XS_CMD_INDEX_SET_DB, 0, 0, name, ""), XS_CMD_OK_DB_CHANGED, 0)
	return err
}

func (x *XSIndex) Close() {
	x.server.Close(false)
	for _, serv := range x.servers {
		serv.Close(false)
	}
}
