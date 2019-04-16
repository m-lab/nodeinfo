package metrics

import (
	"testing"

	"github.com/m-lab/go/prometheusx/promtest"
)

func TestMetrics(t *testing.T) {
	GatherRuns.WithLabelValues("test").Add(1)
	GatherErrors.WithLabelValues("test").Add(1)
	GatherRuntime.WithLabelValues("test").Observe(1)
	promtest.LintMetrics(t)
}
