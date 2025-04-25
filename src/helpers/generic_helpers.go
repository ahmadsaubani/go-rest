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
	"time"
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

	now := time.Now()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		gormTag := field.Tag.Get("gorm")
		dbTag := field.Tag.Get("db")

		if strings.Contains(gormTag, "primaryKey") || dbTag == "id" {
			primaryKeyField = fieldValue
			continue
		}

		// Lewati field tidak valid untuk insert
		if !fieldValue.CanInterface() || dbTag == "-" || dbTag == "" {
			continue
		}

		// Otomatis isi created_at dan updated_at
		if strings.ToLower(dbTag) == "created_at" || strings.ToLower(field.Name) == "CreatedAt" {
			if fieldValue.CanSet() && fieldValue.Type() == reflect.TypeOf(time.Time{}) {
				fieldValue.Set(reflect.ValueOf(now))
			}
		}
		if strings.ToLower(dbTag) == "updated_at" || strings.ToLower(field.Name) == "UpdatedAt" {
			if fieldValue.CanSet() && fieldValue.Type() == reflect.TypeOf(time.Time{}) {
				fieldValue.Set(reflect.ValueOf(now))
			}
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
		structField := elem.Type().Field(i)

		dbTag := structField.Tag.Get("db")
		if dbTag == "-" || dbTag == "" {
			continue // skip fields not mapped to DB
		}

		if field.CanSet() {
			dest = append(dest, field.Addr().Interface())
		}
	}

	err := row.Scan(dest...)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("record not found for %T", model)
	}
	if err != nil {
		return fmt.Errorf("scan error: %w", err)
	}

	return nil
}

func UpdateModelByIDWithMap[T any](updatedFields map[string]interface{}, id any) error {
	if database.GormDB != nil {
		// Menggunakan new(T) untuk memberikan tipe eksplisit ke GORM
		// Dengan new(T), kita bisa memastikan bahwa tipe tersebut sesuai
		return database.GormDB.Model(new(T)).Where("id = ?", id).Updates(updatedFields).Error
	}

	if database.SQLDB == nil {
		return sql.ErrConnDone
	}

	// Menggunakan refleksi untuk mendapatkan nama tabel dengan tipe eksplisit
	table := GetTableName(new(T)) // new(T) memberikan tipe eksplisit

	// Tambahkan updated_at ke map jika belum ada
	if _, exists := updatedFields["updated_at"]; !exists {
		updatedFields["updated_at"] = time.Now()
	}

	var sets []string
	var values []any

	// Memproses map field untuk SQL update
	for column, value := range updatedFields {
		sets = append(sets, fmt.Sprintf("%s = $%d", column, len(values)+1))
		values = append(values, value)
	}

	// Membuat query untuk update
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", table, strings.Join(sets, ", "), len(values)+1)
	values = append(values, id)

	_, err := database.SQLDB.Exec(query, values...)
	return err
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

		// Inject updated_at = now
		if strings.ToLower(tag) == "updated_at" && value.Type() == reflect.TypeOf(time.Time{}) && value.CanSet() {
			value.Set(reflect.ValueOf(time.Now()))
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

func hasDeletedAt(model any) bool {
	val := reflect.ValueOf(model).Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if strings.ToLower(field.Name) == "deletedat" &&
			field.Type == reflect.TypeOf(time.Time{}) {
			return true
		}
	}
	return false
}

func DeleteModelByID[T any](model *T, id any) error {
	if database.GormDB != nil {
		// GORM punya soft delete bawaan, tapi kita handle manual biar konsisten
		if hasDeletedAt(model) {
			return database.GormDB.Model(model).
				Where("id = ?", id).
				Update("deleted_at", time.Now()).Error
		}
		return database.GormDB.Delete(model, id).Error
	}

	if database.SQLDB == nil {
		return sql.ErrConnDone
	}

	table := GetTableName(model)

	if hasDeletedAt(model) {
		query := fmt.Sprintf("UPDATE %s SET deleted_at = $1 WHERE id = $2", table)
		_, err := database.SQLDB.Exec(query, time.Now(), id)
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", table)
	_, err := database.SQLDB.Exec(query, id)
	return err
}

func FindOneByField[T any](model *T, conditions ...any) error {
	if len(conditions)%2 != 0 {
		return fmt.Errorf("conditions must be in key-value pairs")
	}

	if database.GormDB != nil {
		query := database.GormDB
		for i := 0; i < len(conditions); i += 2 {
			field := conditions[i].(string)
			value := conditions[i+1]
			query = query.Where(fmt.Sprintf("%s = ?", field), value)
		}
		return query.First(model).Error
	}

	if database.SQLDB == nil {
		return sql.ErrConnDone
	}

	table := GetTableName(model)
	whereClause := ""
	args := []any{}

	for i := 0; i < len(conditions); i += 2 {
		field := conditions[i].(string)
		value := conditions[i+1]
		if i > 0 {
			whereClause += " AND "
		}
		whereClause += fmt.Sprintf("%s = $%d", field, (i/2)+1)
		args = append(args, value)
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT 1", table, whereClause)

	row := database.SQLDB.QueryRow(query, args...)

	return scanRowIntoStruct(row, model)
}
