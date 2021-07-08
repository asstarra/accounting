package log

var (
	Panic   = "PANIC!"
	Error   = "ERROR!"
	Warning = "WARNING!"
	Info    = "INFO:"
	Debug   = "DEBUG:"

	BeginWindow  = "INFO: BEGIN window %s"
	InitWindow   = "INFO: INIT window %s"
	CreateWindow = "INFO: CREATE window %s"
	RunWindow    = "INFO: RUN window %s"
	EndWindow    = "INFO: END window %s  cmd %v"

	Entity        = "ENTITY"
	Entities      = "ENTITIES"
	EntityRec     = "ENTITY_REC"
	Type          = "TYPE"
	MarkedDetail  = "MARKED_DETAIL"
	MarkedDetails = "MARKED_DETAILS"

	LogOk     = "Ok"
	LogCansel = "Cansel"
	LogAdd    = "Add"
	LogChange = "Change"
	LogDelete = "Delete"
	LogChoose = "Choose"
	LogSearch = "Search"

	LogNotInsertedId = "Not inserted id"
)

var In = struct {
	InEntitiesRunDialog      string
	InEntityRunDialog        string
	InEntityRecRunDialog     string
	InTypeRunDialog          string
	InMarkedDetailRunDialog  string
	InMarkedDetailsRunDialog string

	InSelectEntity            string
	InSelectEntityRecChild    string
	InSelectEntityRec         string
	InSelectId16Title         string
	InSelectMarkingLineNew    string
	InSelectMarkingLineOld    string
	InSelectMarkingLineEntity string
	InSelectMarkedDetails     string

	InSelectPersonTime      string
	InSelectPerson          string
	InSelectDays            string
	InSelectId16            string
	InSelectQualification   string
	InNewQualificationTable string
}{
	InEntitiesRunDialog:      "In EntitiesRunDialog(isChage = %t, IdTitle = %v)",
	InEntityRunDialog:        "In EntityRunDialog(entity = %v)",
	InEntityRecRunDialog:     "In EntityRecRunDialog(child = %v)",
	InTypeRunDialog:          "In TypeRunDialog(tableName = %s)",
	InMarkedDetailRunDialog:  "In MarkedDetailRunDialog(detail = %v)",
	InMarkedDetailsRunDialog: "In MarkedDetailsRunDialog(isChage = %t, parent detail = %v)",

	InSelectEntity:            "In SelectEntities(title = \"%s\", entityType = %d)",
	InSelectEntityRecChild:    "In SelectEntityRecChild(parent = %d)",
	InSelectEntityRec:         "In SelectEntityRec()",
	InSelectId16Title:         "In SelectId16Title(tableName = %s, id = %s, title = %s)",
	InSelectMarkingLineNew:    "In SelectMarkingLineNew()",
	InSelectMarkingLineOld:    "In SelectMarkingLineOld()",
	InSelectMarkingLineEntity: "In SelectMarkingLineEntity(id = %d)",
	InSelectMarkedDetails:     "In SelectMarkedDetails(marking = %v)",

	InSelectPersonTime:      "In SelectPersonTime(person = %s, start = %s, finish = %s, detail = %s, entity = %s, number = %s)",
	InSelectPerson:          "In SelectPerson(start = %s, finish = %s)",
	InSelectDays:            "In SelectDays(start = %s, finish = %s)",
	InSelectId16:            "In SelectId16(tableName = %s)",
	InSelectQualification:   "In SelectQualification()",
	InNewQualificationTable: "In NewQualificationTable()",
}
