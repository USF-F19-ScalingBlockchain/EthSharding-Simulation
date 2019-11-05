package utils

var TOTAL_SHARDS uint32 = 6

/* Transaction pool thresholds */
var MIN_TX_POOL_SIZE = 50 // Minimum 50 Tx for given shard
var THRESHOLD_TX_POOL_SIZE = 150 // Maximum 150 Tx for given shard

/* Shard pool thresholds */
var MIN_SHARD_POOL_THRESHOLD = 1 // Minimum 1 Shard Root hash.
var MAX_SHARD_POOL_THRESHOLD = TOTAL_SHARDS // Maximum TOTAL_SHARDS Shard Root hash.