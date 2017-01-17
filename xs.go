// xs project xs.go
package xs

import (
	"errors"
)

type XS struct {
	search             *XSSearch
	index              *XSIndex
	tokenizer          XSTokenizer
	scheme, bindScheme *XSFieldScheme
	config             map[string]interface{}
}

func NewXS(iniFile string) (*XS, error) {
	xs := &XS{}
	err := xs.LoadIniFile(iniFile)
	if err != nil {
		return nil, err
	}
	return xs, nil
}

func (x *XS) LoadIniFile(iniFile string) error {
	c, err := LoadInIFile(iniFile)
	if err != nil {
		return err
	}
	x.config = c
	scheme := NewFieldScheme()
	for k, v := range x.config {
		vv, ok := v.(map[string]string)
		if ok {
			scheme.AddField(k, vv)
		}
	}
	ok := scheme.CheckValid()
	if !ok {
		return errors.New("check scheme invalid")
	}
	x.scheme, x.bindScheme = scheme, scheme
	return nil
}

func (x *XS) GetScheme() *XSFieldScheme {
	return x.scheme
}

func (x *XS) GetConfig() map[string]interface{} {
	return x.config
}

func (x *XS) GetConfigStringValue(name string) string {
	v, ok := x.config[name]
	if !ok {
		return ""
	}
	vv, ok := v.(string)
	if !ok {
		return ""
	}
	return vv
}

func (x *XS) GetConfigMapValue(name string) map[string]string {
	v, ok := x.config[name]
	if !ok {
		return nil
	}
	vv, ok := v.(map[string]string)
	if !ok {
		return nil
	}
	return vv
}

func (x *XS) GetName() string {
	return x.GetConfigStringValue("project.name")
}

func (x *XS) SetName(value string) {
	x.config["project.name"] = value
}

func (x *XS) GetSearch() (*XSSearch, error) {
	if x.search == nil {
		addr := x.GetConfigStringValue("server.search")
		app := x.GetName()
		server, err := NewServer(addr, nil)
		if err != nil {
			return nil, err
		}
		server.SetProject(app, "")
		x.search = NewSearch(server)
		x.search.xs = x
	}
	return x.search, nil
}

func (x *XS) GetField(name string) *XSFieldMeta {
	return x.scheme.GetField(name)
}

func (x *XS) Close() {
	x.search.Close()
}
