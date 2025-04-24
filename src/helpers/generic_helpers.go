package helpers

import (
	"database/sql"
	"errors"
	"fmt"
	"gin/src/configs/database"
	"log"
	"reflect"
	"regexp"
	"strings"
)

type Tabler interface {
	TableName() string
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func InsertModel[T any](model *T) error {
	if database.GormDB != nil {
		return database.GormDB.Create(model).Error
	}

	if database.SQLDB == nil {
		fmt.Println("❌ No database connection available: %w", sql.ErrConnDone)
		return sql.ErrConnDone
	}

	val := reflect.ValueOf(model).Elem()
	typ := val.Type()

	var tableName string
	if t, ok := any(model).(Tabler); ok {
		tableName = t.TableName()
	} else {
		tableName = ToSnakeCase(typ.Name()) + "s"
	}
	var columns []string
	var placeholders []string
	var values []any
	var primaryKeyField reflect.Value
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Cari primary key (id) berdasarkan tag gorm atau db
		gormTag := field.Tag.Get("gorm")
		dbTag := field.Tag.Get("db")

		if strings.Contains(gormTag, "primaryKey") || dbTag == "id" {
			primaryKeyField = fieldValue
			continue // kita skip ID, karena akan diisi oleh RETURNING
		}

		// Field yang tidak valid untuk insert
		if !fieldValue.CanInterface() || dbTag == "-" || dbTag == "" {
			continue
		}

		columns = append(columns, dbTag)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(columns)))
		values = append(values, fieldValue.Interface())
	}

	if len(columns) == 0 {
		return errors.New("no columns to insert")
	}

	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s) RETURNING id`,
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	log.Printf("🛠 SQL: %s | Values: %#v", query, values)

	if !primaryKeyField.CanAddr() {
		return errors.New("cannot get address of primary key field")
	}

	return database.SQLDB.QueryRow(query, values...).Scan(primaryKeyField.Addr().Interface())
}

func GetAllModels[T any](models *[]T, limit, offset int, orderBy string) error {
	if database.GormDB != nil {
		query := database.GormDB.Limit(limit).Offset(offset)
		if orderBy != "" {
			query = query.Order(orderBy)
		}
		return query.Find(models).Error
	}

	if database.SQLDB == nil {
		return sql.ErrConnDone
	}

	// SQL native
	var model T
	table := GetTableName(&model)

	query := fmt.Sprintf("SELECT * FROM %s", table)

	if orderBy != "" {
		query += fmt.Sprintf(" ORDER BY %s", orderBy)
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := database.SQLDB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var item T
		dest, err := scanRowDestinations(&item)
		if err != nil {
			return err
		}
		if err := rows.Scan(dest...); err != nil {
			return err
		}
		*models = append(*models, item)
	}

	return nil
}

func scanRowDestinations[T any](model *T) ([]any, error) {
	val := reflect.ValueOf(model)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return nil, errors.New("model must be a non-nil pointer")
	}

	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return nil, errors.New("model must point to a struct")
	}

	var dest []any
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		if field.CanSet() {
			dest = append(dest, field.Addr().Interface())
		}
	}
	return dest, nil
}

func GetTableName[T any](model *T) string {
	if t, ok := any(model).(Tabler); ok {
		return t.TableName()
	}
	typ := reflect.TypeOf(model).Elem()
	return ToSnakeCase(typ.Name()) + "s" // fallback User -> users
}

func GetModelByID[T any](model *T, id any) error {
	if database.GormDB != nil {
		return database.GormDB.First(model, id).Error
	}

	if database.SQLDB == nil {
		fmt.Println("❌ No database connection available: %w", sql.ErrConnDone)
		return sql.ErrConnDone
	}

	table := GetTableName(model)
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1 LIMIT 1", table)
	row := database.SQLDB.QueryRow(query, id)

	return scanRowIntoStruct(row, model)
}

func scanRowIntoStruct[T any](row *sql.Row, model *T) error {
	val := reflect.ValueOf(model)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("model must be a non-nil pointer")
	}

	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return errors.New("model must point to a struct")
	}

	var dest []any
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		if field.CanSet() {
			dest = append(dest, field.Addr().Interface())
		}
	}

	// return row.Scan(dest...)
	err := row.Scan(dest...)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("record not found")
	}
	if err != nil {
		return fmt.Errorf("scan error: %w", err)
	}

	return nil
}

func UpdateModelByID[T any](model *T, id any) error {
	if database.GormDB != nil {
		return database.GormDB.Model(model).Where("id = ?", id).Updates(model).Error
	}

	if database.SQLDB == nil {
		return sql.ErrConnDone
	}

	val := reflect.ValueOf(model).Elem()
	typ := val.Type()
	table := GetTableName(model)

	var sets []string
	var values []any

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		if !value.CanInterface() || field.Name == "ID" {
			continue
		}

		tag := field.Tag.Get("db")
		if tag == "" {
			tag = strings.ToLower(field.Name)
		}

		sets = append(sets, fmt.Sprintf("%s = $%d", tag, len(values)+1))
		values = append(values, value.Interface())
	}

	if len(sets) == 0 {
		return errors.New("no fields to update")
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", table,
		strings.Join(sets, ", "), len(values)+1)
	values = append(values, id)

	_, err := database.SQLDB.Exec(query, values...)
	return err
}

func DeleteModelByID[T any](model *T, id any) error {
	if database.GormDB != nil {
		return database.GormDB.Delete(model, id).Error
	}

	if database.SQLDB == nil {
		return sql.ErrConnDone
	}

	table := GetTableName(model)
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", table)

	_, err := database.SQLDB.Exec(query, id)
	return err
}

func FindOneByField[T any](model *T, field string, value any) error {
	if database.GormDB != nil {
		return database.GormDB.Where(fmt.Sprintf("%s = ?", field), value).First(model).Error
	}

	if database.SQLDB == nil {
		return sql.ErrConnDone
	}

	table := GetTableName(model)
	fmt.Println("Table name:", table)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 LIMIT 1", table, field)
	row := database.SQLDB.QueryRow(query, value)

	return scanRowIntoStruct(row, model)
}
