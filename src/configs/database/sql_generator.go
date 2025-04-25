package database

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// GenerateCreateTableSQL generates a SQL string for creating a table based on the given model.
// The `tableName` parameter is the name of the table to be created.
// The `model` parameter must be a struct.
// The function will generate a SQL string for creating the table with the given name,
// with columns and constraints according to the struct's fields and their tags.
// The function will also generate a trigger SQL for automatically setting the `updated_at` field
// to the current time on each update, if the struct has a field with the "db" tag set to "updated_at".
// The function will panic if the `DB_DRIVER` environment variable is not set to "postgres" or "mysql".
func GenerateCreateTableSQL(tableName string, model interface{}) string {
	engine := strings.ToLower(os.Getenv("DB_DRIVER")) // "postgres" atau "mysql"
	if engine != "postgres" && engine != "mysql" {
		panic("DB_DRIVER must be either 'postgres' or 'mysql'")
	}

	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Struct {
		panic("model must be a struct")
	}

	var fields []string
	var constraints []string
	var hasUpdatedAt bool

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")
		if tag == "-" || tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		col := parts[0]
		opts := parts[1:]

		if col == "" {
			continue
		}

		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		colType, isAutoUpdate := GoTypeToSQLType(engine, fieldType, opts)
		if isAutoUpdate {
			hasUpdatedAt = true
		}

		for _, opt := range opts {
			if strings.HasPrefix(opt, "foreign:") {
				ref := strings.TrimPrefix(opt, "foreign:")
				constraints = append(constraints, fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s", quoteIdent(engine, col), ref))
			}
		}

		fields = append(fields, fmt.Sprintf("%s %s", quoteIdent(engine, col), colType))
	}

	allDefs := append(fields, constraints...)
	createTableSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n%s\n);", quoteIdent(engine, tableName), strings.Join(allDefs, ",\n"))

	if hasUpdatedAt {
		if engine == "postgres" {
			createTableSQL += "\n" + GenerateUpdatedAtTriggerSQLPostgres(tableName)
		} else {
			createTableSQL += "\n" + GenerateUpdatedAtTriggerSQLMySQL(tableName)
		}
	}

	return createTableSQL
}

func quoteIdent(engine, ident string) string {
	if engine == "postgres" {
		return `"` + ident + `"`
	}
	return "`" + ident + "`"
}

func GoTypeToSQLType(engine string, goType reflect.Type, opts []string) (string, bool) {
	flags := map[string]bool{}
	var defaultVal string
	isAutoUpdate := false

	for _, opt := range opts {
		switch {
		case opt == "autoupdate":
			isAutoUpdate = true
		case strings.HasPrefix(opt, "default:"):
			defaultVal = strings.TrimPrefix(opt, "default:")
		default:
			flags[opt] = true
		}
	}

	var base string
	switch goType.Kind() {
	case reflect.Int, reflect.Int64:
		if flags["serial"] {
			if engine == "postgres" {
				base = "SERIAL"
			} else {
				base = "BIGINT AUTO_INCREMENT"
			}
		} else {
			base = "BIGINT"
		}
	case reflect.Uint, reflect.Uint64:
		base = "BIGINT UNSIGNED"
	case reflect.String:
		base = "TEXT"
	case reflect.Bool:
		base = "BOOLEAN"
	case reflect.Float32, reflect.Float64:
		base = "DOUBLE"
	case reflect.Struct:
		if goType.PkgPath() == "time" && goType.Name() == "Time" {
			base = "TIMESTAMP"
		} else {
			base = "TEXT"
		}
	default:
		base = "TEXT"
	}

	var parts []string
	parts = append(parts, base)

	if defaultVal != "" {
		if defaultVal == "now" && goType.PkgPath() == "time" && goType.Name() == "Time" {
			if engine == "postgres" {
				parts = append(parts, "DEFAULT CURRENT_TIMESTAMP")
			} else if engine == "mysql" {
				parts = append(parts, "DEFAULT CURRENT_TIMESTAMP")
			}
		} else {
			parts = append(parts, fmt.Sprintf("DEFAULT '%s'", defaultVal))
		}
	}

	if flags["notnull"] {
		parts = append(parts, "NOT NULL")
	}
	if flags["unique"] {
		parts = append(parts, "UNIQUE")
	}
	if flags["primary"] {
		parts = append(parts, "PRIMARY KEY")
	}

	return strings.Join(parts, " "), isAutoUpdate
}

func GenerateUpdatedAtTriggerSQLPostgres(tableName string) string {
	return fmt.Sprintf(`
CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
	NEW.updated_at = NOW();
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS set_updated_at ON %s;
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON %s
FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_at();`, tableName, tableName)
}

func GenerateUpdatedAtTriggerSQLMySQL(tableName string) string {
	return fmt.Sprintf(`
DROP TRIGGER IF EXISTS set_updated_at_%s;
CREATE TRIGGER set_updated_at_%s
BEFORE UPDATE ON %s
FOR EACH ROW
BEGIN
	SET NEW.updated_at = CURRENT_TIMESTAMP;
END;`, tableName, tableName, tableName)
}
