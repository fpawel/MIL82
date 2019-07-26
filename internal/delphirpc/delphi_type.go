package delphirpc

import (
	"fmt"
	r "reflect"
)

type delphiType struct {
	name    string
	members map[string]delphiType
	elem    delphiType
	kind    delphiTypeKind
}

type delphiTypeKind int

const (
	delphiPOD delphiTypeKind = iota
	delphiRecord
	delphiArray
)

func newDelphiType(t r.Type) delphiType {

	pod := func(name string) delphiType {
		return delphiType{name: name, kind: delphiPOD}
	}

	switch t.Kind() {

	case r.Float32:
		return pod("Single")

	case r.Float64:
		return pod("Double")

	case r.Int:
		return pod("Integer")

	case r.Uint8:
		return pod("Byte")

	case r.Uint16:
		return pod("Word")

	case r.Uint32:
		return pod("Cardinal")

	case r.Uint64:
		return pod("UInt64")

	case r.Int8:
		return pod("ShortInt")

	case r.Int16:
		return pod("SmallInt")

	case r.Int32:
		return pod("Integer")

	case r.Int64:
		return pod("Int64")

	case r.Bool:
		return pod("Boolean")

	case r.String:
		return pod("string")

	case r.Array, r.Slice:
		return delphiType{
			name: t.Name(),
			kind: delphiArray,
			elem: newDelphiType(t.Elem()),
		}

	case r.Struct:
		x := delphiType{
			name:    t.Name(),
			kind:    delphiRecord,
			members: make(map[string]delphiType),
		}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			x.members[f.Name] = newDelphiType(f.Type)
		}
		return x

	default:
		panic(fmt.Sprintf("bad type: %s, %v", t.Name(), t.Kind()))
	}
}
