package qwery

import (
	"accounting/data"
	"database/sql"
	"fmt"
	"time"
)

func Prefix(strArr []string) string {
	if len(strArr) == 1 {
		return " WHERE "
	} else {
		return " AND "
	}
}

func Merger(strArr []string) string {
	str := ""
	for _, val := range strArr {
		str += val
	}
	return str
}

// Превращение интерфейса в строку.
func ToStr(ptr interface{}) string {
	null := "NULL"
	if ptr == nil {
		return "nil"
	}
	switch t := ptr.(type) {
	case *time.Time:
		return t.Format(data.TimeLayout.MySql) // GO-TO
	case time.Time:
		return t.Format(data.TimeLayout.MySql) // GO-TO
	case *string:
		return *t
	case *float32:
		return ToStr(*t)
	case *float64:
		return ToStr(*t)
	case *int8:
		return ToStr(*t)
	case *int16:
		return ToStr(*t)
	case *int32:
		return ToStr(*t)
	case *int64:
		return ToStr(*t)
	case *int:
		return ToStr(*t)
	case *uint8:
		return ToStr(*t)
	case *uint16:
		return ToStr(*t)
	case *uint32:
		return ToStr(*t)
	case *uint64:
		return ToStr(*t)
	case *uint:
		return ToStr(*t)
	case *bool:
		return ToStr(*t)
	case string, float32, float64, int8, int16, int32, int64, int, uint8, uint16, uint32, uint64, uint, bool:
		return fmt.Sprint(t)

	case *sql.NullTime:
		return ToStr(*t)
	case sql.NullTime:
		if !t.Valid {
			return null
		}
		return ToStr(t.Time)
	case *sql.NullString:
		return ToStr(*t)
	case sql.NullString:
		if !t.Valid {
			return null
		}
		return ToStr(t.String)
	case *sql.NullFloat64:
		return ToStr(*t)
	case sql.NullFloat64:
		if !t.Valid {
			return null
		}
		return ToStr(t.Float64)
	case *sql.NullInt64:
		return ToStr(*t)
	case sql.NullInt64:
		if !t.Valid {
			return null
		}
		return ToStr(t.Int64)
	case *sql.NullInt32:
		return ToStr(*t)
	case sql.NullInt32:
		if !t.Valid {
			return null
		}
		return ToStr(t.Int32)
	case *sql.NullBool:
		return ToStr(*t)
	case sql.NullBool:
		if !t.Valid {
			return null
		}
		return ToStr(t.Bool)
	}
	return fmt.Sprintf("%v", ptr)
}
