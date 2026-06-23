-- Marmot Ledger unified database bootstrap.
-- Use this file for a fresh local/development database:
--   mysql < sql/marmot_ledger.sql
--
-- This is the single source of truth for the current MVP schema and required
-- reference data. It intentionally does not create demo users or demo records.

CREATE DATABASE IF NOT EXISTS `marmot_ledger`
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE `marmot_ledger`;
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 1;

-- Important domain decisions:
--   * Use financial_event, ledger_entry, and bucket as the core ledger model.
--   * Do not create a transaction table.
--   * account does not store balance; bucket stores balance.
--   * ledger_entry.amount is signed; no direction column is used.

CREATE TABLE IF NOT EXISTS `user` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `account` VARCHAR(64) NOT NULL COMMENT '登录账号',
  `password` VARCHAR(255) NOT NULL COMMENT '登录密码',
  `name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `role` VARCHAR(32) NOT NULL DEFAULT 'user' COMMENT '用户角色: user/admin',
  `desc` VARCHAR(512) NOT NULL DEFAULT '' COMMENT '用户描述',
  `avatar` VARCHAR(512) NOT NULL DEFAULT '' COMMENT '头像地址',
  `extra` TEXT NULL COMMENT '扩展信息',
  `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `gmt_create` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `gmt_modified` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_account` (`account`),
  KEY `idx_user_role` (`role`),
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

CREATE TABLE IF NOT EXISTS `category_template` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `template_code` VARCHAR(64) NOT NULL COMMENT '模板代码',
  `name` VARCHAR(64) NOT NULL COMMENT '模板分类名称',
  `type` VARCHAR(32) NOT NULL COMMENT 'income/expense',
  `icon` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '图标',
  `color` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '颜色',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_category_template_code` (`template_code`),
  KEY `idx_category_template_type_enabled_sort` (`type`, `enabled`, `sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分类模板';

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

CREATE TABLE IF NOT EXISTS `personal_channel` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `channel_template_id` BIGINT NULL COMMENT '来源渠道模板ID',
  `name` VARCHAR(64) NOT NULL COMMENT '渠道名称',
  `channel_type` VARCHAR(32) NOT NULL COMMENT '渠道类型',
  `provider_code` VARCHAR(64) NULL COMMENT '平台或机构代码',
  `supported_event_types` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '支持的事件类型，逗号分隔',
  `icon` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '图标',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_personal_channel_user_name_deleted` (`user_id`, `name`, `is_deleted`),
  KEY `idx_personal_channel_user_active` (`user_id`, `is_deleted`, `is_active`, `sort`),
  KEY `idx_personal_channel_template` (`channel_template_id`),
  KEY `idx_personal_channel_type` (`channel_type`),
  CONSTRAINT `fk_personal_channel_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`),
  CONSTRAINT `fk_personal_channel_template` FOREIGN KEY (`channel_template_id`) REFERENCES `channel_template` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='个人收付款渠道';

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
  `icon` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '图标',
  `color` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '颜色',
  `is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_category_user_type_name_deleted` (`user_id`, `type`, `name`, `is_deleted`),
  KEY `idx_category_user_type_active` (`user_id`, `type`, `is_deleted`, `is_active`),
  CONSTRAINT `fk_category_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户分类';

CREATE TABLE IF NOT EXISTS `financial_event` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `event_group_id` BIGINT NULL COMMENT '事件组ID',
  `related_financial_event_id` BIGINT NULL COMMENT '关联财务事件ID',
  `event_type` VARCHAR(32) NOT NULL COMMENT '事件类型',
  `description` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '描述',
  `category_id` BIGINT NULL COMMENT '用户分类ID',
  `channel_id` BIGINT NULL COMMENT '个人渠道ID',
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
  KEY `idx_financial_event_user_category_time` (`user_id`, `category_id`, `event_time`),
  KEY `idx_financial_event_event_group` (`event_group_id`),
  KEY `idx_financial_event_related` (`related_financial_event_id`),
  KEY `idx_financial_event_user_channel_time` (`user_id`, `channel_id`, `event_time`),
  KEY `idx_financial_event_user_deleted_time` (`user_id`, `is_deleted`, `event_time`),
  KEY `idx_financial_event_currency` (`currency`),
  CONSTRAINT `fk_financial_event_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`),
  CONSTRAINT `fk_financial_event_related` FOREIGN KEY (`related_financial_event_id`) REFERENCES `financial_event` (`id`),
  CONSTRAINT `fk_financial_event_category` FOREIGN KEY (`category_id`) REFERENCES `category` (`id`),
  CONSTRAINT `fk_financial_event_channel` FOREIGN KEY (`channel_id`) REFERENCES `personal_channel` (`id`),
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

CREATE TABLE IF NOT EXISTS `family_category_group` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `family_id` BIGINT NOT NULL COMMENT '家庭ID',
  `name` VARCHAR(64) NOT NULL COMMENT '分组名称',
  `type` VARCHAR(32) NOT NULL COMMENT 'income/expense',
  `icon` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '图标',
  `color` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '颜色',
  `created_by_user_id` BIGINT NOT NULL COMMENT '创建者用户ID',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_family_category_group_family_active` (`family_id`, `is_deleted`, `is_active`),
  KEY `idx_family_category_group_type` (`type`),
  CONSTRAINT `fk_family_category_group_family` FOREIGN KEY (`family_id`) REFERENCES `family` (`id`),
  CONSTRAINT `fk_family_category_group_creator` FOREIGN KEY (`created_by_user_id`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='家庭分类组';

CREATE TABLE IF NOT EXISTS `family_category_group_member` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `family_group_id` BIGINT NOT NULL COMMENT '家庭分类组ID',
  `category_id` BIGINT NOT NULL COMMENT '用户分类ID',
  `added_by_user_id` BIGINT NOT NULL COMMENT '添加者用户ID',
  `added_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '添加时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_family_category_group_member` (`family_group_id`, `category_id`),
  KEY `idx_family_category_group_member_category` (`category_id`),
  CONSTRAINT `fk_family_category_group_member_group` FOREIGN KEY (`family_group_id`) REFERENCES `family_category_group` (`id`),
  CONSTRAINT `fk_family_category_group_member_category` FOREIGN KEY (`category_id`) REFERENCES `category` (`id`),
  CONSTRAINT `fk_family_category_group_member_adder` FOREIGN KEY (`added_by_user_id`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='家庭分类组成员';

-- Required reference data. Safe to rerun.

INSERT INTO `currency` (`code`, `name`, `symbol`, `precision_digits`, `enabled`, `sort`)
VALUES
  ('CNY', '人民币', '¥', 2, 1, 10),
  ('USD', '美元', '$', 2, 1, 20),
  ('HKD', '港币', 'HK$', 2, 1, 30),
  ('EUR', '欧元', '€', 2, 1, 40),
  ('JPY', '日元', '¥', 0, 1, 50),
  ('GBP', '英镑', '£', 2, 1, 60),
  ('SGD', '新加坡元', 'S$', 2, 1, 70),
  ('AUD', '澳大利亚元', 'A$', 2, 1, 80),
  ('NZD', '新西兰元', 'NZ$', 2, 1, 90)
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `symbol` = VALUES(`symbol`),
  `precision_digits` = VALUES(`precision_digits`),
  `enabled` = VALUES(`enabled`),
  `sort` = VALUES(`sort`);

INSERT INTO `account_template` (`provider_code`, `name`, `type`, `icon`, `color`, `sort`, `enabled`)
VALUES
  ('WECHAT', '微信钱包', 'wallet', 'Wallet', '#22c55e', 10, 1),
  ('ALIPAY', '支付宝', 'wallet', 'Wallet', '#1677ff', 20, 1),
  ('CMB', '招商银行', 'bank', 'CreditCard', '#ef4444', 30, 1),
  ('ICBC', '工商银行', 'bank', 'CreditCard', '#f97316', 40, 1),
  ('CCB', '建设银行', 'bank', 'CreditCard', '#1d4ed8', 50, 1),
  ('ABC', '农业银行', 'bank', 'CreditCard', '#16a34a', 60, 1),
  ('BOC', '中国银行', 'bank', 'CreditCard', '#dc2626', 70, 1),
  ('CASH', '现金', 'cash', 'Money', '#f59e0b', 80, 1),
  ('IBKR', 'IBKR', 'investment', 'TrendCharts', '#1f2933', 90, 1),
  ('FUTU', '富途证券', 'investment', 'TrendCharts', '#06b6d4', 100, 1),
  ('OTHER', '其他账户', 'other', 'Wallet', '#2f7d5c', 999, 1)
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `type` = VALUES(`type`),
  `icon` = VALUES(`icon`),
  `color` = VALUES(`color`),
  `sort` = VALUES(`sort`),
  `enabled` = VALUES(`enabled`);

INSERT INTO `category_group` (`group_code`, `name`, `type`, `icon`, `color`, `sort`, `enabled`)
VALUES
  ('FOOD', '餐饮', 'expense', '', '#f97316', 10, 1),
  ('GROCERY', '买菜日用', 'expense', '', '#f97316', 20, 1),
  ('TRANSPORT', '交通', 'expense', '', '#f97316', 30, 1),
  ('HOUSING', '住房', 'expense', '', '#f97316', 40, 1),
  ('UTILITIES', '水电燃气', 'expense', '', '#f97316', 50, 1),
  ('SHOPPING', '购物', 'expense', '', '#f97316', 60, 1),
  ('ENTERTAINMENT', '娱乐', 'expense', '', '#f97316', 70, 1),
  ('HEALTHCARE', '医疗健康', 'expense', '', '#f97316', 80, 1),
  ('EDUCATION', '教育学习', 'expense', '', '#f97316', 90, 1),
  ('TRAVEL', '旅行', 'expense', '', '#f97316', 100, 1),
  ('OTHER_EXPENSE', '其他支出', 'expense', '', '#f97316', 999, 1),
  ('SALARY', '工资', 'income', '', '#ef4444', 1010, 1),
  ('BONUS', '奖金', 'income', '', '#ef4444', 1020, 1),
  ('INVESTMENT_INCOME', '投资收益', 'income', '', '#ef4444', 1030, 1),
  ('REIMBURSEMENT', '报销', 'income', '', '#ef4444', 1040, 1),
  ('OTHER_INCOME', '其他收入', 'income', '', '#ef4444', 1999, 1)
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `type` = VALUES(`type`),
  `icon` = VALUES(`icon`),
  `color` = VALUES(`color`),
  `sort` = VALUES(`sort`),
  `enabled` = VALUES(`enabled`);

INSERT INTO `channel_template` (`channel_code`, `name`, `channel_type`, `provider_code`, `supported_event_types`, `icon`, `sort`, `enabled`, `remark`)
VALUES
  ('WECHAT', '微信', 'wallet', 'WECHAT', 'income,expense,refund,transfer', '💬', 10, 1, '微信收款/微信支付，由事件类型区分'),
  ('ALIPAY', '支付宝', 'wallet', 'ALIPAY', 'income,expense,refund,transfer', '🔵', 20, 1, '支付宝收款/支付宝支付，由事件类型区分'),
  ('BANK_TRANSFER', '银行转账', 'bank', NULL, 'income,expense,refund,transfer', '🏦', 30, 1, ''),
  ('BANK_CARD', '银行卡刷卡', 'card', NULL, 'expense,refund', '💳', 40, 1, ''),
  ('CREDIT_CARD', '信用卡', 'card', NULL, 'expense,refund', '💳', 50, 1, ''),
  ('UNIONPAY', '云闪付', 'card', 'UNIONPAY', 'income,expense,refund', '☁️', 60, 1, ''),
  ('CASH', '现金', 'cash', 'CASH', 'income,expense,refund', '💵', 70, 1, ''),
  ('AUTO_DEBIT', '自动扣款', 'system', NULL, 'expense', '🔁', 80, 1, ''),
  ('REFUND_ORIGINAL', '原路退回', 'system', NULL, 'refund', '↩️', 90, 1, ''),
  ('POS', 'POS刷卡', 'card', NULL, 'expense,refund', '🏧', 100, 1, ''),
  ('QR_CODE', '收付款码', 'wallet', NULL, 'income,expense,refund', '🔳', 110, 1, ''),
  ('OTHER', '其他', 'other', 'OTHER', 'income,expense,refund,transfer', '🔗', 999, 1, '')
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `channel_type` = VALUES(`channel_type`),
  `provider_code` = VALUES(`provider_code`),
  `supported_event_types` = VALUES(`supported_event_types`),
  `icon` = VALUES(`icon`),
  `sort` = VALUES(`sort`),
  `enabled` = VALUES(`enabled`),
  `remark` = VALUES(`remark`);

INSERT INTO `category_template` (`template_code`, `name`, `type`, `icon`, `color`, `sort`, `enabled`)
VALUES
  -- 支出分类模板
  ('FOOD_DINING', '餐饮', 'expense', '🍽️', '#f97316', 10, 1),
  ('FOOD_DELIVERY', '外卖', 'expense', '🥡', '#fb923c', 15, 1),
  ('GROCERY', '买菜日用', 'expense', '🛒', '#f97316', 20, 1),
  ('SNACKS', '零食饮料', 'expense', '🍪', '#fbbf24', 25, 1),
  ('TRANSPORT_TAXI', '打车', 'expense', '🚕', '#3b82f6', 30, 1),
  ('TRANSPORT_PUBLIC', '公共交通', 'expense', '🚌', '#0ea5e9', 35, 1),
  ('GASOLINE', '油费', 'expense', '⛽', '#ef4444', 40, 1),
  ('HOUSING_RENT', '房租', 'expense', '🏠', '#8b5cf6', 50, 1),
  ('HOUSING_MORTGAGE', '房贷', 'expense', '🏘️', '#a855f7', 55, 1),
  ('UTILITIES_ELECTRIC', '电费', 'expense', '⚡', '#facc15', 60, 1),
  ('UTILITIES_WATER', '水费', 'expense', '💧', '#06b6d4', 65, 1),
  ('UTILITIES_GAS', '燃气', 'expense', '🔥', '#ef4444', 70, 1),
  ('UTILITIES_INTERNET', '网费', 'expense', '📡', '#14b8a6', 75, 1),
  ('SHOPPING_CLOTHES', '衣服', 'expense', '👕', '#ec4899', 80, 1),
  ('SHOPPING_ELECTRONICS', '电子产品', 'expense', '📱', '#6366f1', 85, 1),
  ('SHOPPING_DAILY', '日用品', 'expense', '🧴', '#78716c', 90, 1),
  ('ENTERTAINMENT_MOVIE', '电影', 'expense', '🎬', '#8b5cf6', 100, 1),
  ('ENTERTAINMENT_GAME', '游戏', 'expense', '🎮', '#22c55e', 105, 1),
  ('ENTERTAINMENT_SUBSCRIPTION', '会员订阅', 'expense', '📺', '#f97316', 110, 1),
  ('HEALTHCARE_MEDICAL', '就医', 'expense', '🏥', '#ef4444', 120, 1),
  ('HEALTHCARE_MEDICINE', '药品', 'expense', '💊', '#14b8a6', 125, 1),
  ('HEALTHCARE_FITNESS', '健身', 'expense', '💪', '#22c55e', 130, 1),
  ('EDUCATION_BOOKS', '书籍', 'expense', '📚', '#f59e0b', 140, 1),
  ('EDUCATION_COURSE', '课程', 'expense', '🎓', '#3b82f6', 145, 1),
  ('TRAVEL_TRANSPORT', '旅行交通', 'expense', '✈️', '#06b6d4', 150, 1),
  ('TRAVEL_ACCOMMODATION', '旅行住宿', 'expense', '🏨', '#2f7d5c', 155, 1),
  ('TRAVEL_FOOD', '旅行餐饮', 'expense', '🍜', '#f97316', 160, 1),
  ('PERSONAL_CARE', '美容护理', 'expense', '💅', '#ec4899', 170, 1),
  ('PETS', '宠物', 'expense', '🐾', '#f97316', 180, 1),
  ('FAMILY_CHILDREN', '育儿', 'expense', '👨‍👩‍👧', '#f59e0b', 190, 1),

  -- 收入分类模板
  ('SALARY_MONTHLY', '月工资', 'income', '💰', '#ef4444', 1010, 1),
  ('BONUS_YEARLY', '年终奖', 'income', '🎁', '#f59e0b', 1020, 1),
  ('BONUS_PROJECT', '项目奖金', 'income', '🏆', '#22c55e', 1025, 1),
  ('INVESTMENT_DIVIDEND', '股息', 'income', '📈', '#22c55e', 1030, 1),
  ('INVESTMENT_INTEREST', '利息', 'income', '💵', '#14b8a6', 1035, 1),
  ('REIMBURSEMENT_WORK', '工作报销', 'income', '🧾', '#3b82f6', 1040, 1),
  ('SIDE_INCOME', '副业收入', 'income', '💼', '#8b5cf6', 1050, 1),
  ('RENTAL_INCOME', '租金收入', 'income', '🏠', '#2f7d5c', 1060, 1),
  ('GIFT_INCOME', '红包礼金', 'income', '🧧', '#ef4444', 1070, 1),
  ('REFUND_INCOME', '退款', 'income', '↩️', '#64748b', 1080, 1)
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `type` = VALUES(`type`),
  `icon` = VALUES(`icon`),
  `color` = VALUES(`color`),
  `sort` = VALUES(`sort`),
  `enabled` = VALUES(`enabled`);
