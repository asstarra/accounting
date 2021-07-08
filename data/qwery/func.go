package qwery

import (
	. "accounting/data/constants"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
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
		return t.Format(TimeLayoutMySql)
	case time.Time:
		return t.Format(TimeLayoutMySql)
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

// Превращение нескольких интерфейсов в массив строк.
func ToStrs(arr []interface{}) []string {
	strs := make([]string, len(arr), len(arr))
	for argNum, arg := range arr {
		strs[argNum] = ToStr(arg)
	}
	return strs
}

func ToIntfs(strs []string) []interface{} {
	arr := make([]interface{}, len(strs), len(strs))
	for argNum, arg := range strs {
		arr[argNum] = arg
	}
	return arr
}

func Printf(format string, arr ...interface{}) string {
	var strs = ToIntfs(ToStrs(arr))
	return fmt.Sprintf(format, strs...)
}

func Wrapf(err error, format string, arr ...interface{}) error {
	return errors.Wrap(err, Printf(format, arr))
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}
