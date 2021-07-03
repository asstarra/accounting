package optimization

// import (
// 	// "accounting/data"
// 	"database/sql"
// 	"time"
// 	//"log"
// 	// "github.com/pkg/errors"
// )

// type Details struct {
// 	arr []DetailDB
// 	mp  map[int64]*DetailDB
// }

// func InitDetails(db *sql.DB, startPtr, finishPtr *time.Time) ([]Detail, error) {
// 	arrDet, err := SelectDetail(db, startPtr, finishPtr)
// 	if err != nil {
// 		return []Detail{}, err // GO-TO
// 	}
// 	// mapDetail := make(map[int64]*Detail, len(arrDet))
// 	for _, detDB := range arrDet {
// 		var nowStage int
// 		for nowStage = 0; nowStage < len(detDB.Way) && detDB.Way[nowStage].Check; nowStage++ {
// 		}
// 		// det := Detail{
// 		// 	DetailDB: detDB,
// 		// }
// 	}
// 	return []Detail{}, nil
// }
