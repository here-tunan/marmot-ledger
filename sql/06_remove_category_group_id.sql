USE `marmot_ledger`;
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 1;

-- Marmot Ledger schema cleanup:
--   * category groups are reporting mappings, not event/category columns
--   * channel_account_id and bucket_group_key are removed as redundant metadata
--   * templates are copy sources, not persistent category references
--   * unused exchange/investment tables are postponed out of the MVP schema

SET @schema_name = DATABASE();

-- Drop optional tables first, in dependency order.
DROP TABLE IF EXISTS `investment_snapshot`;
DROP TABLE IF EXISTS `investment_asset`;
DROP TABLE IF EXISTS `exchange_detail`;
DROP TABLE IF EXISTS `exchange_rate`;

-- financial_event.category_group_id / channel_account_id cleanup.
SET @sql := IF(
  EXISTS(SELECT 1 FROM information_schema.TABLE_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = @schema_name AND TABLE_NAME = 'financial_event' AND CONSTRAINT_NAME = 'fk_financial_event_category_group'),
  'ALTER TABLE `financial_event` DROP FOREIGN KEY `fk_financial_event_category_group`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  EXISTS(SELECT 1 FROM information_schema.TABLE_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = @schema_name AND TABLE_NAME = 'financial_event' AND CONSTRAINT_NAME = 'fk_financial_event_channel_account'),
  'ALTER TABLE `financial_event` DROP FOREIGN KEY `fk_financial_event_channel_account`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  EXISTS(SELECT 1 FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'financial_event' AND INDEX_NAME = 'idx_financial_event_user_category_group_time'),
  'ALTER TABLE `financial_event` DROP INDEX `idx_financial_event_user_category_group_time`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  NOT EXISTS(SELECT 1 FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'financial_event' AND INDEX_NAME = 'idx_financial_event_user_category_time'),
  'ALTER TABLE `financial_event` ADD KEY `idx_financial_event_user_category_time` (`user_id`, `category_id`, `event_time`)',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'financial_event' AND COLUMN_NAME = 'category_group_id'),
  'ALTER TABLE `financial_event` DROP COLUMN `category_group_id`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'financial_event' AND COLUMN_NAME = 'channel_account_id'),
  'ALTER TABLE `financial_event` DROP COLUMN `channel_account_id`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- category.category_group_id / template_id cleanup, while keeping icon/color.
SET @sql := IF(
  EXISTS(SELECT 1 FROM information_schema.TABLE_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = @schema_name AND TABLE_NAME = 'category' AND CONSTRAINT_NAME = 'fk_category_group'),
  'ALTER TABLE `category` DROP FOREIGN KEY `fk_category_group`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  EXISTS(SELECT 1 FROM information_schema.TABLE_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = @schema_name AND TABLE_NAME = 'category' AND CONSTRAINT_NAME = 'fk_category_category_group'),
  'ALTER TABLE `category` DROP FOREIGN KEY `fk_category_category_group`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  EXISTS(SELECT 1 FROM information_schema.TABLE_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = @schema_name AND TABLE_NAME = 'category' AND CONSTRAINT_NAME = 'fk_category_template'),
  'ALTER TABLE `category` DROP FOREIGN KEY `fk_category_template`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  EXISTS(SELECT 1 FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'category' AND INDEX_NAME = 'idx_category_group_id'),
  'ALTER TABLE `category` DROP INDEX `idx_category_group_id`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  EXISTS(SELECT 1 FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'category' AND INDEX_NAME = 'idx_category_template_id'),
  'ALTER TABLE `category` DROP INDEX `idx_category_template_id`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  NOT EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'category' AND COLUMN_NAME = 'icon'),
  'ALTER TABLE `category` ADD COLUMN `icon` VARCHAR(255) NOT NULL DEFAULT '''' COMMENT ''图标'' AFTER `type`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  NOT EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'category' AND COLUMN_NAME = 'color'),
  'ALTER TABLE `category` ADD COLUMN `color` VARCHAR(32) NOT NULL DEFAULT '''' COMMENT ''颜色'' AFTER `icon`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'category' AND COLUMN_NAME = 'category_group_id'),
  'ALTER TABLE `category` DROP COLUMN `category_group_id`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'category' AND COLUMN_NAME = 'template_id'),
  'ALTER TABLE `category` DROP COLUMN `template_id`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- bucket.bucket_group_key cleanup.
SET @sql := IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'bucket' AND COLUMN_NAME = 'bucket_group_key'),
  'ALTER TABLE `bucket` DROP COLUMN `bucket_group_key`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- category_template no longer stores a default reporting group.
SET @sql := IF(
  EXISTS(SELECT 1 FROM information_schema.TABLE_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = @schema_name AND TABLE_NAME = 'category_template' AND CONSTRAINT_NAME = 'fk_category_template_default_group'),
  'ALTER TABLE `category_template` DROP FOREIGN KEY `fk_category_template_default_group`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'category_template' AND COLUMN_NAME = 'default_category_group_id'),
  'ALTER TABLE `category_template` DROP COLUMN `default_category_group_id`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
