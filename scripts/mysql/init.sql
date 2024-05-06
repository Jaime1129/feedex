CREATE DATABASE IF NOT EXISTS trx_fee;
USE trx_fee;

-- trx_fee.uni_trx_fee definition

CREATE TABLE IF NOT EXISTS `uni_trx_fee` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'auto-generated primary key',
  `symbol` varchar(100) NOT NULL DEFAULT 'WETH/USDC' COMMENT 'symbol',
  `trx_hash` varchar(100) NOT NULL DEFAULT '' COMMENT 'transaction hash',
  `trx_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT 'timestamp of transaction',
  `block_num` bigint unsigned NOT NULL DEFAULT '0',
  `gas_used` bigint unsigned NOT NULL DEFAULT '0',
  `gas_price` bigint unsigned NOT NULL DEFAULT '0',
  `eth_usdt_price` decimal(10,0) NOT NULL DEFAULT '0',
  `trx_fee_usdt` decimal(10,0) NOT NULL DEFAULT '0',
  `gmt_created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `trx_hash_unique` (`trx_hash`),
  KEY `uni_trx_fee_trx_time_IDX` (`trx_time`) USING BTREE,
  KEY `uni_trx_fee_block_num_IDX` (`block_num`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=217 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- trx_fee.block_num_record definition

CREATE TABLE IF NOT EXISTS `block_num_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `max_block` bigint unsigned NOT NULL DEFAULT '0',
  `symbol` varchar(100) NOT NULL DEFAULT 'WETH/USDC',
  PRIMARY KEY (`id`),
  UNIQUE KEY `block_num_record_symbol_IDX` (`symbol`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;