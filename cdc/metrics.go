// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package cdc

import (
	"github.com/amyangfei/tiflow/cdc/entry"
	"github.com/amyangfei/tiflow/cdc/kv"
	"github.com/amyangfei/tiflow/cdc/owner"
	"github.com/amyangfei/tiflow/cdc/processor"
	tablepipeline "github.com/amyangfei/tiflow/cdc/processor/pipeline"
	"github.com/amyangfei/tiflow/cdc/puller"
	redowriter "github.com/amyangfei/tiflow/cdc/redo/writer"
	"github.com/amyangfei/tiflow/cdc/sink"
	"github.com/amyangfei/tiflow/cdc/sorter"
	"github.com/amyangfei/tiflow/cdc/sorter/leveldb"
	"github.com/amyangfei/tiflow/cdc/sorter/memory"
	"github.com/amyangfei/tiflow/cdc/sorter/unified"
	"github.com/amyangfei/tiflow/pkg/actor"
	"github.com/amyangfei/tiflow/pkg/db"
	"github.com/amyangfei/tiflow/pkg/etcd"
	"github.com/amyangfei/tiflow/pkg/orchestrator"
	"github.com/prometheus/client_golang/prometheus"
)

var registry = prometheus.NewRegistry()

func init() {
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	registry.MustRegister(prometheus.NewGoCollector())

	kv.InitMetrics(registry)
	puller.InitMetrics(registry)
	sink.InitMetrics(registry)
	entry.InitMetrics(registry)
	processor.InitMetrics(registry)
	tablepipeline.InitMetrics(registry)
	owner.InitMetrics(registry)
	etcd.InitMetrics(registry)
	initServerMetrics(registry)
	actor.InitMetrics(registry)
	orchestrator.InitMetrics(registry)
	// Sorter metrics
	sorter.InitMetrics(registry)
	memory.InitMetrics(registry)
	unified.InitMetrics(registry)
	leveldb.InitMetrics(registry)
	redowriter.InitMetrics(registry)
	db.InitMetrics(registry)
}
