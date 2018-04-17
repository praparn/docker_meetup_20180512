// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package processors

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kubernetes-incubator/metrics-server/metrics/core"
)

func TestClusterAggregate(t *testing.T) {
	batch := core.DataBatch{
		Timestamp: time.Now(),
		MetricSets: map[string]*core.MetricSet{
			core.PodKey("ns1", "pod1"): {
				Labels: map[string]string{
					core.LabelMetricSetType.Key: core.MetricSetTypeNamespace,
					core.LabelNamespaceName.Key: "ns1",
				},
				MetricValues: map[string]core.MetricValue{
					"m1": {
						ValueType:  core.ValueInt64,
						MetricType: core.MetricGauge,
						IntValue:   10,
					},
					"m2": {
						ValueType:  core.ValueInt64,
						MetricType: core.MetricGauge,
						IntValue:   222,
					},
				},
			},

			core.PodKey("ns1", "pod2"): {
				Labels: map[string]string{
					core.LabelMetricSetType.Key: core.MetricSetTypeNamespace,
					core.LabelNamespaceName.Key: "ns1",
				},
				MetricValues: map[string]core.MetricValue{
					"m1": {
						ValueType:  core.ValueInt64,
						MetricType: core.MetricGauge,
						IntValue:   100,
					},
					"m3": {
						ValueType:  core.ValueInt64,
						MetricType: core.MetricGauge,
						IntValue:   30,
					},
				},
			},
		},
	}
	processor := ClusterAggregator{
		MetricsToAggregate: []string{"m1", "m3"},
	}
	result, err := processor.Process(&batch)
	assert.NoError(t, err)
	cluster, found := result.MetricSets[core.ClusterKey()]
	assert.True(t, found)

	m1, found := cluster.MetricValues["m1"]
	assert.True(t, found)
	assert.Equal(t, int64(110), m1.IntValue)

	m3, found := cluster.MetricValues["m3"]
	assert.True(t, found)
	assert.Equal(t, int64(30), m3.IntValue)
}
