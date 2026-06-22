USE `marmot_ledger`;
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 1;

SET @schema_name = DATABASE();

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

INSERT IGNORE INTO `personal_channel` (`user_id`, `channel_template_id`, `name`, `channel_type`, `provider_code`, `supported_event_types`, `icon`, `sort`, `is_active`, `is_deleted`)
SELECT u.id, ct.id, ct.name, LOWER(ct.channel_type), ct.provider_code, ct.supported_event_types, ct.icon, ct.sort, 1, 0
FROM `user` u
JOIN `channel_template` ct ON ct.enabled = 1
WHERE u.is_deleted = 0;

SET @sql := IF(
  NOT EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'financial_event' AND COLUMN_NAME = 'channel_id'),
  'ALTER TABLE `financial_event` ADD COLUMN `channel_id` BIGINT NULL COMMENT ''个人渠道ID'' AFTER `category_id`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  NOT EXISTS(SELECT 1 FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'financial_event' AND INDEX_NAME = 'idx_financial_event_user_channel_time'),
  'ALTER TABLE `financial_event` ADD KEY `idx_financial_event_user_channel_time` (`user_id`, `channel_id`, `event_time`)',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

UPDATE `financial_event` fe
JOIN `channel_template` ct ON ct.channel_code = fe.channel_type
JOIN `personal_channel` pc ON pc.user_id = fe.user_id AND pc.channel_template_id = ct.id AND pc.is_deleted = 0
SET fe.channel_id = pc.id
WHERE fe.channel_id IS NULL
  AND fe.channel_type IS NOT NULL;

SET @sql := IF(
  EXISTS(SELECT 1 FROM information_schema.TABLE_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = @schema_name AND TABLE_NAME = 'financial_event' AND CONSTRAINT_NAME = 'fk_financial_event_channel'),
  'ALTER TABLE `financial_event` DROP FOREIGN KEY `fk_financial_event_channel`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  EXISTS(SELECT 1 FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'financial_event' AND INDEX_NAME = 'idx_financial_event_channel_time'),
  'ALTER TABLE `financial_event` DROP INDEX `idx_financial_event_channel_time`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = @schema_name AND TABLE_NAME = 'financial_event' AND COLUMN_NAME = 'channel_type'),
  'ALTER TABLE `financial_event` DROP COLUMN `channel_type`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  NOT EXISTS(SELECT 1 FROM information_schema.TABLE_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = @schema_name AND TABLE_NAME = 'financial_event' AND CONSTRAINT_NAME = 'fk_financial_event_channel'),
  'ALTER TABLE `financial_event` ADD CONSTRAINT `fk_financial_event_channel` FOREIGN KEY (`channel_id`) REFERENCES `personal_channel` (`id`)',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
