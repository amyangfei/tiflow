// Copyright 2023 PingCAP, Inc.
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

package observer

import (
	"context"
	"strings"

	"github.com/pingcap/tiflow/cdc/contextutil"
	"github.com/pingcap/tiflow/pkg/config"
	"github.com/pingcap/tiflow/pkg/sink"
	pmysql "github.com/pingcap/tiflow/pkg/sink/mysql"
)

// Observer defines an interface of downstream performance observer.
type Observer interface {
	// Tick is called periodically, Observer fetches performance metrics and
	// records them in each Tick.
	Tick(ctx context.Context) error
}

// NewObserverOpt represents available options when creating a new observer.
type NewObserverOpt struct {
	dbConnFactory pmysql.Factory
}

// NewObserverOption configures NewObserverOpt.
type NewObserverOption func(*NewObserverOpt)

// WithDBConnFactory specifies factory to create db connection.
func WithDBConnFactory(factory pmysql.Factory) NewObserverOption {
	return func(opt *NewObserverOpt) {
		opt.dbConnFactory = factory
	}
}

// NewObserver creates a new Observer
func NewObserver(
	ctx context.Context,
	sinkURIStr string,
	replCfg *config.ReplicaConfig,
	opts ...NewObserverOption,
) (Observer, error) {
	options := &NewObserverOpt{dbConnFactory: pmysql.CreateMySQLDBConn}
	for _, opt := range opts {
		opt(options)
	}
	sinkURI, err := config.GetSinkURIAndAdjustConfigWithSinkURI(sinkURIStr, replCfg)
	if err != nil {
		return nil, err
	}

	scheme := strings.ToLower(sinkURI.Scheme)
	if !sink.IsMySQLCompatibleScheme(scheme) {
		return NewDummyObserver(), nil
	}

	changefeedID := contextutil.ChangefeedIDFromCtx(ctx)
	cfg := pmysql.NewConfig()
	err = cfg.Apply(ctx, changefeedID, sinkURI, replCfg)
	if err != nil {
		return nil, err
	}

	dsnStr, err := pmysql.GenerateDSN(ctx, sinkURI, cfg, options.dbConnFactory)
	if err != nil {
		return nil, err
	}
	db, err := options.dbConnFactory(ctx, dsnStr)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(2)

	isTiDB, err := pmysql.CheckIsTiDB(ctx, db)
	if err != nil {
		return nil, err
	}
	if isTiDB {
		return NewTiDBObserver(db), nil
	}
	return NewDummyObserver(), nil
}
