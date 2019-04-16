package metrics

import (
	"testing"

	"github.com/m-lab/go/prometheusx/promtest"
)

func TestMetrics(t *testing.T) {
	// Label the metrics and set them to a value to ensure they show up in the output.
	GatherRuns.WithLabelValues("test").Add(1)
	GatherErrors.WithLabelValues("test").Add(1)
	GatherRuntime.WithLabelValues("test").Observe(1)
	promtest.LintMetrics(t)
}
