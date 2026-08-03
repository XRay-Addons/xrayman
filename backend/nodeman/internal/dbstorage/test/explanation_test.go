package dbstoragetest

import (
	"fmt"
	"strings"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
)

type Explanations struct {
	mode  ExplainMode
	name  string
	expl  []string
	start time.Time
}

func (e *Explanations) Add(s string) {
	e.expl = append(e.expl, s)
}

func (e *Explanations) Metrics() (*ExplainMetrics, error) {
	if e.mode == ExplainNone || e.mode == ExplainText {
		return nil, xerr.New("non-explained mode")
	}
	return parseRequestExpl(e.name, e.expl)
}

func (e *Explanations) Text() string {
	return fmt.Sprintf("%s\n:%s", e.name, strings.Join(e.expl, "\n"))
}
