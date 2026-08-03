package dbstoragetest

import (
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

type ExplainMetrics struct {
	name          string
	ExecutionTime time.Duration

	SeqScans         int
	Sorts            int
	NestedLoops      int
	BitmapHeapScans  int
	BitmapIndexScans int

	TotalNodes int
}

func (m *ExplainMetrics) Print(log *zap.Logger) {
	log.Info(m.name,
		zap.Duration("Exec Time", m.ExecutionTime),
		zap.Int("Seq Scans", m.SeqScans),
		zap.Int("Nested Loops", m.NestedLoops),
		zap.Int("Bitmap Heap Scans", m.BitmapHeapScans),
		zap.Int("Bitmap Index Scans", m.BitmapIndexScans),
		zap.Int("Total Nodes", m.TotalNodes),
	)
}

func parseRequestExpl(name string, inputs []string) (*ExplainMetrics, error) {
	m := &ExplainMetrics{
		name: name,
	}
	for _, input := range inputs {
		var raw []map[string]any

		if err := json.Unmarshal([]byte(input), &raw); err != nil {
			return nil, err
		}

		for _, item := range raw {

			// execution time (ms → duration)
			if v, ok := item["Execution Time"].(float64); ok {
				m.ExecutionTime += msToDuration(v)
			}

			if plan, ok := item["Plan"].(map[string]any); ok {
				walk(plan, m)
			}
		}
	}

	return m, nil
}

func msToDuration(ms float64) time.Duration {
	return time.Duration(ms * float64(time.Millisecond))
}

func walk(node map[string]any, m *ExplainMetrics) {
	if node == nil {
		return
	}

	m.TotalNodes++

	switch node["Node Type"] {
	case "Seq Scan":
		m.SeqScans++
	case "Sort":
		m.Sorts++
	case "Nested Loop":
		m.NestedLoops++
	case "Bitmap Heap Scan":
		m.BitmapHeapScans++
	case "Bitmap Index Scan":
		m.BitmapIndexScans++
	}

	if children, ok := node["Plans"].([]any); ok {
		for _, c := range children {
			if cm, ok := c.(map[string]any); ok {
				walk(cm, m)
			}
		}
	}
}
