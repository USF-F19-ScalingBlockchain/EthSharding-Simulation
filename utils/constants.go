package utils

var TOTAL_SHARDS uint32 = 2

/* Transaction pool thresholds */
var MIN_TX_POOL_SIZE = 1         // Minimum 50 Tx for given shard
var THRESHOLD_TX_POOL_SIZE = 150 // Maximum 150 Tx for given shard

/* Shard pool thresholds */
var MIN_SHARD_POOL_THRESHOLD = 1            // Minimum 1 Shard Root hash.
var MAX_SHARD_POOL_THRESHOLD = TOTAL_SHARDS // Maximum TOTAL_SHARDS Shard Root hash.
var SHARD_INTERVAL int32 = 1                // Push to beacon every nth block

/* Beacon Id */
var BEACON_ID uint32 = 9999

/*data construction*/
var Dataset = "cryptokitties"

//var NoOfShards = "4"
//var crossOrSame = suffix
var SubmissionRate = "1"
var ShardProductionType = "constant"
var BeaconProductionType = "constant"
var ShardProductionRate = "5"
var BeaconProductionRate = "10"
