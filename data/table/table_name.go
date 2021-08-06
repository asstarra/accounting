package table

import (
	"accounting/data/text"
)

type TableName int8

const (
	TableEntityType TableName = iota
	TableEntity
	TableEntityRec
	TableLine
	TableLineId
	TableMarkedDetail
	TableStatusType
	TableStatus
	TableStatusDetail
	TablePerson
	TableOperation
)

func (tn TableName) String() string {
	switch tn {
	case TableEntityType:
		return "EntityType"
	case TableEntity:
		return "Entity"
	case TableEntityRec:
		return "EntityRec"
	case TableLine:
		return "MarkingLine"
	case TableLineId:
		return "Marking"
	case TableMarkedDetail:
		return "MarkedDetail"
	case TableStatusType:
		return "StatusType"
	case TableStatus:
		return "Status"
	case TableStatusDetail:
		return "StatusDetail"
	case TablePerson:
		return "Person"
	case TableOperation:
		return "Operation"
	}
	return ""
}

func (tn TableName) Heading() string { // TO-DO
	switch tn {
	case TableEntityType:
		return text.T.HeadingEntityType
	case TableStatusType:
		return "StatusType"
	case TablePerson:
		return "Person"
	case TableOperation:
		return "Operation"
	}
	return ""
}

func (tn TableName) Title() string { // TO-DO
	switch tn {
	case TableEntityType:
		return text.T.ColumnTitle
	case TableStatusType:
		return text.T.ColumnTitle
	case TablePerson:
		return "Person"
	case TableOperation:
		return text.T.ColumnTitle
	}
	return ""
}
