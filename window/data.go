package window

type IdTitle struct {
	Id    int64
	Title string
}

type EntityRecChild struct {
	IdTitle
	Count int
}

type Entity struct {
	Id             int64
	Title          string
	Type           int64
	Specification  string
	ProductionLine bool
	Note           string
	Children       *[]*EntityRecChild
}

func Empty(s string) string {
	return "'" + s + "'"
}
