package helpers

import (
	"database/sql"
	"errors"
	"fmt"
	"gin/src/configs/database"
	"gin/src/utils/filters"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const maxBatchSize = 500

type Tabler interface {
	TableName() string
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// ToSnakeCase converts a given CamelCase string to snake_case.
// This function uses regular expressions to identify capital letters
// and insert underscores before them, then converts the entire string
// to lowercase. For example, "CamelCase" becomes "camel_case".

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// InsertModelBatch inserts a batch of models into the database. It will use GORM
// if the USE_GORM environment variable is set to "true", otherwise it will use
// native SQL. It will automatically set the created_at and updated_at fields to
// the current time if they are present in the model and are of type time.Time.
// The returned error is the error from the database operation.
//
// InsertModelBatch will automatically build a batch insert query from the given
// models. It will use the maximum batch size of 500 records for each batch.
// If the database connection is not available, InsertModelBatch returns
// sql.ErrConnDone.
func InsertModelBatch[T any](models []T) error {
	if len(models) == 0 {
		return nil
	}

	now := time.Now()
	useGorm := os.Getenv("USE_GORM") == "true"
	useSQL := !useGorm && database.SQLDB != nil

	if !useGorm && !useSQL {
		return fmt.Errorf("❌ No valid database connection available")
	}

	for start := 0; start < len(models); start += maxBatchSize {
		end := start + maxBatchSize
		if end > len(models) {
			end = len(models)
		}
		batch := models[start:end]

		// Transaksi GORM
		if useGorm {
			err := database.GormDB.Transaction(func(tx *gorm.DB) error {
				for i := range batch {
					val := reflect.ValueOf(&batch[i]).Elem()
					typ := val.Type()

					// Otomatis set created_at dan updated_at
					for j := 0; j < val.NumField(); j++ {
						field := typ.Field(j)
						fieldValue := val.Field(j)
						dbTag := field.Tag.Get("db")

						// Set created_at dan updated_at
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
					}
				}

				// Insert batch menggunakan GORM
				if err := tx.Create(&batch).Error; err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("❌ GORM transaction failed: %w", err)
			}
		}

		// Transaksi Native SQL
		if useSQL {
			tx, err := database.SQLDB.Begin()
			if err != nil {
				return fmt.Errorf("❌ SQL transaction begin failed: %w", err)
			}
			defer tx.Rollback() // Ensure rollback on failure

			firstVal := reflect.ValueOf(batch[0])
			typ := firstVal.Type()
			var tableName string
			if t, ok := any(batch[0]).(Tabler); ok {
				tableName = t.TableName()
			} else {
				tableName = ToSnakeCase(typ.Name()) + "s"
			}

			var columns []string
			for i := 0; i < typ.NumField(); i++ {
				field := typ.Field(i)
				dbTag := field.Tag.Get("db")
				gormTag := field.Tag.Get("gorm")

				if strings.Contains(gormTag, "primaryKey") || dbTag == "id" {
					continue
				}
				if dbTag != "" && dbTag != "-" {
					columns = append(columns, dbTag)
				}
			}

			if len(columns) == 0 {
				return errors.New("no columns to insert")
			}

			// Placeholder for batch values
			placeholderRows := []string{}
			allValues := []any{}
			paramIdx := 1

			for _, m := range batch {
				val := reflect.ValueOf(m)
				rowPlaceholders := []string{}

				// Set created_at dan updated_at
				for i := 0; i < val.NumField(); i++ {
					field := typ.Field(i)
					fieldValue := val.Field(i)
					dbTag := field.Tag.Get("db")
					gormTag := field.Tag.Get("gorm")

					if strings.Contains(gormTag, "primaryKey") || dbTag == "id" || dbTag == "-" || dbTag == "" {
						continue
					}

					// Set created_at dan updated_at
					if strings.ToLower(dbTag) == "created_at" || strings.ToLower(field.Name) == "CreatedAt" {
						allValues = append(allValues, now)
					} else if strings.ToLower(dbTag) == "updated_at" || strings.ToLower(field.Name) == "UpdatedAt" {
						allValues = append(allValues, now)
					} else {
						allValues = append(allValues, fieldValue.Interface())
					}
					rowPlaceholders = append(rowPlaceholders, fmt.Sprintf("$%d", paramIdx))
					paramIdx++
				}
				placeholderRows = append(placeholderRows, "("+strings.Join(rowPlaceholders, ", ")+")")
			}

			query := fmt.Sprintf(
				"INSERT INTO %s (%s) VALUES %s",
				tableName,
				strings.Join(columns, ", "),
				strings.Join(placeholderRows, ", "),
			)

			// Execute SQL batch insert
			if _, err := tx.Exec(query, allValues...); err != nil {
				return fmt.Errorf("❌ SQL batch insert failed: %w", err)
			}

			// Commit the transaction
			if err := tx.Commit(); err != nil {
				return fmt.Errorf("❌ SQL transaction commit failed: %w", err)
			}
		}
	}

	return nil
}

// InsertModel inserts a model into the database.
//
// If GORM is enabled, it will use GORM's Create method.
// If GORM is disabled, it will use the SQLDB to execute an INSERT query.
//
// The returned error is the error from the database operation.
//
// InsertModel will automatically set the created_at and updated_at fields to the current time
// if they are present in the model and are of type time.Time.
//
// InsertModel will return an error if the model has no valid columns to insert.
//
// InsertModel will return an error if the primary key field is not addressable.
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

	if !primaryKeyField.CanAddr() {
		return errors.New("cannot get address of primary key field")
	}

	return database.SQLDB.QueryRow(query, values...).Scan(primaryKeyField.Addr().Interface())
}

// func GetAllModels[T any](models *[]T, limit, offset int, orderBy string) error {
// 	if database.GormDB != nil {
// 		query := database.GormDB.Limit(limit).Offset(offset)
// 		if orderBy != "" {
// 			query = query.Order(orderBy)
// 		}
// 		return query.Find(models).Error
// 	}

// 	if database.SQLDB == nil {
// 		return sql.ErrConnDone
// 	}

// 	// SQL native
// 	var model T
// 	table := GetTableName(&model)

// 	query := fmt.Sprintf("SELECT * FROM %s", table)

// 	if orderBy != "" {
// 		query += fmt.Sprintf(" ORDER BY %s", orderBy)
// 	}

// 	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

// 	rows, err := database.SQLDB.Query(query)
// 	if err != nil {
// 		return err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var item T
// 		dest, err := scanRowDestinations(&item)
// 		if err != nil {
// 			return err
// 		}
// 		if err := rows.Scan(dest...); err != nil {
// 			return err
// 		}
// 		*models = append(*models, item)
// 	}

// 	return nil
// }

// GetAllModels will fetch all records from the database based on the given limit,
// offset, and orderBy. It will use GORM if the USE_GORM environment variable is set
// to "true", otherwise it will use native SQL. It will automatically build a WHERE
// clause from the query string parameters of the given gin.Context.
func GetAllModels[T any](ctx *gin.Context, models *[]T, limit, offset int, orderBy string) error {
	useGORM := os.Getenv("USE_GORM") == "true"
	if useGORM {
		useGORM = true
	} else {
		useGORM = false
	}

	// GORM
	if useGORM && database.GormDB != nil {
		query := database.GormDB.Limit(limit).Offset(offset)

		if orderBy != "" {
			query = query.Order(orderBy)
		}

		// Bangun filter dari query string
		whereClause, args, err := filters.BuildFilters(ctx, true)
		if err != nil {
			return err
		}
		if whereClause != "" {
			query = query.Where(whereClause, args...)
		}

		return query.Find(models).Error
	}

	// Native SQL
	if database.SQLDB == nil {
		return sql.ErrConnDone
	}

	var model T
	table := GetTableName(&model)
	query := fmt.Sprintf("SELECT * FROM %s", table)

	whereClause, args, err := filters.BuildFilters(ctx, false)
	if err != nil {
		return err
	}
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	if orderBy != "" {
		query += fmt.Sprintf(" ORDER BY %s", orderBy)
	}
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := database.SQLDB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var item T
		dest, err := scanRowDestinations(&item)
		if err != nil {
			return fmt.Errorf("error scanning row destinations: %w", err)
		}
		if err := rows.Scan(dest...); err != nil {
			return fmt.Errorf("error scanning row: %w", err)
		}
		*models = append(*models, item)
	}

	return nil
}

// scanRowDestinations returns a slice of addresses of the fields of the given struct that can be set.
// It returns an error if the given model is not a non-nil pointer to a struct.
// The returned slice is suitable for passing to the Scan method of a sql.Row.
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

// GetTableName returns the name of the database table for the given model.
// If the model implements the Tabler interface, the table name is obtained
// from the TableName method.
// Otherwise, the table name is obtained by converting the model name to
// snake case and appending "s".
// For example, the model named "User" is mapped to the table named "users".
func GetTableName[T any](model *T) string {
	if t, ok := any(model).(Tabler); ok {
		return t.TableName()
	}
	typ := reflect.TypeOf(model).Elem()
	return ToSnakeCase(typ.Name()) + "s" // fallback User -> users
}

// GetModelByID returns the model with the given ID from the database.
// If the given model implements the Tabler interface, the table name is obtained
// from the TableName method.
// Otherwise, the table name is obtained by converting the model name to
// snake case and appending "s".
// For example, the model named "User" is mapped to the table named "users".
// If the model does not exist in the database, GetModelByID returns
// sql.ErrNoRows.
// If the database connection is not available, GetModelByID returns
// sql.ErrConnDone.
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

// scanRowIntoStruct scans a single SQL row into the provided model struct.
// The model must be a non-nil pointer to a struct with fields that match the
// database column names via `db` tags. Fields with a `db` tag set to "-" or
// not set at all are ignored. If the row does not exist, it returns an error
// indicating the record was not found. If scanning fails, it returns an error
// with details about the scan failure.

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

// UpdateModelByIDWithMap updates a single record in the database with the given ID.
// The updatedFields map specifies the fields to update and their new values.
// If the database connection uses GORM, it will use GORM's Update method.
// Otherwise, it will generate an UPDATE query using the given map.
// The updated_at field will automatically be set to the current time if it is not
// present in the map.
// If the record is not found, UpdateModelByIDWithMap returns an error with a message
// indicating the record was not found. If the update fails, it returns an error with
// details about the failure.
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

// UpdateModelByID updates a single record in the database with the given ID.
// The updated fields are set from the given model.
// If the database connection uses GORM, it will use GORM's Update method.
// Otherwise, it will generate an UPDATE query using reflection.
// The updated_at field will automatically be set to the current time if it is not
// present in the model.
// If the record is not found, UpdateModelByID returns an error with a message
// indicating the record was not found. If the update fails, it returns an error with
// details about the failure.
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

// hasDeletedAt checks if the given model has a field named "DeletedAt" of type time.Time.
// It checks the model's fields using reflection and returns true if the field exists, false otherwise.
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

// DeleteModelByID deletes a single record from the database with the given ID.
// It checks if the given model has a field named "DeletedAt" of type time.Time.
// If it does, it will perform a soft delete by setting the "DeletedAt" field to the current time.
// If it doesn't, it will perform a hard delete.
// If the database connection is not available, DeleteModelByID returns sql.ErrConnDone.
// If the delete operation fails, it returns an error with details about the failure.
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

// FindOneByField retrieves a single record from the database that matches the given conditions.
// The conditions must be provided as key-value pairs. For example, to find a user with a specific
// email and username, you would call: FindOneByField(&user, "email", emailValue, "username", usernameValue).
// If GORM is enabled, it uses GORM's querying capabilities. Otherwise, it uses native SQL.
// The function returns sql.ErrConnDone if no database connection is available.
// If the record is found, it populates the provided model with the record's data.
// If no record matches the conditions, it returns an error indicating the record was not found.

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
