package xs

import (
	"errors"
)

const (
	MIXED_VNO = 255
)

type XSFieldScheme struct {
	fields  map[string]*XSFieldMeta
	typeMap map[int]string
	vnoMap  map[int]string
}

func NewFieldScheme() *XSFieldScheme {
	return &XSFieldScheme{
		fields:  make(map[string]*XSFieldMeta),
		typeMap: make(map[int]string),
		vnoMap:  make(map[int]string),
	}
}

func (x *XSFieldScheme) GetAllFields() map[string]*XSFieldMeta {
	return x.fields
}

func (x *XSFieldScheme) GetVnoMap() map[int]string {
	return x.vnoMap
}

func (x *XSFieldScheme) GetFileByVno(vno int) *XSFieldMeta {
	name, ok := x.vnoMap[vno]
	if ok {
		return x.GetField(name)
	} else {
		return nil
	}
}

func (x *XSFieldScheme) GetField(name string) *XSFieldMeta {
	f, ok := x.fields[name]
	if ok {
		return f
	}
	return nil
}

func (x *XSFieldScheme) GetBodyField() *XSFieldMeta {
	name, ok := x.typeMap[TYPE_BODY]
	if ok {
		return x.GetField(name)
	}
	return nil
}

func (x *XSFieldScheme) GetTitleField() *XSFieldMeta {
	name, ok := x.typeMap[TYPE_TITLE]
	if ok {
		return x.GetField(name)
	}
	return nil
}

func (x *XSFieldScheme) GetIDField() *XSFieldMeta {
	name, ok := x.typeMap[TYPE_ID]
	if ok {
		return x.GetField(name)
	}
	return nil
}

func (x *XSFieldScheme) AddField(name string, config map[string]string) error {
	field := NewFieldMeta(name, config)
	_, ok := x.fields[field.Name]
	if ok {
		return errors.New("Duplicated field name: " + field.Name)
	}
	if field.IsSpeical() {
		_, ok := x.typeMap[field.Type]
		if ok {
			return errors.New("Duplicated type:  " + field.Name)
		}
		x.typeMap[field.Type] = field.Name
	}
	if field.Type == TYPE_BODY {
		field.Vno = MIXED_VNO
	} else {
		field.Vno = len(x.vnoMap)
	}
	x.vnoMap[field.Vno] = field.Name
	x.fields[field.Name] = field
	return nil
}

func (x *XSFieldScheme) CheckValid() bool {
	_, ok := x.typeMap[TYPE_ID]
	if !ok {
		return false
	}
	return true
}
