package utils

var TOTAL_SHARDS uint32 = 0

/* Transaction pool thresholds */
var MIN_TX_POOL_SIZE = 1         // Minimum 50 Tx for given shard
var THRESHOLD_TX_POOL_SIZE = 150 // Maximum 150 Tx for given shard

/* Shard pool thresholds */
var MIN_SHARD_POOL_THRESHOLD = 1            // Minimum 1 Shard Root hash.
var MAX_SHARD_POOL_THRESHOLD = TOTAL_SHARDS // Maximum TOTAL_SHARDS Shard Root hash.
var SHARD_INTERVAL int32 = 2

/* Beacon Id */
var BEACON_ID uint32 = 9999
