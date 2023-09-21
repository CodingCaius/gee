// Schema对象通常用于ORM（对象-关系映射）框架中，它描述了数据表的结构，包括表名、字段名、字段类型等信息

package schema

import (
	"geeorm/dialect"
	"go/ast"
	"reflect"
)

// Field represents a column of database
type Field struct {
	Name string
	Type string
	Tag  string
}

// Schema represents a table of database
type Schema struct {
	Model      interface{} // 关联的 Go 结构体（模型）
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field // 一个映射，用于将字段名称映射到对应的 Field 结构体
}

// GetField 根据字段名称获取对应的 Field 结构体
func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// Parse 解析传入的目标模型，提取其中可导出的字段信息，创建Field对象，并将这些信息填充到Schema对象中，以便后续在ORM操作中使用
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	// 使用反射（reflection）获取目标模型的类型
	// 因为设计的入参是一个对象的指针，因此需要 reflect.Indirect() 获取指针指向的实例
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		// 检查字段是否是可导出的（即字段名是否以大写字母开头），同时排除匿名字段（非导出的字段）
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

// RecordValues 从一个给定的目标（dest）中提取字段的值，并返回一个包含这些值的切片
func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	// reflect.Indirect 用于获取目标对象的实际值，如果目标对象是指针，则获取指针指向的值。
	// 这是因为通常会传递一个指向结构体的指针作为 dest，所以需要获取结构体的实际值
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}