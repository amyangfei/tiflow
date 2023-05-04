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

package encoding

import (
	"testing"

	"github.com/pingcap/tiflow/cdc/model"
)

// bench cmd:
// go test -run='^$' -benchmem -bench '^(BenchmarkEncodKey)$' \
// github.com/pingcap/tiflow/cdc/processor/sourcemanager/engine/pebble/encoding
func BenchmarkEncodKey(b *testing.B) {
	b.Run("encode-key-bench", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = EncodeKey(1, 2, model.NewPolymorphicEvent(&model.RawKVEntry{
				OpType:  model.OpTypePut,
				Key:     []byte{3},
				StartTs: 4,
				CRTs:    5,
			}))
		}
	})
}
