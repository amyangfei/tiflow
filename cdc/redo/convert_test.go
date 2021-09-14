// Copyright 2021 PingCAP, Inc.
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

package redo

import (
	"testing"

	timodel "github.com/pingcap/parser/model"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/ticdc/cdc/model"
	"github.com/stretchr/testify/require"
)

func TestRowRedoConvert(t *testing.T) {
	t.Parallel()
	row := &model.RowChangedEvent{
		StartTs:  100,
		CommitTs: 120,
		Table:    &model.TableName{Schema: "test", Table: "table1", TableID: 57},
		PreColumns: []*model.Column{{
			Name:  "a1",
			Type:  mysql.TypeLong,
			Flag:  model.BinaryFlag | model.MultipleKeyFlag | model.HandleKeyFlag,
			Value: int64(1),
		}, {
			Name:  "a2",
			Type:  mysql.TypeVarchar,
			Value: "char",
		}, {
			Name:  "a3",
			Type:  mysql.TypeLong,
			Flag:  model.BinaryFlag | model.MultipleKeyFlag | model.HandleKeyFlag,
			Value: int64(1),
		}, nil},
		Columns: []*model.Column{{
			Name:  "a1",
			Type:  mysql.TypeLong,
			Flag:  model.BinaryFlag | model.MultipleKeyFlag | model.HandleKeyFlag,
			Value: int64(2),
		}, {
			Name:  "a2",
			Type:  mysql.TypeVarchar,
			Value: "char-updated",
		}, {
			Name:  "a3",
			Type:  mysql.TypeLong,
			Flag:  model.BinaryFlag | model.MultipleKeyFlag | model.HandleKeyFlag,
			Value: int64(2),
		}, nil},
		IndexColumns: [][]int{{1, 3}},
	}
	rowRedo := RowToRedo(row)
	require.Equal(t, 4, len(rowRedo.PreColumns))
	require.Equal(t, 4, len(rowRedo.Columns))

	redoLog := &model.RedoLog{
		Row:  rowRedo,
		Type: model.RedoLogTypeRow,
	}
	data, err := redoLog.MarshalMsg(nil)
	require.Nil(t, err)
	redoLog2 := &model.RedoLog{}
	_, err = redoLog2.UnmarshalMsg(data)
	require.Nil(t, err)
	require.Equal(t, row, LogToRow(redoLog2.Row))
}

func TestDDLRedoConvert(t *testing.T) {
	t.Parallel()
	ddl := &model.DDLEvent{
		StartTs:  1020,
		CommitTs: 1030,
		TableInfo: &model.SimpleTableInfo{
			Schema: "test",
			Table:  "t2",
		},
		Type:  timodel.ActionAddColumn,
		Query: "ALTER TABLE test.t1 ADD COLUMN a int",
	}
	redoDDL := DDLToRedo(ddl)

	redoLog := &model.RedoLog{
		DDL:  redoDDL,
		Type: model.RedoLogTypeDDL,
	}
	data, err := redoLog.MarshalMsg(nil)
	require.Nil(t, err)
	redoLog2 := &model.RedoLog{}
	_, err = redoLog2.UnmarshalMsg(data)
	require.Nil(t, err)
	require.Equal(t, ddl, LogToDDL(redoLog2.DDL))
}
