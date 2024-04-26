package metrics_test

import (
	"testing"

	"github.com/djpiper28/cards-against-humanity/backend/metrics"
	"github.com/stretchr/testify/assert"
)

const PrometheusRegEx = "[a-zA-Z_]+ [0-9]+\\n?"

func TestMetrics(t *testing.T) {
	metrics := metrics.GetMetrics()

	assert.NotEmpty(t, metrics)
	assert.Regexp(t, PrometheusRegEx, metrics)
}
