package filters

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func parseFilterParam(param string) (field string, operator string) {
	parts := strings.Split(param, "[")
	if len(parts) != 2 {
		return "", ""
	}

	field = parts[0]
	operator = strings.TrimRight(parts[1], "]")
	return field, operator
}

func BuildFilters(ctx *gin.Context, useGORM bool) (string, []interface{}, error) {
	var filters []string
	var args []interface{}
	argIndex := 1

	// Iterasi semua query parameter
	for param, values := range ctx.Request.URL.Query() {
		if !strings.Contains(param, "[") || len(values) == 0 {
			continue
		}

		field, operator := parseFilterParam(param)
		if field == "" {
			return "", nil, fmt.Errorf("invalid filter parameter: %s", param)
		}

		value := values[0]

		switch operator {
		case "like", "ilike":
			valueLower := strings.ToLower(value)

			if useGORM {
				driver := os.Getenv("DB_DRIVER")
				if driver == "postgres" && operator == "ilike" {
					filters = append(filters, fmt.Sprintf("%s ILIKE ?", field))
					args = append(args, "%"+valueLower+"%")
				} else {
					// Untuk MySQL/GORM, pakai LOWER agar aman dari collation
					filters = append(filters, fmt.Sprintf("LOWER(%s) LIKE ?", field))
					args = append(args, "%"+valueLower+"%")
				}
			} else {
				driver := os.Getenv("DB_DRIVER")
				op := "LIKE"
				placeholder := "?"
				if driver == "postgres" && operator == "ilike" {
					op = "ILIKE"
					placeholder = fmt.Sprintf("$%d", argIndex)
				} else if driver == "postgres" {
					placeholder = fmt.Sprintf("$%d", argIndex)
				}
				filters = append(filters, fmt.Sprintf("LOWER(%s) %s %s", field, op, placeholder))
				args = append(args, "%"+valueLower+"%")
				argIndex++
			}
		case "moreThan":
			val, err := strconv.Atoi(value)
			if err != nil {
				return "", nil, fmt.Errorf("invalid value for 'moreThan': %s", value)
			}
			if useGORM {
				filters = append(filters, fmt.Sprintf("%s > ?", field))
				args = append(args, val)
			} else {
				filters = append(filters, fmt.Sprintf("%s > $%d", field, argIndex))
				args = append(args, val)
				argIndex++
			}

		case "lessThan":
			val, err := strconv.Atoi(value)
			if err != nil {
				return "", nil, fmt.Errorf("invalid value for 'lessThan': %s", value)
			}
			if useGORM {
				filters = append(filters, fmt.Sprintf("%s < ?", field))
				args = append(args, val)
			} else {
				filters = append(filters, fmt.Sprintf("%s < $%d", field, argIndex))
				args = append(args, val)
				argIndex++
			}

		case "equals":
			if useGORM {
				filters = append(filters, fmt.Sprintf("%s = ?", field))
				args = append(args, value)
			} else {
				filters = append(filters, fmt.Sprintf("%s = $%d", field, argIndex))
				args = append(args, value)
				argIndex++
			}

		case "notEquals":
			if useGORM {
				filters = append(filters, fmt.Sprintf("%s != ?", field))
				args = append(args, value)
			} else {
				filters = append(filters, fmt.Sprintf("%s != $%d", field, argIndex))
				args = append(args, value)
				argIndex++
			}

		case "greaterThanOrEqual":
			if useGORM {
				filters = append(filters, fmt.Sprintf("%s >= ?", field))
				args = append(args, value)
			} else {
				filters = append(filters, fmt.Sprintf("%s >= $%d", field, argIndex))
				args = append(args, value)
				argIndex++
			}

		case "lessThanOrEqual":
			if useGORM {
				filters = append(filters, fmt.Sprintf("%s <= ?", field))
				args = append(args, value)
			} else {
				filters = append(filters, fmt.Sprintf("%s <= $%d", field, argIndex))
				args = append(args, value)
				argIndex++
			}

		case "in":
			valList := strings.Split(value, ",")
			if useGORM {
				filters = append(filters, fmt.Sprintf("%s IN (?)", field))
				args = append(args, valList)
			} else {
				placeholders := []string{}
				for _, val := range valList {
					placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
					args = append(args, val)
					argIndex++
				}
				filters = append(filters, fmt.Sprintf("%s IN (%s)", field, strings.Join(placeholders, ",")))
			}

		case "notIn":
			valList := strings.Split(value, ",")
			if useGORM {
				filters = append(filters, fmt.Sprintf("%s NOT IN (?)", field))
				args = append(args, valList)
			} else {
				placeholders := []string{}
				for _, val := range valList {
					placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
					args = append(args, val)
					argIndex++
				}
				filters = append(filters, fmt.Sprintf("%s NOT IN (%s)", field, strings.Join(placeholders, ",")))
			}

		default:
			return "", nil, fmt.Errorf("unsupported operator: %s", operator)
		}
	}

	whereClause := strings.Join(filters, " AND ")
	return whereClause, args, nil
}
