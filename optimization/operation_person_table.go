package optimization

import (
	"accounting/data"
	"database/sql"

	//"log"

	"github.com/pkg/errors"
)

type PersonIdLevel struct {
	Id    int16 // Ид человека.
	Level int8  // Уровень квалификации.
}

type OperationIdLevel struct {
	Id    int16 // Ид операции.
	Level int8  // Уровень квалификации.
}

type Qualification struct {
	Person    int16 // Ид человека.
	Operation int16 // Ид операции.
	Level     int8  // Уровень квалификации.
}

// Структура хранящая информацию о квалификации персонала.
type OperationPersonTable struct {
	person    map[int16]int
	operation map[int16]int
	table     [][]int8 // opt.table[personNumber][operationNumber] = val.Level
}

func lookMap16(mp *map[int16]int, value int) (int16, bool) {
	for key, val := range *mp {
		if val == value {
			return key, true
		}
	}
	return 0, false
}

// Выбор всех ид из таблицы "tableName".
func SelectId16(db *sql.DB, tableName string) ([]int16, error) {
	arr := make([]int16, 0)
	if err := (func() error {
		if err := db.Ping(); err != nil {
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		QwStr := data.SelectId(tableName)
		rows, err := db.Query(QwStr)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorQueryDB+QwStr)
		}
		defer rows.Close()
		var id int16
		for rows.Next() {
			err := rows.Scan(&id)
			if err != nil {
				return errors.Wrap(err, data.S.ErrorDecryptRow)
			}
			arr = append(arr, id)
		}
		return nil
	}()); err != nil {
		return arr, errors.Wrapf(err, "In SelectId16"+tableName)
	}
	return arr, nil
}

// Выбор уровня квалификации из БД.
func SelectQualification(db *sql.DB) ([]Qualification, error) {
	arr := make([]Qualification, 0)
	if err := (func() error {
		if err := db.Ping(); err != nil {
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		QwStr := data.SelectQualification(nil, nil, nil)
		rows, err := db.Query(QwStr)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorQueryDB+QwStr)
		}
		defer rows.Close()
		for rows.Next() {
			pol := Qualification{}
			err := rows.Scan(&pol.Person, &pol.Operation, &pol.Level)
			if err != nil {
				return errors.Wrap(err, data.S.ErrorDecryptRow)
			}
			arr = append(arr, pol)
		}
		return nil
	}()); err != nil {
		return arr, errors.Wrapf(err, "In SelectQualification") //GO-TO строка
	}
	return arr, nil
}

func NewQualificationTable(db *sql.DB) (OperationPersonTable, error) {
	opt := OperationPersonTable{}
	if err := (func() error {
		person, err := SelectId16(db, "Person")
		if err != nil {
			return err
		}
		if len(person) == 0 {
			return errors.New("len(person) == 0") //GO-TO добавить обработку ошибок
		}
		opt.person = make(map[int16]int, len(person))
		for index, val := range person {
			opt.person[val] = index + 1
		}

		operation, err := SelectId16(db, "Operation")
		if err != nil {
			return err
		}
		if len(operation) == 0 {
			return errors.New("len(operation) == 0") //GO-TO добавить обработку ошибок
		}
		opt.operation = make(map[int16]int, len(operation))
		for index, val := range operation {
			opt.operation[val] = index + 1
		}

		cells := make([]int8, len(person)*len(operation))
		opt.table = make([][]int8, len(person))
		for index := range opt.table {
			opt.table[index], cells = cells[:len(operation)], cells[len(operation):]
		}
		level, err := SelectQualification(db)
		if err != nil {
			return err
		}
		if len(level) == 0 {
			return errors.New("len(level) == 0") //GO-TO добавить обработку ошибок
		}
		for _, val := range level {
			personNumber := opt.person[val.Person] - 1
			operationNumber := opt.operation[val.Operation] - 1
			opt.table[personNumber][operationNumber] = val.Level
		}
		return nil
	}()); err != nil {
		return opt, errors.Wrapf(err, "In NewQualificationTable") //GO-TO строка
	}
	return opt, nil
}

// Получения уровня квалифкации данного человека по данной операции.
func (opt OperationPersonTable) GetLevel(personId, operationId int16) int8 {
	if personNumber, operationNumber := opt.person[personId], opt.operation[operationId]; personNumber == 0 || operationNumber == 0 {
		//GO-TO print error
		return 0
	} else {
		return opt.table[personNumber-1][operationNumber-1]
	}
}

// Получение списка людей, у которых есть какая-либо квалификация в данной операции.
func (opt OperationPersonTable) GetPersons(operation int16) []int16 {
	arr := make([]int16, 0, len(opt.person))
	for index, number := range opt.person {
		if opt.table[number-1][operation-1] > 0 {
			arr = append(arr, index)
		}
	}
	return arr
}
