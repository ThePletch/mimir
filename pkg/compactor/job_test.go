// SPDX-License-Identifier: AGPL-3.0-only

package compactor

import (
	"context"
	"path"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thanos-io/objstore"

	"github.com/grafana/mimir/pkg/storage/bucket"
	"github.com/grafana/mimir/pkg/storage/tsdb/block"
)

func TestJob_MinCompactionLevel(t *testing.T) {
	job := newJob("user-1", "group-1", labels.EmptyLabels(), 0, true, 2, "shard-1")
	require.NoError(t, job.AppendMeta(&block.Meta{BlockMeta: tsdb.BlockMeta{ULID: ulid.MustNew(1, nil), Compaction: tsdb.BlockMetaCompaction{Level: 2}}}))
	assert.Equal(t, 2, job.MinCompactionLevel())

	require.NoError(t, job.AppendMeta(&block.Meta{BlockMeta: tsdb.BlockMeta{ULID: ulid.MustNew(2, nil), Compaction: tsdb.BlockMetaCompaction{Level: 3}}}))
	assert.Equal(t, 2, job.MinCompactionLevel())

	require.NoError(t, job.AppendMeta(&block.Meta{BlockMeta: tsdb.BlockMeta{ULID: ulid.MustNew(3, nil), Compaction: tsdb.BlockMetaCompaction{Level: 1}}}))
	assert.Equal(t, 1, job.MinCompactionLevel())
}

func TestJobWaitPeriodElapsed(t *testing.T) {
	type jobBlock struct {
		meta     *block.Meta
		attrs    objstore.ObjectAttributes
		attrsErr error
	}

	// Blocks with compaction level 1.
	meta1 := &block.Meta{BlockMeta: tsdb.BlockMeta{ULID: ulid.MustNew(1, nil), Compaction: tsdb.BlockMetaCompaction{Level: 1}}}
	meta2 := &block.Meta{BlockMeta: tsdb.BlockMeta{ULID: ulid.MustNew(2, nil), Compaction: tsdb.BlockMetaCompaction{Level: 1}}}

	// Blocks with compaction level 2.
	meta3 := &block.Meta{BlockMeta: tsdb.BlockMeta{ULID: ulid.MustNew(3, nil), Compaction: tsdb.BlockMetaCompaction{Level: 2}}}
	meta4 := &block.Meta{BlockMeta: tsdb.BlockMeta{ULID: ulid.MustNew(4, nil), Compaction: tsdb.BlockMetaCompaction{Level: 2}}}

	// OOO blocks
	meta5 := &block.Meta{BlockMeta: tsdb.BlockMeta{ULID: ulid.MustNew(5, nil), Compaction: tsdb.BlockMetaCompaction{Level: 1}}}
	meta6 := &block.Meta{BlockMeta: tsdb.BlockMeta{ULID: ulid.MustNew(6, nil), Compaction: tsdb.BlockMetaCompaction{Level: 1}}}
	meta5.Compaction.SetOutOfOrder()
	meta6.Compaction.SetOutOfOrder()

	tests := map[string]struct {
		waitPeriod      time.Duration
		jobBlocks       []jobBlock
		expectedElapsed bool
		expectedMeta    *block.Meta
		expectedErr     string
	}{
		"wait period disabled": {
			waitPeriod: 0,
			jobBlocks: []jobBlock{
				{meta: meta1, attrs: objstore.ObjectAttributes{LastModified: time.Now().Add(-20 * time.Minute)}},
				{meta: meta2, attrs: objstore.ObjectAttributes{LastModified: time.Now().Add(-5 * time.Minute)}},
			},
			expectedElapsed: true,
			expectedMeta:    nil,
		},
		"blocks uploaded since more than the wait period": {
			waitPeriod: 10 * time.Minute,
			jobBlocks: []jobBlock{
				{meta: meta1, attrs: objstore.ObjectAttributes{LastModified: time.Now().Add(-20 * time.Minute)}},
				{meta: meta2, attrs: objstore.ObjectAttributes{LastModified: time.Now().Add(-25 * time.Minute)}},
			},
			expectedElapsed: true,
			expectedMeta:    nil,
		},
		"blocks uploaded since less than the wait period": {
			waitPeriod: 10 * time.Minute,
			jobBlocks: []jobBlock{
				{meta: meta1, attrs: objstore.ObjectAttributes{LastModified: time.Now().Add(-20 * time.Minute)}},
				{meta: meta2, attrs: objstore.ObjectAttributes{LastModified: time.Now().Add(-5 * time.Minute)}},
			},
			expectedElapsed: false,
			expectedMeta:    meta2,
		},
		"blocks uploaded since less than the wait period but their compaction level is > 1": {
			waitPeriod: 10 * time.Minute,
			jobBlocks: []jobBlock{
				{meta: meta3, attrs: objstore.ObjectAttributes{LastModified: time.Now().Add(-4 * time.Minute)}},
				{meta: meta4, attrs: objstore.ObjectAttributes{LastModified: time.Now().Add(-5 * time.Minute)}},
			},
			expectedElapsed: true,
			expectedMeta:    nil,
		},
		"out of order block": {
			waitPeriod: 10 * time.Minute,
			jobBlocks: []jobBlock{
				{meta: meta5, attrs: objstore.ObjectAttributes{LastModified: time.Now().Add(-20 * time.Minute)}},
				{meta: meta6, attrs: objstore.ObjectAttributes{LastModified: time.Now().Add(-5 * time.Minute)}},
			},
			expectedElapsed: true,
			expectedMeta:    nil,
		},
		"an error occurred while checking the blocks upload timestamp": {
			waitPeriod: 10 * time.Minute,
			jobBlocks: []jobBlock{
				// This block has been uploaded since more than the wait period.
				{meta: meta1, attrs: objstore.ObjectAttributes{LastModified: time.Now().Add(-20 * time.Minute)}},

				// This block has been uploaded since less than the wait period, but we failed getting its attributes.
				{meta: meta2, attrs: objstore.ObjectAttributes{LastModified: time.Now().Add(-5 * time.Minute)}, attrsErr: errors.New("mocked error")},
			},
			expectedErr:  "mocked error",
			expectedMeta: meta2,
		},
	}

	for testName, testData := range tests {
		t.Run(testName, func(t *testing.T) {
			job := newJob("user-1", "group-1", labels.EmptyLabels(), 0, true, 2, "shard-1")
			for _, b := range testData.jobBlocks {
				require.NoError(t, job.AppendMeta(b.meta))
			}

			userBucket := &bucket.ClientMock{}
			for _, b := range testData.jobBlocks {
				userBucket.MockAttributes(path.Join(b.meta.ULID.String(), block.MetaFilename), b.attrs, b.attrsErr)
			}

			elapsed, meta, err := jobWaitPeriodElapsed(context.Background(), job, testData.waitPeriod, userBucket)
			if testData.expectedErr != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, testData.expectedErr)
				assert.False(t, elapsed)
				assert.Equal(t, testData.expectedMeta, meta)
			} else {
				require.NoError(t, err)
				assert.Equal(t, testData.expectedElapsed, elapsed)
				assert.Equal(t, testData.expectedMeta, meta)
			}
		})
	}
}
