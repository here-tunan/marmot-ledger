-- Marmot Ledger initial schema.
-- Execute after sql/00_create_database.sql.
-- Important domain decisions:
--   * Use financial_event, ledger_entry, and bucket as the core ledger model.
--   * Do not create a transaction table.
--   * account does not store balance; bucket stores balance.
--   * ledger_entry.amount is signed; no direction column is used.

USE `marmot_ledger`;
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 1;

CREATE TABLE IF NOT EXISTS `user` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `account` VARCHAR(64) NOT NULL COMMENT '登录账号',
  `password` VARCHAR(255) NOT NULL COMMENT '登录密码',
  `name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `desc` VARCHAR(512) NOT NULL DEFAULT '' COMMENT '用户描述',
  `avatar` VARCHAR(512) NOT NULL DEFAULT '' COMMENT '头像地址',
  `extra` TEXT NULL COMMENT '扩展信息',
  `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `gmt_create` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `gmt_modified` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_account` (`account`),
  KEY `idx_user_is_deleted` (`is_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户';

CREATE TABLE IF NOT EXISTS `currency` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `code` CHAR(3) NOT NULL COMMENT '币种代码',
  `name` VARCHAR(64) NOT NULL COMMENT '币种名称',
  `symbol` VARCHAR(16) NOT NULL DEFAULT '' COMMENT '币种符号',
  `precision_digits` TINYINT UNSIGNED NOT NULL DEFAULT 2 COMMENT '金额精度',
  `enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_currency_code` (`code`),
  KEY `idx_currency_enabled_sort` (`enabled`, `sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='币种字典';

CREATE TABLE IF NOT EXISTS `account_template` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `provider_code` VARCHAR(64) NOT NULL COMMENT '模板代码',
  `name` VARCHAR(64) NOT NULL COMMENT '模板名称',
  `type` VARCHAR(32) NOT NULL COMMENT '账户类型',
  `icon` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Element Plus 图标 key',
  `color` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '颜色',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_account_template_provider_code` (`provider_code`),
  KEY `idx_account_template_type_enabled_sort` (`type`, `enabled`, `sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户模板';

CREATE TABLE IF NOT EXISTS `category_group` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `group_code` VARCHAR(64) NOT NULL COMMENT '聚合分类代码',
  `name` VARCHAR(64) NOT NULL COMMENT '聚合分类名称',
  `type` VARCHAR(32) NOT NULL COMMENT 'income/expense',
  `icon` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '图标',
  `color` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '颜色',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_category_group_code` (`group_code`),
  KEY `idx_category_group_type_enabled_sort` (`type`, `enabled`, `sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='统计聚合分类';

CREATE TABLE IF NOT EXISTS `channel_template` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `channel_code` VARCHAR(64) NOT NULL COMMENT '渠道代码',
  `name` VARCHAR(64) NOT NULL COMMENT '渠道名称',
  `channel_type` VARCHAR(32) NOT NULL COMMENT '渠道类型',
  `provider_code` VARCHAR(64) NULL COMMENT '平台或机构代码',
  `supported_event_types` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '支持的事件类型，逗号分隔',
  `icon` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '图标',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  `remark` VARCHAR(512) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_channel_template_code` (`channel_code`),
  KEY `idx_channel_template_type_enabled_sort` (`channel_type`, `enabled`, `sort`),
  KEY `idx_channel_template_provider_code` (`provider_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收付款渠道模板';

CREATE TABLE IF NOT EXISTS `exchange_rate` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` BIGINT NULL COMMENT '用户ID，NULL表示全局汇率',
  `from_currency` CHAR(3) NOT NULL COMMENT '源币种',
  `to_currency` CHAR(3) NOT NULL COMMENT '目标币种',
  `rate` DECIMAL(24,12) NOT NULL COMMENT '汇率',
  `rate_date` DATE NOT NULL COMMENT '汇率日期',
  `source` VARCHAR(64) NOT NULL DEFAULT 'manual' COMMENT '汇率来源',
  `is_manual` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否手动录入',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_exchange_rate_lookup` (`from_currency`, `to_currency`, `rate_date`),
  KEY `idx_exchange_rate_user_lookup` (`user_id`, `from_currency`, `to_currency`, `rate_date`),
  CONSTRAINT `fk_exchange_rate_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`),
  CONSTRAINT `fk_exchange_rate_from_currency` FOREIGN KEY (`from_currency`) REFERENCES `currency` (`code`),
  CONSTRAINT `fk_exchange_rate_to_currency` FOREIGN KEY (`to_currency`) REFERENCES `currency` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='汇率缓存';

CREATE TABLE IF NOT EXISTS `account` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `name` VARCHAR(128) NOT NULL COMMENT '账户名称',
  `type` VARCHAR(32) NOT NULL COMMENT '账户类型',
  `icon` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Element Plus 图标 key',
  `color` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '颜色',
  `is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_account_user_active` (`user_id`, `is_deleted`, `is_active`),
  KEY `idx_account_user_type` (`user_id`, `type`),
  CONSTRAINT `fk_account_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户，不保存余额';

CREATE TABLE IF NOT EXISTS `bucket` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `account_id` BIGINT NOT NULL COMMENT '账户ID',
  `name` VARCHAR(128) NOT NULL COMMENT 'Bucket名称',
  `currency` CHAR(3) NOT NULL COMMENT '币种',
  `balance` DECIMAL(20,4) NOT NULL DEFAULT 0.0000 COMMENT '当前余额',
  `initial_balance` DECIMAL(20,4) NOT NULL DEFAULT 0.0000 COMMENT '初始余额',
  `bucket_type` VARCHAR(32) NOT NULL COMMENT 'Bucket类型',
  `bucket_nature` VARCHAR(32) NOT NULL COMMENT 'asset/liability',
  `bucket_group_key` VARCHAR(64) NULL COMMENT 'Bucket聚合键',
  `is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_bucket_user_account_active` (`user_id`, `account_id`, `is_deleted`, `is_active`),
  KEY `idx_bucket_user_type` (`user_id`, `bucket_type`),
  KEY `idx_bucket_user_nature_currency` (`user_id`, `bucket_nature`, `currency`),
  KEY `idx_bucket_currency` (`currency`),
  CONSTRAINT `fk_bucket_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`),
  CONSTRAINT `fk_bucket_account` FOREIGN KEY (`account_id`) REFERENCES `account` (`id`),
  CONSTRAINT `fk_bucket_currency` FOREIGN KEY (`currency`) REFERENCES `currency` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='资产/负债/虚拟资金容器';

CREATE TABLE IF NOT EXISTS `category` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `name` VARCHAR(64) NOT NULL COMMENT '分类名称',
  `type` VARCHAR(32) NOT NULL COMMENT 'income/expense',
  `category_group_id` BIGINT NOT NULL COMMENT '聚合分类ID',
  `is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_category_user_type_name_deleted` (`user_id`, `type`, `name`, `is_deleted`),
  KEY `idx_category_user_type_active` (`user_id`, `type`, `is_deleted`, `is_active`),
  KEY `idx_category_group_id` (`category_group_id`),
  CONSTRAINT `fk_category_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`),
  CONSTRAINT `fk_category_group` FOREIGN KEY (`category_group_id`) REFERENCES `category_group` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户分类';

CREATE TABLE IF NOT EXISTS `financial_event` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `event_group_id` BIGINT NULL COMMENT '事件组ID',
  `related_financial_event_id` BIGINT NULL COMMENT '关联财务事件ID',
  `event_type` VARCHAR(32) NOT NULL COMMENT '事件类型',
  `description` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '描述',
  `category_id` BIGINT NULL COMMENT '用户分类ID',
  `category_group_id` BIGINT NULL COMMENT '聚合分类ID',
  `channel_type` VARCHAR(64) NULL COMMENT '渠道代码，对应channel_template.channel_code',
  `channel_account_id` BIGINT NULL COMMENT '渠道账户ID',
  `event_time` DATETIME(3) NOT NULL COMMENT '事件发生时间',
  `currency` CHAR(3) NOT NULL COMMENT '原始币种',
  `amount` DECIMAL(20,4) NOT NULL DEFAULT 0.0000 COMMENT '原始金额',
  `include_in_statistics` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否进入普通收支统计',
  `source` VARCHAR(32) NOT NULL DEFAULT 'manual' COMMENT '来源',
  `status` VARCHAR(32) NOT NULL DEFAULT 'active' COMMENT '状态',
  `remark` TEXT NULL COMMENT '备注',
  `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_financial_event_user_time` (`user_id`, `event_time`),
  KEY `idx_financial_event_user_type_time` (`user_id`, `event_type`, `event_time`),
  KEY `idx_financial_event_user_category_group_time` (`user_id`, `category_group_id`, `event_time`),
  KEY `idx_financial_event_event_group` (`event_group_id`),
  KEY `idx_financial_event_related` (`related_financial_event_id`),
  KEY `idx_financial_event_channel_time` (`channel_type`, `event_time`),
  KEY `idx_financial_event_user_deleted_time` (`user_id`, `is_deleted`, `event_time`),
  KEY `idx_financial_event_currency` (`currency`),
  CONSTRAINT `fk_financial_event_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`),
  CONSTRAINT `fk_financial_event_related` FOREIGN KEY (`related_financial_event_id`) REFERENCES `financial_event` (`id`),
  CONSTRAINT `fk_financial_event_category` FOREIGN KEY (`category_id`) REFERENCES `category` (`id`),
  CONSTRAINT `fk_financial_event_category_group` FOREIGN KEY (`category_group_id`) REFERENCES `category_group` (`id`),
  CONSTRAINT `fk_financial_event_channel` FOREIGN KEY (`channel_type`) REFERENCES `channel_template` (`channel_code`),
  CONSTRAINT `fk_financial_event_channel_account` FOREIGN KEY (`channel_account_id`) REFERENCES `account` (`id`),
  CONSTRAINT `fk_financial_event_currency` FOREIGN KEY (`currency`) REFERENCES `currency` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='财务事件';

CREATE TABLE IF NOT EXISTS `ledger_entry` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `financial_event_id` BIGINT NOT NULL COMMENT '财务事件ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `bucket_id` BIGINT NOT NULL COMMENT 'Bucket ID',
  `currency` CHAR(3) NOT NULL COMMENT '币种',
  `amount` DECIMAL(20,4) NOT NULL COMMENT '余额变化金额，正数增加，负数减少',
  `balance_after` DECIMAL(20,4) NOT NULL COMMENT '变更后余额',
  `entry_role` VARCHAR(32) NULL COMMENT '分录角色',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_ledger_entry_event` (`financial_event_id`),
  KEY `idx_ledger_entry_user_bucket_created` (`user_id`, `bucket_id`, `created_at`),
  KEY `idx_ledger_entry_bucket_created` (`bucket_id`, `created_at`),
  KEY `idx_ledger_entry_user_created` (`user_id`, `created_at`),
  KEY `idx_ledger_entry_currency` (`currency`),
  CONSTRAINT `fk_ledger_entry_event` FOREIGN KEY (`financial_event_id`) REFERENCES `financial_event` (`id`),
  CONSTRAINT `fk_ledger_entry_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`),
  CONSTRAINT `fk_ledger_entry_bucket` FOREIGN KEY (`bucket_id`) REFERENCES `bucket` (`id`),
  CONSTRAINT `fk_ledger_entry_currency` FOREIGN KEY (`currency`) REFERENCES `currency` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='余额分录';

CREATE TABLE IF NOT EXISTS `exchange_detail` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `financial_event_id` BIGINT NOT NULL COMMENT '财务事件ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `from_bucket_id` BIGINT NOT NULL COMMENT '换出Bucket ID',
  `to_bucket_id` BIGINT NOT NULL COMMENT '换入Bucket ID',
  `from_currency` CHAR(3) NOT NULL COMMENT '换出币种',
  `to_currency` CHAR(3) NOT NULL COMMENT '换入币种',
  `from_amount` DECIMAL(20,4) NOT NULL COMMENT '换出金额',
  `to_amount` DECIMAL(20,4) NOT NULL COMMENT '换入金额',
  `exchange_rate` DECIMAL(24,12) NOT NULL COMMENT '汇率',
  `fee_amount` DECIMAL(20,4) NOT NULL DEFAULT 0.0000 COMMENT '手续费金额',
  `fee_currency` CHAR(3) NULL COMMENT '手续费币种',
  `fee_bucket_id` BIGINT NULL COMMENT '手续费Bucket ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_exchange_detail_event` (`financial_event_id`),
  KEY `idx_exchange_detail_user_created` (`user_id`, `created_at`),
  KEY `idx_exchange_detail_from_to_currency` (`from_currency`, `to_currency`),
  KEY `idx_exchange_detail_to_currency` (`to_currency`),
  KEY `idx_exchange_detail_fee_currency` (`fee_currency`),
  KEY `idx_exchange_detail_from_bucket` (`from_bucket_id`),
  KEY `idx_exchange_detail_to_bucket` (`to_bucket_id`),
  KEY `idx_exchange_detail_fee_bucket` (`fee_bucket_id`),
  CONSTRAINT `fk_exchange_detail_event` FOREIGN KEY (`financial_event_id`) REFERENCES `financial_event` (`id`),
  CONSTRAINT `fk_exchange_detail_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`),
  CONSTRAINT `fk_exchange_detail_from_bucket` FOREIGN KEY (`from_bucket_id`) REFERENCES `bucket` (`id`),
  CONSTRAINT `fk_exchange_detail_to_bucket` FOREIGN KEY (`to_bucket_id`) REFERENCES `bucket` (`id`),
  CONSTRAINT `fk_exchange_detail_fee_bucket` FOREIGN KEY (`fee_bucket_id`) REFERENCES `bucket` (`id`),
  CONSTRAINT `fk_exchange_detail_from_currency` FOREIGN KEY (`from_currency`) REFERENCES `currency` (`code`),
  CONSTRAINT `fk_exchange_detail_to_currency` FOREIGN KEY (`to_currency`) REFERENCES `currency` (`code`),
  CONSTRAINT `fk_exchange_detail_fee_currency` FOREIGN KEY (`fee_currency`) REFERENCES `currency` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='换汇详情';

CREATE TABLE IF NOT EXISTS `investment_asset` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `account_id` BIGINT NOT NULL COMMENT '投资账户ID',
  `name` VARCHAR(128) NOT NULL COMMENT '资产名称',
  `code` VARCHAR(64) NOT NULL COMMENT '资产代码',
  `asset_type` VARCHAR(32) NOT NULL COMMENT '资产类型',
  `currency` CHAR(3) NOT NULL COMMENT '币种',
  `market` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '市场',
  `is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_investment_asset_user_market_code_deleted` (`user_id`, `market`, `code`, `is_deleted`),
  KEY `idx_investment_asset_user_account` (`user_id`, `account_id`),
  KEY `idx_investment_asset_currency` (`currency`),
  CONSTRAINT `fk_investment_asset_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`),
  CONSTRAINT `fk_investment_asset_account` FOREIGN KEY (`account_id`) REFERENCES `account` (`id`),
  CONSTRAINT `fk_investment_asset_currency` FOREIGN KEY (`currency`) REFERENCES `currency` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='投资资产';

CREATE TABLE IF NOT EXISTS `investment_snapshot` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `account_id` BIGINT NOT NULL COMMENT '投资账户ID',
  `asset_id` BIGINT NOT NULL COMMENT '投资资产ID',
  `snapshot_date` DATE NOT NULL COMMENT '快照日期',
  `currency` CHAR(3) NOT NULL COMMENT '币种',
  `principal_amount` DECIMAL(20,4) NOT NULL DEFAULT 0.0000 COMMENT '本金',
  `market_value` DECIMAL(20,4) NOT NULL DEFAULT 0.0000 COMMENT '市值',
  `realized_profit` DECIMAL(20,4) NOT NULL DEFAULT 0.0000 COMMENT '已实现收益',
  `unrealized_profit` DECIMAL(20,4) NOT NULL DEFAULT 0.0000 COMMENT '未实现收益',
  `total_profit` DECIMAL(20,4) NOT NULL DEFAULT 0.0000 COMMENT '总收益',
  `profit_rate` DECIMAL(18,8) NOT NULL DEFAULT 0.00000000 COMMENT '收益率',
  `base_currency` CHAR(3) NOT NULL COMMENT '本位币',
  `base_market_value` DECIMAL(20,4) NOT NULL DEFAULT 0.0000 COMMENT '本位币市值',
  `base_total_profit` DECIMAL(20,4) NOT NULL DEFAULT 0.0000 COMMENT '本位币总收益',
  `exchange_rate` DECIMAL(24,12) NOT NULL DEFAULT 1.000000000000 COMMENT '汇率快照',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_investment_snapshot_asset_date` (`asset_id`, `snapshot_date`),
  KEY `idx_investment_snapshot_user_date` (`user_id`, `snapshot_date`),
  KEY `idx_investment_snapshot_account_date` (`account_id`, `snapshot_date`),
  KEY `idx_investment_snapshot_currency` (`currency`),
  KEY `idx_investment_snapshot_base_currency` (`base_currency`),
  CONSTRAINT `fk_investment_snapshot_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`),
  CONSTRAINT `fk_investment_snapshot_account` FOREIGN KEY (`account_id`) REFERENCES `account` (`id`),
  CONSTRAINT `fk_investment_snapshot_asset` FOREIGN KEY (`asset_id`) REFERENCES `investment_asset` (`id`),
  CONSTRAINT `fk_investment_snapshot_currency` FOREIGN KEY (`currency`) REFERENCES `currency` (`code`),
  CONSTRAINT `fk_investment_snapshot_base_currency` FOREIGN KEY (`base_currency`) REFERENCES `currency` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='投资估值与收益快照';

CREATE TABLE IF NOT EXISTS `family` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` VARCHAR(128) NOT NULL COMMENT '家庭名称',
  `owner_user_id` BIGINT NOT NULL COMMENT '创建者用户ID',
  `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_family_owner` (`owner_user_id`, `is_deleted`),
  CONSTRAINT `fk_family_owner_user` FOREIGN KEY (`owner_user_id`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='家庭';

CREATE TABLE IF NOT EXISTS `family_member` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `family_id` BIGINT NOT NULL COMMENT '家庭ID',
  `user_id` BIGINT NOT NULL COMMENT '成员用户ID',
  `role` VARCHAR(32) NOT NULL DEFAULT 'member' COMMENT 'owner/admin/member',
  `status` VARCHAR(32) NOT NULL DEFAULT 'invited' COMMENT 'invited/active/rejected/left',
  `display_name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '家庭内展示名',
  `invited_by_user_id` BIGINT NULL COMMENT '邀请人用户ID',
  `invited_at` DATETIME(3) NULL COMMENT '邀请时间',
  `joined_at` DATETIME(3) NULL COMMENT '加入时间',
  `left_at` DATETIME(3) NULL COMMENT '离开时间',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_family_member_family_user` (`family_id`, `user_id`),
  KEY `idx_family_member_user_status` (`user_id`, `status`),
  KEY `idx_family_member_family_status` (`family_id`, `status`),
  KEY `idx_family_member_inviter` (`invited_by_user_id`),
  CONSTRAINT `fk_family_member_family` FOREIGN KEY (`family_id`) REFERENCES `family` (`id`),
  CONSTRAINT `fk_family_member_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`),
  CONSTRAINT `fk_family_member_inviter` FOREIGN KEY (`invited_by_user_id`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='家庭成员';
