-- Simplify account model and make account templates creation defaults.
-- Safe to rerun on MySQL 8+.

USE `marmot_ledger`;
SET NAMES utf8mb4;

SET @schema_name = DATABASE();

-- Drop account -> account_template linkage first.
SET @has_fk_account_template = (
  SELECT COUNT(*)
  FROM information_schema.TABLE_CONSTRAINTS
  WHERE CONSTRAINT_SCHEMA = @schema_name
    AND TABLE_NAME = 'account'
    AND CONSTRAINT_NAME = 'fk_account_template'
    AND CONSTRAINT_TYPE = 'FOREIGN KEY'
);
SET @sql = IF(@has_fk_account_template > 0,
  'ALTER TABLE `account` DROP FOREIGN KEY `fk_account_template`',
  'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @has_idx_account_standard_account_id = (
  SELECT COUNT(*)
  FROM information_schema.STATISTICS
  WHERE TABLE_SCHEMA = @schema_name
    AND TABLE_NAME = 'account'
    AND INDEX_NAME = 'idx_account_standard_account_id'
);
SET @sql = IF(@has_idx_account_standard_account_id > 0,
  'ALTER TABLE `account` DROP INDEX `idx_account_standard_account_id`',
  'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @has_idx_account_provider_code = (
  SELECT COUNT(*)
  FROM information_schema.STATISTICS
  WHERE TABLE_SCHEMA = @schema_name
    AND TABLE_NAME = 'account'
    AND INDEX_NAME = 'idx_account_provider_code'
);
SET @sql = IF(@has_idx_account_provider_code > 0,
  'ALTER TABLE `account` DROP INDEX `idx_account_provider_code`',
  'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Preserve a meaningful visible name before dropping display_name.
SET @has_display_name = (
  SELECT COUNT(*)
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = @schema_name
    AND TABLE_NAME = 'account'
    AND COLUMN_NAME = 'display_name'
);
SET @sql = IF(@has_display_name > 0,
  'UPDATE `account` SET `name` = `display_name` WHERE `display_name` IS NOT NULL AND `display_name` <> ''''',
  'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @has_provider_code = (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'account' AND COLUMN_NAME = 'provider_code'
);
SET @sql = IF(@has_provider_code > 0, 'ALTER TABLE `account` DROP COLUMN `provider_code`', 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @has_account_group_key = (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'account' AND COLUMN_NAME = 'account_group_key'
);
SET @sql = IF(@has_account_group_key > 0, 'ALTER TABLE `account` DROP COLUMN `account_group_key`', 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @has_standard_account_id = (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'account' AND COLUMN_NAME = 'standard_account_id'
);
SET @sql = IF(@has_standard_account_id > 0, 'ALTER TABLE `account` DROP COLUMN `standard_account_id`', 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @has_bank_code = (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'account' AND COLUMN_NAME = 'bank_code'
);
SET @sql = IF(@has_bank_code > 0, 'ALTER TABLE `account` DROP COLUMN `bank_code`', 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @has_display_name = (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'account' AND COLUMN_NAME = 'display_name'
);
SET @sql = IF(@has_display_name > 0, 'ALTER TABLE `account` DROP COLUMN `display_name`', 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

ALTER TABLE `account`
  MODIFY COLUMN `icon` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Element Plus 图标 key';

SET @has_account_template_color = (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'account_template' AND COLUMN_NAME = 'color'
);
SET @sql = IF(@has_account_template_color = 0,
  'ALTER TABLE `account_template` ADD COLUMN `color` VARCHAR(32) NOT NULL DEFAULT '''' COMMENT ''颜色'' AFTER `icon`',
  'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

ALTER TABLE `account_template`
  MODIFY COLUMN `provider_code` VARCHAR(64) NOT NULL COMMENT '模板代码',
  MODIFY COLUMN `icon` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Element Plus 图标 key';

UPDATE `account_template`
SET
  `type` = LOWER(`type`),
  `icon` = CASE
    WHEN `icon` <> '' THEN `icon`
    WHEN LOWER(`type`) = 'cash' THEN 'Money'
    WHEN LOWER(`type`) = 'wallet' THEN 'Wallet'
    WHEN LOWER(`type`) IN ('bank', 'credit') THEN 'CreditCard'
    WHEN LOWER(`type`) = 'investment' THEN 'TrendCharts'
    WHEN LOWER(`type`) = 'liability' THEN 'Warning'
    ELSE 'Wallet'
  END,
  `color` = CASE
    WHEN `color` <> '' THEN `color`
    WHEN LOWER(`type`) = 'cash' THEN '#f59e0b'
    WHEN LOWER(`type`) = 'wallet' THEN '#22c55e'
    WHEN LOWER(`type`) = 'bank' THEN '#3b82f6'
    WHEN LOWER(`type`) = 'credit' THEN '#ef4444'
    WHEN LOWER(`type`) = 'investment' THEN '#1f2933'
    WHEN LOWER(`type`) = 'liability' THEN '#f97316'
    ELSE '#2f7d5c'
  END;

UPDATE `account`
SET
  `type` = LOWER(`type`),
  `icon` = CASE
    WHEN `icon` <> '' THEN `icon`
    WHEN LOWER(`type`) = 'cash' THEN 'Money'
    WHEN LOWER(`type`) = 'wallet' THEN 'Wallet'
    WHEN LOWER(`type`) IN ('bank', 'credit') THEN 'CreditCard'
    WHEN LOWER(`type`) = 'investment' THEN 'TrendCharts'
    WHEN LOWER(`type`) = 'liability' THEN 'Warning'
    ELSE 'Wallet'
  END,
  `color` = CASE
    WHEN `color` <> '' THEN `color`
    WHEN LOWER(`type`) = 'cash' THEN '#f59e0b'
    WHEN LOWER(`type`) = 'wallet' THEN '#22c55e'
    WHEN LOWER(`type`) = 'bank' THEN '#3b82f6'
    WHEN LOWER(`type`) = 'credit' THEN '#ef4444'
    WHEN LOWER(`type`) = 'investment' THEN '#1f2933'
    WHEN LOWER(`type`) = 'liability' THEN '#f97316'
    ELSE '#2f7d5c'
  END;
