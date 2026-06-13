package dbstoragetest

type ExplainMode int

const (
	ExplainNone ExplainMode = iota + 1
	ExplainText
	ExplainJson
	ExplainAnalyze
)
