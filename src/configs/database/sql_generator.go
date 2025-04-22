package database

import (
	"fmt"
	"reflect"
	"strings"
)

func GenerateCreateTableSQL(tableName string, model interface{}) string {
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

		// Unwrap pointer
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		colType, isAutoUpdate := GoTypeToPostgresType(fieldType, opts)
		if isAutoUpdate {
			hasUpdatedAt = true
		}

		// Foreign key
		for _, opt := range opts {
			if strings.HasPrefix(opt, "foreign:") {
				ref := strings.TrimPrefix(opt, "foreign:")
				constraints = append(constraints, fmt.Sprintf("FOREIGN KEY (\"%s\") REFERENCES %s", col, ref))
			}
		}

		fields = append(fields, fmt.Sprintf(`"%s" %s`, col, colType))
	}

	allDefs := append(fields, constraints...)
	createTableSQL := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
%s
);`, tableName, strings.Join(allDefs, ",\n"))

	if hasUpdatedAt {
		triggerSQL := GenerateUpdatedAtTriggerSQL(tableName)
		createTableSQL += "\n" + triggerSQL
	}

	return createTableSQL
}

func GoTypeToPostgresType(goType reflect.Type, opts []string) (string, bool) {
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
			base = "SERIAL"
		} else {
			base = "BIGINT"
		}
	case reflect.Uint, reflect.Uint64:
		if flags["serial"] {
			base = "BIGSERIAL"
		} else {
			base = "BIGINT"
		}
	case reflect.String:
		base = "TEXT"
	case reflect.Bool:
		base = "BOOLEAN"
	case reflect.Float32, reflect.Float64:
		base = "DOUBLE PRECISION"
	case reflect.Struct:
		if goType.PkgPath() == "time" && goType.Name() == "Time" {
			base = "TIMESTAMP"
		} else {
			base = "TEXT"
		}
	default:
		base = "TEXT"
	}

	// Format SQL dalam urutan yang benar
	var parts []string
	parts = append(parts, base)

	if defaultVal != "" {
		if defaultVal == "now" && goType.PkgPath() == "time" && goType.Name() == "Time" {
			parts = append(parts, "DEFAULT CURRENT_TIMESTAMP")
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

// func GoTypeToPostgresType(goType reflect.Type, opts []string) (string, bool) {
// 	flags := map[string]bool{}
// 	var defaultVal string
// 	var isAutoUpdate bool

// 	// Parse tags for default and autoupdate
// 	for _, opt := range opts {
// 		if strings.HasPrefix(opt, "default:") {
// 			defaultVal = strings.TrimPrefix(opt, "default:")
// 		} else if opt == "autoupdate" {
// 			isAutoUpdate = true
// 		} else {
// 			flags[opt] = true
// 		}
// 	}

// 	var base string
// 	switch goType.Kind() {
// 	case reflect.Int, reflect.Int64:
// 		if flags["serial"] {
// 			base = "SERIAL"
// 		} else {
// 			base = "BIGINT"
// 		}
// 	case reflect.Uint, reflect.Uint64:
// 		if flags["serial"] {
// 			base = "BIGSERIAL"
// 		} else {
// 			base = "BIGINT"
// 		}
// 	case reflect.String:
// 		base = "TEXT"
// 	case reflect.Bool:
// 		base = "BOOLEAN"
// 	case reflect.Float32, reflect.Float64:
// 		base = "DOUBLE PRECISION"
// 	case reflect.Struct:
// 		if goType.PkgPath() == "time" && goType.Name() == "Time" {
// 			base = "TIMESTAMP"
// 		} else {
// 			base = "TEXT"
// 		}
// 	default:
// 		base = "TEXT"
// 	}

// 	var constraints []string
// 	if flags["primary"] {
// 		constraints = append(constraints, "PRIMARY KEY")
// 	}
// 	if flags["notnull"] {
// 		constraints = append(constraints, "NOT NULL")
// 	}
// 	if flags["unique"] {
// 		constraints = append(constraints, "UNIQUE")
// 	}
// 	if defaultVal != "" {
// 		if defaultVal == "now" && base == "TIMESTAMP" {
// 			constraints = append(constraints, "DEFAULT NOW()")
// 		} else {
// 			constraints = append(constraints, fmt.Sprintf("DEFAULT '%s'", defaultVal))
// 		}
// 	}
// 	if len(constraints) > 0 {
// 		base += " " + strings.Join(constraints, " ")
// 	}

// 	return base, isAutoUpdate
// }

func GenerateUpdatedAtTriggerSQL(tableName string) string {
	triggerFunc := fmt.Sprintf(`
		CREATE OR REPLACE FUNCTION trigger_set_updated_at()
		RETURNS TRIGGER AS $$
		BEGIN
		NEW.updated_at = NOW();
		RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;`)

	trigger := fmt.Sprintf(`
		CREATE TRIGGER set_updated_at
		BEFORE UPDATE ON %s
		FOR EACH ROW
		EXECUTE FUNCTION trigger_set_updated_at();`, tableName)

	return triggerFunc + "\n" + trigger
}
