package qbconsts

type QueryType string

const (
	SELECT QueryType = "SELECT"
	INSERT QueryType = "INSERT"
	UPDATE QueryType = "UPDATE"
	DELETE QueryType = "DELETE"
)

type QueryOrder string

const (
	ASCENDING  QueryOrder = "ASC"
	DESCENDING QueryOrder = "DESC"
)

type QueryConnector string

const (
	AND QueryConnector = "AND"
	OR  QueryConnector = "OR"
	NOT QueryConnector = "NOT"
)
