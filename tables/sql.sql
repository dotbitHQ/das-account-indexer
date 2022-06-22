SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE DATABASE IF NOT EXISTS `das_account_indexer`;
USE `das_account_indexer`;

-- ----------------------------
-- Table structure for t_block_info
-- ----------------------------
CREATE TABLE IF NOT EXISTS `t_block_info`
(
    `id`           bigint(20) unsigned                                           NOT NULL AUTO_INCREMENT COMMENT '',
    `block_number` bigint(20) unsigned                                           NOT NULL DEFAULT '0' COMMENT '',
    `block_hash`   varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '',
    `parent_hash`  varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '',
    `created_at`   timestamp                                                     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '',
    `updated_at`   timestamp                                                     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_block_number` (`block_number`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='for block rollback';

-- ----------------------------
-- Table structure for t_account_info
-- ----------------------------
CREATE TABLE IF NOT EXISTS `t_account_info`
(
    `id`                      bigint(20) unsigned                                           NOT NULL AUTO_INCREMENT COMMENT '',
    `block_number`            bigint(20) unsigned                                           NOT NULL DEFAULT '0' COMMENT '',
    `block_timestamp`         BIGINT(20)                                                    NOT NULL DEFAULT '0' COMMENT '',
    `outpoint`                varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'Hash-Index',
    `account_id`              varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'hash of account',
    `parent_account_id`       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '',
    `next_account_id`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'hash of next account',
    `account`                 varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '',
    `owner_chain_type`        smallint(6)                                                   NOT NULL DEFAULT '0' COMMENT '',
    `owner`                   varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'owner address',
    `owner_algorithm_id`      smallint(6)                                                   NOT NULL DEFAULT '0' COMMENT '',
    `manager_chain_type`      smallint(6)                                                   NOT NULL DEFAULT '0' COMMENT '',
    `manager`                 varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'manager address',
    `manager_algorithm_id`    smallint(6)                                                   NOT NULL DEFAULT '0' COMMENT '',
    `status`                  smallint(6)                                                   NOT NULL DEFAULT '0' COMMENT '',
    `enable_sub_account`      smallint(6)                                                   NOT NULL DEFAULT '0' COMMENT '',
    `renew_sub_account_price` bigint(20)                                                    NOT NULL DEFAULT '0' COMMENT '',
    `nonce`                   bigint(20)                                                    NOT NULL DEFAULT '0' COMMENT '',
    `registered_at`           bigint(20) unsigned                                           NOT NULL DEFAULT '0' COMMENT '',
    `expired_at`              bigint(20) unsigned                                           NOT NULL DEFAULT '0' COMMENT '',
    `created_at`              timestamp                                                     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '',
    `updated_at`              timestamp                                                     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_account_id` (`account_id`) USING BTREE,
    KEY `k_account` (`account`) USING BTREE,
    KEY `k_next_account_id` (`next_account_id`) USING BTREE,
    KEY `k_oct_o` (`owner_chain_type`, `owner`) USING BTREE,
    KEY `k_mct_m` (`manager_chain_type`, `manager`) USING BTREE,
    KEY `k_parent_account_id` (`parent_account_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='current account info';

-- ----------------------------
-- Table structure for t_records_info
-- ----------------------------
CREATE TABLE IF NOT EXISTS `t_records_info`
(
    `id`                bigint(20) unsigned                                            NOT NULL AUTO_INCREMENT COMMENT '',
    `account_id`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci  NOT NULL DEFAULT '' COMMENT 'hash of account',
    `parent_account_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci  NOT NULL DEFAULT '' COMMENT '',
    `account`           varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci  NOT NULL DEFAULT '' COMMENT '',
    `key`               varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci  NOT NULL DEFAULT '',
    `type`              varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci  NOT NULL DEFAULT '',
    `label`             varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci  NOT NULL DEFAULT '',
    `value`             varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
    `ttl`               varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci  NOT NULL DEFAULT '',
    `created_at`        timestamp                                                      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '',
    `updated_at`        timestamp                                                      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `k_account_id` (`account_id`) USING BTREE,
    KEY `k_account` (`account`) USING BTREE,
    KEY `k_value` (`value`(768)) USING BTREE,
    KEY `k_parent_account_id` (`parent_account_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='records info in DAS account setting';

-- ----------------------------
-- Table structure for t_reverse_info
-- ----------------------------
CREATE TABLE IF NOT EXISTS `t_reverse_info`
(
    `id`              BIGINT(20) UNSIGNED                                           NOT NULL AUTO_INCREMENT COMMENT '',
    `block_number`    BIGINT(20)                                                    NOT NULL DEFAULT '0' COMMENT '',
    `block_timestamp` BIGINT(20)                                                    NOT NULL DEFAULT '0' COMMENT '',
    `outpoint`        VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '',
    `algorithm_id`    SMALLINT(6)                                                   NOT NULL DEFAULT '0' COMMENT '',
    `chain_type`      SMALLINT(6)                                                   NOT NULL DEFAULT '0' COMMENT '',
    `address`         VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '',
    `account`         VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '',
    `capacity`        BIGINT(20)                                                    NOT NULL DEFAULT '0' COMMENT '',
    `created_at`      timestamp                                                     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '',
    `updated_at`      timestamp                                                     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '',
    PRIMARY KEY (id),
    UNIQUE KEY uk_outpoint (outpoint),
    KEY k_address (chain_type, address),
    KEY k_account (account)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='reverse records info';

-- # DROP TABLES
-- # DROP TABLE IF EXISTS `t_block_info`;
-- # DROP TABLE IF EXISTS `t_account_info`;
-- # DROP TABLE IF EXISTS `t_records_info`;
-- # DROP TABLE IF EXISTS `t_reverse_info`;