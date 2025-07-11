package utils

import (
	"reflect"
	"unicode"
)

// Interface untuk Tabler (mirip GORM)
type Tabler interface {
	TableName() string
}

// ToSnakeCase converts PascalCase to snake_case.
func ToSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

func GetTableNameRuntime(model any) string {
	if t, ok := model.(Tabler); ok {
		return t.TableName()
	}
	typ := reflect.TypeOf(model)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return ToSnakeCase(typ.Name()) + "s"
}

// GetTableName will return the table name from struct or fallback from struct name.
func GetTableName(model any) string {
	if tabler, ok := model.(Tabler); ok {
		return tabler.TableName()
	}
	typ := reflect.TypeOf(model)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return ToSnakeCase(typ.Name()) + "s"
}
