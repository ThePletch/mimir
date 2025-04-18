// SPDX-License-Identifier: AGPL-3.0-only
// Provenance-includes-location: https://github.com/cortexproject/cortex/blob/master/pkg/ingester/metrics_test.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: The Cortex Authors.

package ingester

import (
	"bytes"
	"testing"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
)

func TestTSDBMetrics(t *testing.T) {
	mainReg := prometheus.NewPedanticRegistry()

	tsdbMetrics := newTSDBMetrics(mainReg, log.NewNopLogger())

	tsdbMetrics.setRegistryForUser("user1", populateTSDBMetrics(12345))
	tsdbMetrics.setRegistryForUser("user2", populateTSDBMetrics(85787))
	tsdbMetrics.setRegistryForUser("user3", populateTSDBMetrics(999))

	err := testutil.GatherAndCompare(mainReg, bytes.NewBufferString(`
			# HELP cortex_ingester_tsdb_compactions_total Total number of TSDB compactions that were executed.
			# TYPE cortex_ingester_tsdb_compactions_total counter
			cortex_ingester_tsdb_compactions_total 693917

			# HELP cortex_ingester_tsdb_compaction_duration_seconds Duration of TSDB compaction runs.
			# TYPE cortex_ingester_tsdb_compaction_duration_seconds histogram
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="1"} 0
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="2"} 0
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="4"} 0
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="8"} 0
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="16"} 3
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="32"} 3
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="64"} 3
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="128"} 3
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="256"} 3
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="512"} 3
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="+Inf"} 3
			cortex_ingester_tsdb_compaction_duration_seconds_sum 27
			cortex_ingester_tsdb_compaction_duration_seconds_count 3

			# HELP cortex_ingester_tsdb_wal_fsync_duration_seconds Duration of TSDB WAL fsync.
			# TYPE cortex_ingester_tsdb_wal_fsync_duration_seconds summary
			cortex_ingester_tsdb_wal_fsync_duration_seconds{quantile="0.5"} 30
			cortex_ingester_tsdb_wal_fsync_duration_seconds{quantile="0.9"} 30
			cortex_ingester_tsdb_wal_fsync_duration_seconds{quantile="0.99"} 30
			cortex_ingester_tsdb_wal_fsync_duration_seconds_sum 30
			cortex_ingester_tsdb_wal_fsync_duration_seconds_count 3

			# HELP cortex_ingester_tsdb_wal_page_flushes_total Total number of TSDB WAL page flushes.
			# TYPE cortex_ingester_tsdb_wal_page_flushes_total counter
			cortex_ingester_tsdb_wal_page_flushes_total 1090441

			# HELP cortex_ingester_tsdb_wal_completed_pages_total Total number of TSDB WAL completed pages.
			# TYPE cortex_ingester_tsdb_wal_completed_pages_total counter
			cortex_ingester_tsdb_wal_completed_pages_total 1189572

			# HELP cortex_ingester_tsdb_wal_truncations_failed_total Total number of TSDB WAL truncations that failed.
			# TYPE cortex_ingester_tsdb_wal_truncations_failed_total counter
			cortex_ingester_tsdb_wal_truncations_failed_total 1288703

			# HELP cortex_ingester_tsdb_wal_truncations_total Total number of TSDB  WAL truncations attempted.
			# TYPE cortex_ingester_tsdb_wal_truncations_total counter
			cortex_ingester_tsdb_wal_truncations_total 1387834

			# HELP cortex_ingester_tsdb_wal_corruptions_total Total number of TSDB WAL corruptions.
			# TYPE cortex_ingester_tsdb_wal_corruptions_total counter
			cortex_ingester_tsdb_wal_corruptions_total 2.676537e+06

			# HELP cortex_ingester_tsdb_wal_writes_failed_total Total number of TSDB WAL writes that failed.
			# TYPE cortex_ingester_tsdb_wal_writes_failed_total counter
			cortex_ingester_tsdb_wal_writes_failed_total 1486965

			# HELP cortex_ingester_tsdb_head_truncations_failed_total Total number of TSDB head truncations that failed.
			# TYPE cortex_ingester_tsdb_head_truncations_failed_total counter
			cortex_ingester_tsdb_head_truncations_failed_total 2.775668e+06

			# HELP cortex_ingester_tsdb_head_truncations_total Total number of TSDB head truncations attempted.
			# TYPE cortex_ingester_tsdb_head_truncations_total counter
			cortex_ingester_tsdb_head_truncations_total 2.874799e+06

			# HELP cortex_ingester_tsdb_head_gc_duration_seconds Runtime of garbage collection in the TSDB head.
			# TYPE cortex_ingester_tsdb_head_gc_duration_seconds summary
			cortex_ingester_tsdb_head_gc_duration_seconds_sum 9
			cortex_ingester_tsdb_head_gc_duration_seconds_count 3

			# HELP cortex_ingester_tsdb_checkpoint_deletions_failed_total Total number of TSDB checkpoint deletions that failed.
			# TYPE cortex_ingester_tsdb_checkpoint_deletions_failed_total counter
			cortex_ingester_tsdb_checkpoint_deletions_failed_total 1586096

			# HELP cortex_ingester_tsdb_checkpoint_deletions_total Total number of TSDB checkpoint deletions attempted.
			# TYPE cortex_ingester_tsdb_checkpoint_deletions_total counter
			cortex_ingester_tsdb_checkpoint_deletions_total 1685227

			# HELP cortex_ingester_tsdb_checkpoint_creations_failed_total Total number of TSDB checkpoint creations that failed.
			# TYPE cortex_ingester_tsdb_checkpoint_creations_failed_total counter
			cortex_ingester_tsdb_checkpoint_creations_failed_total 1784358

			# HELP cortex_ingester_tsdb_checkpoint_creations_total Total number of TSDB checkpoint creations attempted.
			# TYPE cortex_ingester_tsdb_checkpoint_creations_total counter
			cortex_ingester_tsdb_checkpoint_creations_total 1883489

			# HELP cortex_ingester_memory_series The current number of series in memory.
			# TYPE cortex_ingester_memory_series gauge
			cortex_ingester_memory_series 396524

			# HELP cortex_ingester_memory_series_created_total The total number of series that were created per user.
			# TYPE cortex_ingester_memory_series_created_total counter
			# 5 * (12345, 85787 and 999 respectively)
			cortex_ingester_memory_series_created_total{user="user1"} 61725
			cortex_ingester_memory_series_created_total{user="user2"} 428935
			cortex_ingester_memory_series_created_total{user="user3"} 4995

			# HELP cortex_ingester_memory_series_removed_total The total number of series that were removed per user.
			# TYPE cortex_ingester_memory_series_removed_total counter
			# 6 * (12345, 85787 and 999 respectively)
			cortex_ingester_memory_series_removed_total{user="user1"} 74070
			cortex_ingester_memory_series_removed_total{user="user2"} 514722
			cortex_ingester_memory_series_removed_total{user="user3"} 5994

			# HELP cortex_ingester_tsdb_head_active_appenders Number of currently active TSDB appender transactions.
			# TYPE cortex_ingester_tsdb_head_active_appenders gauge
			cortex_ingester_tsdb_head_active_appenders 1982620

			# HELP cortex_ingester_tsdb_head_series_not_found_total Total number of TSDB requests for series that were not found.
			# TYPE cortex_ingester_tsdb_head_series_not_found_total counter
			cortex_ingester_tsdb_head_series_not_found_total 2081751

			# HELP cortex_ingester_tsdb_head_chunks Total number of chunks in the TSDB head block.
			# TYPE cortex_ingester_tsdb_head_chunks gauge
			cortex_ingester_tsdb_head_chunks 2180882

			# HELP cortex_ingester_tsdb_head_chunks_created_total Total number of series created in the TSDB head.
			# TYPE cortex_ingester_tsdb_head_chunks_created_total counter
			cortex_ingester_tsdb_head_chunks_created_total{user="user1"} 283935
			cortex_ingester_tsdb_head_chunks_created_total{user="user2"} 1973101
			cortex_ingester_tsdb_head_chunks_created_total{user="user3"} 22977

			# HELP cortex_ingester_tsdb_head_chunks_removed_total Total number of series removed in the TSDB head.
			# TYPE cortex_ingester_tsdb_head_chunks_removed_total counter
			cortex_ingester_tsdb_head_chunks_removed_total{user="user1"} 296280
			cortex_ingester_tsdb_head_chunks_removed_total{user="user2"} 2058888
			cortex_ingester_tsdb_head_chunks_removed_total{user="user3"} 23976

			# HELP cortex_ingester_tsdb_wal_truncate_duration_seconds Duration of TSDB WAL truncation.
			# TYPE cortex_ingester_tsdb_wal_truncate_duration_seconds summary
			cortex_ingester_tsdb_wal_truncate_duration_seconds_sum 75
			cortex_ingester_tsdb_wal_truncate_duration_seconds_count 3

			# HELP cortex_ingester_tsdb_mmap_chunk_corruptions_total Total number of memory-mapped TSDB chunk corruptions.
			# TYPE cortex_ingester_tsdb_mmap_chunk_corruptions_total counter
			cortex_ingester_tsdb_mmap_chunk_corruptions_total 2577406

			# HELP cortex_ingester_tsdb_mmap_chunk_write_queue_operations_total Total number of memory-mapped TSDB chunk operations.
			# TYPE cortex_ingester_tsdb_mmap_chunk_write_queue_operations_total counter
			cortex_ingester_tsdb_mmap_chunk_write_queue_operations_total{operation="add"} 150
			cortex_ingester_tsdb_mmap_chunk_write_queue_operations_total{operation="complete"} 120
			cortex_ingester_tsdb_mmap_chunk_write_queue_operations_total{operation="get"} 30

			# HELP cortex_ingester_tsdb_mmap_chunks_total Total number of chunks that were memory-mapped.
			# TYPE cortex_ingester_tsdb_mmap_chunks_total counter
			cortex_ingester_tsdb_mmap_chunks_total 2973930

			# HELP cortex_ingester_tsdb_blocks_loaded Number of currently loaded data blocks
			# TYPE cortex_ingester_tsdb_blocks_loaded gauge
			cortex_ingester_tsdb_blocks_loaded 15

			# HELP cortex_ingester_tsdb_reloads_total Number of times the database reloaded block data from disk.
			# TYPE cortex_ingester_tsdb_reloads_total counter
			cortex_ingester_tsdb_reloads_total 30

			# HELP cortex_ingester_tsdb_sample_out_of_order_delta_seconds Delta in seconds by which a sample is considered out-of-order.
			# TYPE cortex_ingester_tsdb_sample_out_of_order_delta_seconds histogram
			# observations        buckets
			#                     600
			# 7*999            -> 86400
			# 7*12245, 7*85787 -> inf
			cortex_ingester_tsdb_sample_out_of_order_delta_seconds_bucket{le="600"} 0
			cortex_ingester_tsdb_sample_out_of_order_delta_seconds_bucket{le="86400"} 1
			cortex_ingester_tsdb_sample_out_of_order_delta_seconds_bucket{le="+Inf"} 3
			cortex_ingester_tsdb_sample_out_of_order_delta_seconds_sum 693917
			cortex_ingester_tsdb_sample_out_of_order_delta_seconds_count 3

			# HELP cortex_ingester_tsdb_reloads_failures_total Number of times the database failed to reloadBlocks block data from disk.
			# TYPE cortex_ingester_tsdb_reloads_failures_total counter
			cortex_ingester_tsdb_reloads_failures_total 21

			# HELP cortex_ingester_tsdb_symbol_table_size_bytes Size of symbol table in memory for loaded blocks
			# TYPE cortex_ingester_tsdb_symbol_table_size_bytes gauge
			cortex_ingester_tsdb_symbol_table_size_bytes{user="user1"} 12641280
			cortex_ingester_tsdb_symbol_table_size_bytes{user="user2"} 87845888
			cortex_ingester_tsdb_symbol_table_size_bytes{user="user3"} 1022976

			# HELP cortex_ingester_tsdb_storage_blocks_bytes The number of bytes that are currently used for local storage by all blocks.
			# TYPE cortex_ingester_tsdb_storage_blocks_bytes gauge
			cortex_ingester_tsdb_storage_blocks_bytes{user="user1"} 50565120
			cortex_ingester_tsdb_storage_blocks_bytes{user="user2"} 351383552
			cortex_ingester_tsdb_storage_blocks_bytes{user="user3"} 4091904

			# HELP cortex_ingester_tsdb_time_retentions_total The number of times that blocks were deleted because the maximum time limit was exceeded.
			# TYPE cortex_ingester_tsdb_time_retentions_total counter
			cortex_ingester_tsdb_time_retentions_total 33

			# HELP cortex_ingester_tsdb_exemplar_last_exemplars_timestamp_seconds The timestamp of the oldest exemplar stored in circular storage. Useful to check for what time range the current exemplar buffer limit allows. This usually means the last timestamp for all exemplars for a typical setup. This is not true though if one of the series timestamp is in future compared to rest series.
			# TYPE cortex_ingester_tsdb_exemplar_last_exemplars_timestamp_seconds gauge
			cortex_ingester_tsdb_exemplar_last_exemplars_timestamp_seconds{user="user1"} 1234
			cortex_ingester_tsdb_exemplar_last_exemplars_timestamp_seconds{user="user2"} 1234
			cortex_ingester_tsdb_exemplar_last_exemplars_timestamp_seconds{user="user3"} 1234

			# HELP cortex_ingester_tsdb_exemplar_out_of_order_exemplars_total Total number of out-of-order exemplar ingestion failed attempts.
			# TYPE cortex_ingester_tsdb_exemplar_out_of_order_exemplars_total counter
			cortex_ingester_tsdb_exemplar_out_of_order_exemplars_total 9

			# HELP cortex_ingester_tsdb_exemplar_series_with_exemplars_in_storage Number of TSDB series with exemplars currently in storage.
			# TYPE cortex_ingester_tsdb_exemplar_series_with_exemplars_in_storage gauge
			cortex_ingester_tsdb_exemplar_series_with_exemplars_in_storage{user="user1"} 1
			cortex_ingester_tsdb_exemplar_series_with_exemplars_in_storage{user="user2"} 1
			cortex_ingester_tsdb_exemplar_series_with_exemplars_in_storage{user="user3"} 1

			# HELP cortex_ingester_tsdb_exemplar_exemplars_appended_total Total number of TSDB exemplars appended.
			# TYPE cortex_ingester_tsdb_exemplar_exemplars_appended_total counter
			cortex_ingester_tsdb_exemplar_exemplars_appended_total{user="user1"} 100
			cortex_ingester_tsdb_exemplar_exemplars_appended_total{user="user2"} 100
			cortex_ingester_tsdb_exemplar_exemplars_appended_total{user="user3"} 100

			# HELP cortex_ingester_tsdb_out_of_order_samples_appended_total Total number of out-of-order samples appended.
			# TYPE cortex_ingester_tsdb_out_of_order_samples_appended_total counter
			cortex_ingester_tsdb_out_of_order_samples_appended_total{user="user1"} 3
			cortex_ingester_tsdb_out_of_order_samples_appended_total{user="user2"} 3
			cortex_ingester_tsdb_out_of_order_samples_appended_total{user="user3"} 3

			# HELP cortex_ingester_tsdb_exemplar_exemplars_in_storage Number of TSDB exemplars currently in storage.
			# TYPE cortex_ingester_tsdb_exemplar_exemplars_in_storage gauge
			cortex_ingester_tsdb_exemplar_exemplars_in_storage 30

			# HELP cortex_ingester_tsdb_head_postings_for_matchers_cache_hits_total Total number of postings lists returned from the PostingsForMatchers cache.
			# TYPE cortex_ingester_tsdb_head_postings_for_matchers_cache_hits_total counter
			cortex_ingester_tsdb_head_postings_for_matchers_cache_hits_total 0

			# HELP cortex_ingester_tsdb_head_postings_for_matchers_cache_misses_total Total number of requests to the PostingsForMatchers cache for which there is no valid cached entry. The subsequent result is cached.
			# TYPE cortex_ingester_tsdb_head_postings_for_matchers_cache_misses_total counter
			cortex_ingester_tsdb_head_postings_for_matchers_cache_misses_total 0

			# HELP cortex_ingester_tsdb_head_postings_for_matchers_cache_requests_total Total number of requests to the PostingsForMatchers cache.
			# TYPE cortex_ingester_tsdb_head_postings_for_matchers_cache_requests_total counter
			cortex_ingester_tsdb_head_postings_for_matchers_cache_requests_total 0

			# HELP cortex_ingester_tsdb_head_postings_for_matchers_cache_skips_total Total number of requests to the PostingsForMatchers cache that have been skipped the cache. The subsequent result is not cached.
			# TYPE cortex_ingester_tsdb_head_postings_for_matchers_cache_skips_total counter
			cortex_ingester_tsdb_head_postings_for_matchers_cache_skips_total{reason="canceled-cached-entry"} 0
			cortex_ingester_tsdb_head_postings_for_matchers_cache_skips_total{reason="ineligible"} 0
			cortex_ingester_tsdb_head_postings_for_matchers_cache_skips_total{reason="stale-cached-entry"} 0

			# HELP cortex_ingester_tsdb_head_postings_for_matchers_cache_evictions_total Total number of evictions from the PostingsForMatchers cache.
			# TYPE cortex_ingester_tsdb_head_postings_for_matchers_cache_evictions_total counter
			cortex_ingester_tsdb_head_postings_for_matchers_cache_evictions_total{reason="max-bytes-reached"} 0
			cortex_ingester_tsdb_head_postings_for_matchers_cache_evictions_total{reason="max-items-reached"} 0
			cortex_ingester_tsdb_head_postings_for_matchers_cache_evictions_total{reason="ttl-expired"} 0
			cortex_ingester_tsdb_head_postings_for_matchers_cache_evictions_total{reason="unknown"} 0

			# HELP cortex_ingester_tsdb_block_postings_for_matchers_cache_hits_total Total number of postings lists returned from the PostingsForMatchers cache.
			# TYPE cortex_ingester_tsdb_block_postings_for_matchers_cache_hits_total counter
			cortex_ingester_tsdb_block_postings_for_matchers_cache_hits_total 0

			# HELP cortex_ingester_tsdb_block_postings_for_matchers_cache_misses_total Total number of requests to the PostingsForMatchers cache for which there is no valid cached entry. The subsequent result is cached.
			# TYPE cortex_ingester_tsdb_block_postings_for_matchers_cache_misses_total counter
			cortex_ingester_tsdb_block_postings_for_matchers_cache_misses_total 0

			# HELP cortex_ingester_tsdb_block_postings_for_matchers_cache_requests_total Total number of requests to the PostingsForMatchers cache.
			# TYPE cortex_ingester_tsdb_block_postings_for_matchers_cache_requests_total counter
			cortex_ingester_tsdb_block_postings_for_matchers_cache_requests_total 0

			# HELP cortex_ingester_tsdb_block_postings_for_matchers_cache_skips_total Total number of requests to the PostingsForMatchers cache that have been skipped the cache. The subsequent result is not cached.
			# TYPE cortex_ingester_tsdb_block_postings_for_matchers_cache_skips_total counter
			cortex_ingester_tsdb_block_postings_for_matchers_cache_skips_total{reason="canceled-cached-entry"} 0
			cortex_ingester_tsdb_block_postings_for_matchers_cache_skips_total{reason="ineligible"} 0
			cortex_ingester_tsdb_block_postings_for_matchers_cache_skips_total{reason="stale-cached-entry"} 0

			# HELP cortex_ingester_tsdb_block_postings_for_matchers_cache_evictions_total Total number of evictions from the PostingsForMatchers cache.
			# TYPE cortex_ingester_tsdb_block_postings_for_matchers_cache_evictions_total counter
			cortex_ingester_tsdb_block_postings_for_matchers_cache_evictions_total{reason="max-bytes-reached"} 0
			cortex_ingester_tsdb_block_postings_for_matchers_cache_evictions_total{reason="max-items-reached"} 0
			cortex_ingester_tsdb_block_postings_for_matchers_cache_evictions_total{reason="ttl-expired"} 0
			cortex_ingester_tsdb_block_postings_for_matchers_cache_evictions_total{reason="unknown"} 0

			# HELP cortex_ingester_tsdb_wal_replay_unknown_refs_total Total number of unknown series references encountered during WAL replay.
			# TYPE cortex_ingester_tsdb_wal_replay_unknown_refs_total counter
			cortex_ingester_tsdb_wal_replay_unknown_refs_total{type="series", user="user1"} 12345
			cortex_ingester_tsdb_wal_replay_unknown_refs_total{type="samples", user="user1"} 24690
			cortex_ingester_tsdb_wal_replay_unknown_refs_total{type="series", user="user2"}  85787
			cortex_ingester_tsdb_wal_replay_unknown_refs_total{type="samples", user="user2"} 171574
			cortex_ingester_tsdb_wal_replay_unknown_refs_total{type="series", user="user3"} 999
			cortex_ingester_tsdb_wal_replay_unknown_refs_total{type="samples", user="user3"} 1998

			# HELP cortex_ingester_tsdb_wbl_replay_unknown_refs_total Total number of unknown series references encountered during WBL replay.
			# TYPE cortex_ingester_tsdb_wbl_replay_unknown_refs_total counter
			cortex_ingester_tsdb_wbl_replay_unknown_refs_total{type="exemplars", user="user1"} 12345
			cortex_ingester_tsdb_wbl_replay_unknown_refs_total{type="exemplars", user="user2"} 85787
			cortex_ingester_tsdb_wbl_replay_unknown_refs_total{type="exemplars", user="user3"} 999
	`))
	require.NoError(t, err)
}

func TestTSDBMetricsWithRemoval(t *testing.T) {
	mainReg := prometheus.NewPedanticRegistry()

	tsdbMetrics := newTSDBMetrics(mainReg, log.NewNopLogger())

	tsdbMetrics.setRegistryForUser("user1", populateTSDBMetrics(12345))
	tsdbMetrics.setRegistryForUser("user2", populateTSDBMetrics(85787))
	tsdbMetrics.setRegistryForUser("user3", populateTSDBMetrics(999))
	tsdbMetrics.removeRegistryForUser("user3")

	err := testutil.GatherAndCompare(mainReg, bytes.NewBufferString(`
			# HELP cortex_ingester_tsdb_compactions_total Total number of TSDB compactions that were executed.
			# TYPE cortex_ingester_tsdb_compactions_total counter
			cortex_ingester_tsdb_compactions_total 693917

			# HELP cortex_ingester_tsdb_compaction_duration_seconds Duration of TSDB compaction runs.
			# TYPE cortex_ingester_tsdb_compaction_duration_seconds histogram
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="1"} 0
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="2"} 0
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="4"} 0
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="8"} 0
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="16"} 3
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="32"} 3
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="64"} 3
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="128"} 3
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="256"} 3
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="512"} 3
			cortex_ingester_tsdb_compaction_duration_seconds_bucket{le="+Inf"} 3
			cortex_ingester_tsdb_compaction_duration_seconds_sum 27
			cortex_ingester_tsdb_compaction_duration_seconds_count 3

			# HELP cortex_ingester_tsdb_wal_fsync_duration_seconds Duration of TSDB WAL fsync.
			# TYPE cortex_ingester_tsdb_wal_fsync_duration_seconds summary
			cortex_ingester_tsdb_wal_fsync_duration_seconds{quantile="0.5"} 30
			cortex_ingester_tsdb_wal_fsync_duration_seconds{quantile="0.9"} 30
			cortex_ingester_tsdb_wal_fsync_duration_seconds{quantile="0.99"} 30
			cortex_ingester_tsdb_wal_fsync_duration_seconds_sum 30
			cortex_ingester_tsdb_wal_fsync_duration_seconds_count 3

			# HELP cortex_ingester_tsdb_wal_page_flushes_total Total number of TSDB WAL page flushes.
			# TYPE cortex_ingester_tsdb_wal_page_flushes_total counter
			cortex_ingester_tsdb_wal_page_flushes_total 1090441

			# HELP cortex_ingester_tsdb_wal_completed_pages_total Total number of TSDB WAL completed pages.
			# TYPE cortex_ingester_tsdb_wal_completed_pages_total counter
			cortex_ingester_tsdb_wal_completed_pages_total 1189572

			# HELP cortex_ingester_tsdb_wal_truncations_failed_total Total number of TSDB WAL truncations that failed.
			# TYPE cortex_ingester_tsdb_wal_truncations_failed_total counter
			cortex_ingester_tsdb_wal_truncations_failed_total 1288703

			# HELP cortex_ingester_tsdb_wal_truncations_total Total number of TSDB  WAL truncations attempted.
			# TYPE cortex_ingester_tsdb_wal_truncations_total counter
			cortex_ingester_tsdb_wal_truncations_total 1387834

			# HELP cortex_ingester_tsdb_wal_corruptions_total Total number of TSDB WAL corruptions.
			# TYPE cortex_ingester_tsdb_wal_corruptions_total counter
			cortex_ingester_tsdb_wal_corruptions_total 2.676537e+06

			# HELP cortex_ingester_tsdb_wal_writes_failed_total Total number of TSDB WAL writes that failed.
			# TYPE cortex_ingester_tsdb_wal_writes_failed_total counter
			cortex_ingester_tsdb_wal_writes_failed_total 1486965

			# HELP cortex_ingester_tsdb_head_truncations_failed_total Total number of TSDB head truncations that failed.
			# TYPE cortex_ingester_tsdb_head_truncations_failed_total counter
			cortex_ingester_tsdb_head_truncations_failed_total 2.775668e+06

			# HELP cortex_ingester_tsdb_head_truncations_total Total number of TSDB head truncations attempted.
			# TYPE cortex_ingester_tsdb_head_truncations_total counter
			cortex_ingester_tsdb_head_truncations_total 2.874799e+06

			# HELP cortex_ingester_tsdb_head_gc_duration_seconds Runtime of garbage collection in the TSDB head.
			# TYPE cortex_ingester_tsdb_head_gc_duration_seconds summary
			cortex_ingester_tsdb_head_gc_duration_seconds_sum 9
			cortex_ingester_tsdb_head_gc_duration_seconds_count 3

			# HELP cortex_ingester_tsdb_checkpoint_deletions_failed_total Total number of TSDB checkpoint deletions that failed.
			# TYPE cortex_ingester_tsdb_checkpoint_deletions_failed_total counter
			cortex_ingester_tsdb_checkpoint_deletions_failed_total 1586096

			# HELP cortex_ingester_tsdb_checkpoint_deletions_total Total number of TSDB checkpoint deletions attempted.
			# TYPE cortex_ingester_tsdb_checkpoint_deletions_total counter
			cortex_ingester_tsdb_checkpoint_deletions_total 1685227

			# HELP cortex_ingester_tsdb_checkpoint_creations_failed_total Total number of TSDB checkpoint creations that failed.
			# TYPE cortex_ingester_tsdb_checkpoint_creations_failed_total counter
			cortex_ingester_tsdb_checkpoint_creations_failed_total 1784358

			# HELP cortex_ingester_tsdb_checkpoint_creations_total Total number of TSDB checkpoint creations attempted.
			# TYPE cortex_ingester_tsdb_checkpoint_creations_total counter
			cortex_ingester_tsdb_checkpoint_creations_total 1883489

			# HELP cortex_ingester_memory_series The current number of series in memory.
			# TYPE cortex_ingester_memory_series gauge
			cortex_ingester_memory_series 392528

			# HELP cortex_ingester_memory_series_created_total The total number of series that were created per user.
			# TYPE cortex_ingester_memory_series_created_total counter
			# 5 * (12345, 85787 and 999 respectively)
			cortex_ingester_memory_series_created_total{user="user1"} 61725
			cortex_ingester_memory_series_created_total{user="user2"} 428935

			# HELP cortex_ingester_memory_series_removed_total The total number of series that were removed per user.
			# TYPE cortex_ingester_memory_series_removed_total counter
			# 6 * (12345, 85787 and 999 respectively)
			cortex_ingester_memory_series_removed_total{user="user1"} 74070
			cortex_ingester_memory_series_removed_total{user="user2"} 514722

			# HELP cortex_ingester_tsdb_head_active_appenders Number of currently active TSDB appender transactions.
			# TYPE cortex_ingester_tsdb_head_active_appenders gauge
			cortex_ingester_tsdb_head_active_appenders 1962640

			# HELP cortex_ingester_tsdb_head_series_not_found_total Total number of TSDB requests for series that were not found.
			# TYPE cortex_ingester_tsdb_head_series_not_found_total counter
			cortex_ingester_tsdb_head_series_not_found_total 2081751

			# HELP cortex_ingester_tsdb_head_chunks Total number of chunks in the TSDB head block.
			# TYPE cortex_ingester_tsdb_head_chunks gauge
			cortex_ingester_tsdb_head_chunks 2158904

			# HELP cortex_ingester_tsdb_head_chunks_created_total Total number of series created in the TSDB head.
			# TYPE cortex_ingester_tsdb_head_chunks_created_total counter
			cortex_ingester_tsdb_head_chunks_created_total{user="user1"} 283935
			cortex_ingester_tsdb_head_chunks_created_total{user="user2"} 1973101

			# HELP cortex_ingester_tsdb_head_chunks_removed_total Total number of series removed in the TSDB head.
			# TYPE cortex_ingester_tsdb_head_chunks_removed_total counter
			cortex_ingester_tsdb_head_chunks_removed_total{user="user1"} 296280
			cortex_ingester_tsdb_head_chunks_removed_total{user="user2"} 2058888

			# HELP cortex_ingester_tsdb_mmap_chunks_total Total number of chunks that were memory-mapped.
			# TYPE cortex_ingester_tsdb_mmap_chunks_total counter
			cortex_ingester_tsdb_mmap_chunks_total 2973930

			# HELP cortex_ingester_tsdb_wal_truncate_duration_seconds Duration of TSDB WAL truncation.
			# TYPE cortex_ingester_tsdb_wal_truncate_duration_seconds summary
			cortex_ingester_tsdb_wal_truncate_duration_seconds_sum 75
			cortex_ingester_tsdb_wal_truncate_duration_seconds_count 3

			# HELP cortex_ingester_tsdb_mmap_chunk_corruptions_total Total number of memory-mapped TSDB chunk corruptions.
			# TYPE cortex_ingester_tsdb_mmap_chunk_corruptions_total counter
			cortex_ingester_tsdb_mmap_chunk_corruptions_total 2577406

			# HELP cortex_ingester_tsdb_mmap_chunk_write_queue_operations_total Total number of memory-mapped TSDB chunk operations.
			# TYPE cortex_ingester_tsdb_mmap_chunk_write_queue_operations_total counter
			cortex_ingester_tsdb_mmap_chunk_write_queue_operations_total{operation="add"} 150
			cortex_ingester_tsdb_mmap_chunk_write_queue_operations_total{operation="complete"} 120
			cortex_ingester_tsdb_mmap_chunk_write_queue_operations_total{operation="get"} 30

			# HELP cortex_ingester_tsdb_blocks_loaded Number of currently loaded data blocks
			# TYPE cortex_ingester_tsdb_blocks_loaded gauge
			cortex_ingester_tsdb_blocks_loaded 10

			# HELP cortex_ingester_tsdb_reloads_total Number of times the database reloaded block data from disk.
			# TYPE cortex_ingester_tsdb_reloads_total counter
			cortex_ingester_tsdb_reloads_total 30

			# HELP cortex_ingester_tsdb_sample_out_of_order_delta_seconds Delta in seconds by which a sample is considered out-of-order.
			# TYPE cortex_ingester_tsdb_sample_out_of_order_delta_seconds histogram
			# observations        buckets
			#                     600
			# 7*999            -> 86400
			# 7*12245, 7*85787 -> inf
			cortex_ingester_tsdb_sample_out_of_order_delta_seconds_bucket{le="600"} 0
			cortex_ingester_tsdb_sample_out_of_order_delta_seconds_bucket{le="86400"} 1
			cortex_ingester_tsdb_sample_out_of_order_delta_seconds_bucket{le="+Inf"} 3
			cortex_ingester_tsdb_sample_out_of_order_delta_seconds_sum 693917
			cortex_ingester_tsdb_sample_out_of_order_delta_seconds_count 3

			# HELP cortex_ingester_tsdb_reloads_failures_total Number of times the database failed to reloadBlocks block data from disk.
			# TYPE cortex_ingester_tsdb_reloads_failures_total counter
			cortex_ingester_tsdb_reloads_failures_total 21

			# HELP cortex_ingester_tsdb_symbol_table_size_bytes Size of symbol table in memory for loaded blocks
			# TYPE cortex_ingester_tsdb_symbol_table_size_bytes gauge
			cortex_ingester_tsdb_symbol_table_size_bytes{user="user1"} 12641280
			cortex_ingester_tsdb_symbol_table_size_bytes{user="user2"} 87845888

			# HELP cortex_ingester_tsdb_storage_blocks_bytes The number of bytes that are currently used for local storage by all blocks.
			# TYPE cortex_ingester_tsdb_storage_blocks_bytes gauge
			cortex_ingester_tsdb_storage_blocks_bytes{user="user1"} 50565120
			cortex_ingester_tsdb_storage_blocks_bytes{user="user2"} 351383552

			# HELP cortex_ingester_tsdb_time_retentions_total The number of times that blocks were deleted because the maximum time limit was exceeded.
			# TYPE cortex_ingester_tsdb_time_retentions_total counter
			cortex_ingester_tsdb_time_retentions_total 33

			# HELP cortex_ingester_tsdb_exemplar_last_exemplars_timestamp_seconds The timestamp of the oldest exemplar stored in circular storage. Useful to check for what time range the current exemplar buffer limit allows. This usually means the last timestamp for all exemplars for a typical setup. This is not true though if one of the series timestamp is in future compared to rest series.
			# TYPE cortex_ingester_tsdb_exemplar_last_exemplars_timestamp_seconds gauge
			cortex_ingester_tsdb_exemplar_last_exemplars_timestamp_seconds{user="user1"} 1234
			cortex_ingester_tsdb_exemplar_last_exemplars_timestamp_seconds{user="user2"} 1234

			# HELP cortex_ingester_tsdb_exemplar_out_of_order_exemplars_total Total number of out-of-order exemplar ingestion failed attempts.
			# TYPE cortex_ingester_tsdb_exemplar_out_of_order_exemplars_total counter
			cortex_ingester_tsdb_exemplar_out_of_order_exemplars_total 9

			# HELP cortex_ingester_tsdb_exemplar_series_with_exemplars_in_storage Number of TSDB series with exemplars currently in storage.
			# TYPE cortex_ingester_tsdb_exemplar_series_with_exemplars_in_storage gauge
			cortex_ingester_tsdb_exemplar_series_with_exemplars_in_storage{user="user1"} 1
			cortex_ingester_tsdb_exemplar_series_with_exemplars_in_storage{user="user2"} 1

			# HELP cortex_ingester_tsdb_exemplar_exemplars_appended_total Total number of TSDB exemplars appended.
			# TYPE cortex_ingester_tsdb_exemplar_exemplars_appended_total counter
			cortex_ingester_tsdb_exemplar_exemplars_appended_total{user="user1"} 100
			cortex_ingester_tsdb_exemplar_exemplars_appended_total{user="user2"} 100

			# HELP cortex_ingester_tsdb_exemplar_exemplars_in_storage Number of TSDB exemplars currently in storage.
			# TYPE cortex_ingester_tsdb_exemplar_exemplars_in_storage gauge
			cortex_ingester_tsdb_exemplar_exemplars_in_storage 20

			# HELP cortex_ingester_tsdb_out_of_order_samples_appended_total Total number of out-of-order samples appended.
			# TYPE cortex_ingester_tsdb_out_of_order_samples_appended_total counter
			cortex_ingester_tsdb_out_of_order_samples_appended_total{user="user1"} 3
			cortex_ingester_tsdb_out_of_order_samples_appended_total{user="user2"} 3

			# HELP cortex_ingester_tsdb_head_postings_for_matchers_cache_hits_total Total number of postings lists returned from the PostingsForMatchers cache.
			# TYPE cortex_ingester_tsdb_head_postings_for_matchers_cache_hits_total counter
			cortex_ingester_tsdb_head_postings_for_matchers_cache_hits_total 0

			# HELP cortex_ingester_tsdb_head_postings_for_matchers_cache_misses_total Total number of requests to the PostingsForMatchers cache for which there is no valid cached entry. The subsequent result is cached.
			# TYPE cortex_ingester_tsdb_head_postings_for_matchers_cache_misses_total counter
			cortex_ingester_tsdb_head_postings_for_matchers_cache_misses_total 0

			# HELP cortex_ingester_tsdb_head_postings_for_matchers_cache_requests_total Total number of requests to the PostingsForMatchers cache.
			# TYPE cortex_ingester_tsdb_head_postings_for_matchers_cache_requests_total counter
			cortex_ingester_tsdb_head_postings_for_matchers_cache_requests_total 0

			# HELP cortex_ingester_tsdb_head_postings_for_matchers_cache_skips_total Total number of requests to the PostingsForMatchers cache that have been skipped the cache. The subsequent result is not cached.
			# TYPE cortex_ingester_tsdb_head_postings_for_matchers_cache_skips_total counter
			cortex_ingester_tsdb_head_postings_for_matchers_cache_skips_total{reason="canceled-cached-entry"} 0
			cortex_ingester_tsdb_head_postings_for_matchers_cache_skips_total{reason="ineligible"} 0
			cortex_ingester_tsdb_head_postings_for_matchers_cache_skips_total{reason="stale-cached-entry"} 0

			# HELP cortex_ingester_tsdb_head_postings_for_matchers_cache_evictions_total Total number of evictions from the PostingsForMatchers cache.
			# TYPE cortex_ingester_tsdb_head_postings_for_matchers_cache_evictions_total counter
			cortex_ingester_tsdb_head_postings_for_matchers_cache_evictions_total{reason="max-bytes-reached"} 0
			cortex_ingester_tsdb_head_postings_for_matchers_cache_evictions_total{reason="max-items-reached"} 0
			cortex_ingester_tsdb_head_postings_for_matchers_cache_evictions_total{reason="ttl-expired"} 0
			cortex_ingester_tsdb_head_postings_for_matchers_cache_evictions_total{reason="unknown"} 0

			# HELP cortex_ingester_tsdb_block_postings_for_matchers_cache_hits_total Total number of postings lists returned from the PostingsForMatchers cache.
			# TYPE cortex_ingester_tsdb_block_postings_for_matchers_cache_hits_total counter
			cortex_ingester_tsdb_block_postings_for_matchers_cache_hits_total 0

			# HELP cortex_ingester_tsdb_block_postings_for_matchers_cache_misses_total Total number of requests to the PostingsForMatchers cache for which there is no valid cached entry. The subsequent result is cached.
			# TYPE cortex_ingester_tsdb_block_postings_for_matchers_cache_misses_total counter
			cortex_ingester_tsdb_block_postings_for_matchers_cache_misses_total 0

			# HELP cortex_ingester_tsdb_block_postings_for_matchers_cache_requests_total Total number of requests to the PostingsForMatchers cache.
			# TYPE cortex_ingester_tsdb_block_postings_for_matchers_cache_requests_total counter
			cortex_ingester_tsdb_block_postings_for_matchers_cache_requests_total 0

			# HELP cortex_ingester_tsdb_block_postings_for_matchers_cache_skips_total Total number of requests to the PostingsForMatchers cache that have been skipped the cache. The subsequent result is not cached.
			# TYPE cortex_ingester_tsdb_block_postings_for_matchers_cache_skips_total counter
			cortex_ingester_tsdb_block_postings_for_matchers_cache_skips_total{reason="canceled-cached-entry"} 0
			cortex_ingester_tsdb_block_postings_for_matchers_cache_skips_total{reason="ineligible"} 0
			cortex_ingester_tsdb_block_postings_for_matchers_cache_skips_total{reason="stale-cached-entry"} 0

			# HELP cortex_ingester_tsdb_block_postings_for_matchers_cache_evictions_total Total number of evictions from the PostingsForMatchers cache.
			# TYPE cortex_ingester_tsdb_block_postings_for_matchers_cache_evictions_total counter
			cortex_ingester_tsdb_block_postings_for_matchers_cache_evictions_total{reason="max-bytes-reached"} 0
			cortex_ingester_tsdb_block_postings_for_matchers_cache_evictions_total{reason="max-items-reached"} 0
			cortex_ingester_tsdb_block_postings_for_matchers_cache_evictions_total{reason="ttl-expired"} 0
			cortex_ingester_tsdb_block_postings_for_matchers_cache_evictions_total{reason="unknown"} 0

			# HELP cortex_ingester_tsdb_wal_replay_unknown_refs_total Total number of unknown series references encountered during WAL replay.
			# TYPE cortex_ingester_tsdb_wal_replay_unknown_refs_total counter
			cortex_ingester_tsdb_wal_replay_unknown_refs_total{type="series", user="user1"} 12345
			cortex_ingester_tsdb_wal_replay_unknown_refs_total{type="samples", user="user1"} 24690
			cortex_ingester_tsdb_wal_replay_unknown_refs_total{type="series", user="user2"}  85787
			cortex_ingester_tsdb_wal_replay_unknown_refs_total{type="samples", user="user2"} 171574

			# HELP cortex_ingester_tsdb_wbl_replay_unknown_refs_total Total number of unknown series references encountered during WBL replay.
			# TYPE cortex_ingester_tsdb_wbl_replay_unknown_refs_total counter
			cortex_ingester_tsdb_wbl_replay_unknown_refs_total{type="exemplars", user="user1"} 12345
			cortex_ingester_tsdb_wbl_replay_unknown_refs_total{type="exemplars", user="user2"} 85787
	`))
	require.NoError(t, err)
}

func populateTSDBMetrics(base float64) *prometheus.Registry {
	r := prometheus.NewRegistry()

	// TSDB Head
	headSeries := promauto.With(r).NewGauge(prometheus.GaugeOpts{
		Name: "prometheus_tsdb_head_series",
	})
	headSeries.Add(4 * base)

	seriesCreated := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_head_series_created_total",
	})
	seriesCreated.Add(5 * base)

	seriesRemoved := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_head_series_removed_total",
	})
	seriesRemoved.Add(6 * base)

	ran := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_compactions_total",
		Help: "Total number of compactions that were executed for the partition.",
	})
	ran.Add(7 * base)

	duration := promauto.With(r).NewHistogram(prometheus.HistogramOpts{
		Name:    "prometheus_tsdb_compaction_duration_seconds",
		Help:    "Duration of compaction runs",
		Buckets: prometheus.ExponentialBuckets(1, 2, 10),
	})
	duration.Observe(9)

	fsyncDuration := promauto.With(r).NewSummary(prometheus.SummaryOpts{
		Name:       "prometheus_tsdb_wal_fsync_duration_seconds",
		Help:       "Duration of WAL fsync.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
	fsyncDuration.Observe(10)

	pageFlushes := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_wal_page_flushes_total",
		Help: "Total number of page flushes.",
	})
	pageFlushes.Add(11 * base)

	pageCompletions := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_wal_completed_pages_total",
		Help: "Total number of completed pages.",
	})
	pageCompletions.Add(12 * base)

	truncateFail := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_wal_truncations_failed_total",
		Help: "Total number of WAL truncations that failed.",
	})
	truncateFail.Add(13 * base)

	truncateTotal := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_wal_truncations_total",
		Help: "Total number of WAL truncations attempted.",
	})
	truncateTotal.Add(14 * base)

	writesFailed := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_wal_writes_failed_total",
		Help: "Total number of WAL writes that failed.",
	})
	writesFailed.Add(15 * base)

	checkpointDeleteFail := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_checkpoint_deletions_failed_total",
		Help: "Total number of checkpoint deletions that failed.",
	})
	checkpointDeleteFail.Add(16 * base)

	checkpointDeleteTotal := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_checkpoint_deletions_total",
		Help: "Total number of checkpoint deletions attempted.",
	})
	checkpointDeleteTotal.Add(17 * base)

	checkpointCreationFail := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_checkpoint_creations_failed_total",
		Help: "Total number of checkpoint creations that failed.",
	})
	checkpointCreationFail.Add(18 * base)

	checkpointCreationTotal := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_checkpoint_creations_total",
		Help: "Total number of checkpoint creations attempted.",
	})
	checkpointCreationTotal.Add(19 * base)

	activeAppenders := promauto.With(r).NewGauge(prometheus.GaugeOpts{
		Name: "prometheus_tsdb_head_active_appenders",
		Help: "Number of currently active appender transactions",
	})
	activeAppenders.Set(20 * base)

	seriesNotFound := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_head_series_not_found_total",
		Help: "Total number of requests for series that were not found.",
	})
	seriesNotFound.Add(21 * base)

	chunks := promauto.With(r).NewGauge(prometheus.GaugeOpts{
		Name: "prometheus_tsdb_head_chunks",
		Help: "Total number of chunks in the head block.",
	})
	chunks.Set(22 * base)

	chunksCreated := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_head_chunks_created_total",
		Help: "Total number of chunks created in the head",
	})
	chunksCreated.Add(23 * base)

	chunksRemoved := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_head_chunks_removed_total",
		Help: "Total number of chunks removed in the head",
	})
	chunksRemoved.Add(24 * base)

	walTruncateDuration := promauto.With(r).NewSummary(prometheus.SummaryOpts{
		Name: "prometheus_tsdb_wal_truncate_duration_seconds",
		Help: "Duration of WAL truncation.",
	})
	walTruncateDuration.Observe(25)

	mmapChunkCorruptionTotal := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_mmap_chunk_corruptions_total",
		Help: "Total number of memory-mapped chunk corruptions.",
	})
	mmapChunkCorruptionTotal.Add(26 * base)

	mmapChunkQueueOperationsTotal := promauto.With(r).NewCounterVec(prometheus.CounterOpts{
		Name: "prometheus_tsdb_chunk_write_queue_operations_total",
		Help: "Number of operations on the chunk_write_queue.",
	}, []string{"operation"})
	mmapChunkQueueOperationsTotal.WithLabelValues("add").Add(50)
	mmapChunkQueueOperationsTotal.WithLabelValues("get").Add(10)
	mmapChunkQueueOperationsTotal.WithLabelValues("complete").Add(40)

	tsdbOOOHistogram := promauto.With(r).NewHistogram(prometheus.HistogramOpts{
		Name:    "prometheus_tsdb_sample_ooo_delta",
		Help:    "Delta in seconds by which a sample is considered out-of-order.",
		Buckets: []float64{60 * 10, 60 * 60 * 24}, // for testing: 3 buckets: 10 min, 24 hour, and inf
	})
	tsdbOOOHistogram.Observe(7 * base)

	walCorruptionsTotal := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_wal_corruptions_total",
		Help: "Total number of WAL corruptions.",
	})
	walCorruptionsTotal.Add(27 * base)

	headTruncateFail := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_head_truncations_failed_total",
		Help: "Total number of head truncations that failed.",
	})
	headTruncateFail.Add(28 * base)

	headTruncateTotal := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_head_truncations_total",
		Help: "Total number of head truncations attempted.",
	})
	headTruncateTotal.Add(29 * base)

	gcDuration := promauto.With(r).NewSummary(prometheus.SummaryOpts{
		Name: "prometheus_tsdb_head_gc_duration_seconds",
		Help: "Runtime of garbage collection in the head block.",
	})
	gcDuration.Observe(3)

	loadedBlocks := promauto.With(r).NewGauge(prometheus.GaugeOpts{
		Name: "prometheus_tsdb_blocks_loaded",
		Help: "Number of currently loaded data blocks",
	})
	loadedBlocks.Set(5)

	reloadsTotal := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_reloads_total",
		Help: "Number of times the database reloaded block data from disk.",
	})
	reloadsTotal.Add(10)

	reloadsFailed := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_reloads_failures_total",
		Help: "Number of times the database failed to reloadBlocks block data from disk.",
	})
	reloadsFailed.Add(7)

	symbolTableSize := promauto.With(r).NewGauge(prometheus.GaugeOpts{
		Name: "prometheus_tsdb_symbol_table_size_bytes",
		Help: "Size of symbol table in memory for loaded blocks",
	})
	symbolTableSize.Set(1024 * base)

	blocksSize := promauto.With(r).NewGauge(prometheus.GaugeOpts{
		Name: "prometheus_tsdb_storage_blocks_bytes",
		Help: "The number of bytes that are currently used for local storage by all blocks.",
	})
	blocksSize.Set(4096 * base)

	retentionsTotal := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_time_retentions_total",
		Help: "The number of times that blocks were deleted because the maximum time limit was exceeded.",
	})
	retentionsTotal.Add(11)

	exemplarsAppendedTotal := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_exemplar_exemplars_appended_total",
		Help: "Total number of appended exemplars.",
	})
	exemplarsAppendedTotal.Add(100)

	exemplarsStored := promauto.With(r).NewGauge(prometheus.GaugeOpts{
		Name: "prometheus_tsdb_exemplar_exemplars_in_storage",
		Help: "Number of exemplars currently in circular storage.",
	})
	exemplarsStored.Set(10)

	exemplarsSeriesStored := promauto.With(r).NewGauge(prometheus.GaugeOpts{
		Name: "prometheus_tsdb_exemplar_series_with_exemplars_in_storage",
		Help: "Number of series with exemplars currently in circular storage.",
	})
	exemplarsSeriesStored.Set(1)

	exemplarsLastTs := promauto.With(r).NewGauge(prometheus.GaugeOpts{
		Name: "prometheus_tsdb_exemplar_last_exemplars_timestamp_seconds",
		Help: "The timestamp of the oldest exemplar stored in circular storage. Useful to check for what time" +
			"range the current exemplar buffer limit allows. This usually means the last timestamp" +
			"for all exemplars for a typical setup. This is not true though if one of the series timestamp is in future compared to rest series.",
	})
	exemplarsLastTs.Set(1234)

	exemplarsOutOfOrderTotal := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_exemplar_out_of_order_exemplars_total",
		Help: "Total number of out-of-order exemplar ingestion failed attempts.",
	})
	exemplarsOutOfOrderTotal.Add(3)

	outOfOrderSamplesAppendedTotal := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_head_out_of_order_samples_appended_total",
		Help: "Total number of appended out-of-order samples.",
	})
	outOfOrderSamplesAppendedTotal.Add(3)

	chunksMmappedTotal := promauto.With(r).NewCounter(prometheus.CounterOpts{
		Name: "prometheus_tsdb_mmap_chunks_total",
		Help: "Total number of chunks that were memory-mapped.",
	})
	chunksMmappedTotal.Add(30 * base)

	tsdbWalReplayUnknownRefsTotal := promauto.With(r).NewCounterVec(prometheus.CounterOpts{
		Name: "prometheus_tsdb_wal_replay_unknown_refs_total",
		Help: "Total number of unknown series references encountered during WAL replay.",
	}, []string{"type"})
	tsdbWalReplayUnknownRefsTotal.WithLabelValues("series").Add(base)
	tsdbWalReplayUnknownRefsTotal.WithLabelValues("samples").Add(base * 2)

	tsdbWblReplayUnknownRefsTotal := promauto.With(r).NewCounterVec(prometheus.CounterOpts{
		Name: "prometheus_tsdb_wbl_replay_unknown_refs_total",
		Help: "Total number of unknown series references encountered during WBL replay pprus.",
	}, []string{"type"})
	tsdbWblReplayUnknownRefsTotal.WithLabelValues("exemplars").Add(base)

	return r
}
