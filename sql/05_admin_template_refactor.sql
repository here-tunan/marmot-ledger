USE `marmot_ledger`;
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 1;

-- ============================================
-- 1. 添加用户角色字段
-- ============================================
ALTER TABLE `user` ADD COLUMN `role` VARCHAR(32) NOT NULL DEFAULT 'user' COMMENT '用户角色: user/admin' AFTER `name`;
ALTER TABLE `user` ADD KEY `idx_user_role` (`role`);

-- 设置开发环境的管理员用户
UPDATE `user` SET `role` = 'admin' WHERE `account` = 'admin';

-- ============================================
-- 2. 分类模板表
-- ============================================
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

-- ============================================
-- 3. 家庭分类组表（家庭维度自定义分组）
-- ============================================
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

-- ============================================
-- 4. 家庭分类组成员表（多对多关联）
-- ============================================
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

-- ============================================
-- 5. 预置分类模板数据
-- ============================================
INSERT INTO `category_template` (`template_code`, `name`, `type`, `sort`) VALUES
-- 收入类
('SALARY', '工资', 'income', 1),
('BONUS', '奖金', 'income', 2),
('INVESTMENT_INCOME', '投资收益', 'income', 3),
('SIDE_INCOME', '副业收入', 'income', 4),
('REFUND', '退款', 'income', 5),
('OTHER_INCOME', '其他收入', 'income', 10),
-- 支出类
('FOOD', '餐饮', 'expense', 1),
('TRANSPORT', '交通', 'expense', 2),
('SHOPPING', '购物', 'expense', 3),
('ENTERTAINMENT', '娱乐', 'expense', 4),
('HOUSING', '居住', 'expense', 5),
('UTILITIES', '水电燃气', 'expense', 6),
('HEALTHCARE', '医疗健康', 'expense', 7),
('EDUCATION', '教育学习', 'expense', 8),
('FAMILY', '家人开支', 'expense', 9),
('TRAVEL', '旅行', 'expense', 10),
('OTHER_EXPENSE', '其他支出', 'expense', 20)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- ============================================
-- 6. 补充 channel_template 预置数据
-- ============================================
INSERT INTO `channel_template` (`channel_code`, `name`, `channel_type`, `provider_code`, `supported_event_types`, `icon`, `sort`, `remark`) VALUES
('WECHAT', '微信', 'wallet', 'WECHAT', 'income,expense,refund,transfer', '💬', 1, '微信收款/微信支付，由事件类型区分'),
('ALIPAY', '支付宝', 'wallet', 'ALIPAY', 'income,expense,refund,transfer', '🔵', 2, '支付宝收款/支付宝支付，由事件类型区分'),
('CASH', '现金', 'cash', 'CASH', 'income,expense,refund', '💵', 3, ''),
('BANK_TRANSFER', '银行转账', 'bank', 'BANK', 'transfer,income,expense,refund', '🏦', 4, '')
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `channel_type` = VALUES(`channel_type`),
  `provider_code` = VALUES(`provider_code`),
  `supported_event_types` = VALUES(`supported_event_types`),
  `icon` = VALUES(`icon`),
  `sort` = VALUES(`sort`),
  `remark` = VALUES(`remark`);

-- ============================================
-- 7. 补充 account_template 预置数据
-- ============================================
INSERT INTO `account_template` (`provider_code`, `name`, `type`, `icon`, `color`, `sort`) VALUES
('WECHAT', '微信钱包', 'wallet', 'Wallet', '#22c55e', 1),
('ALIPAY', '支付宝', 'wallet', 'Wallet', '#1677ff', 2),
('ICBC', '中国工商银行', 'bank', 'CreditCard', '#f97316', 3),
('CCB', '中国建设银行', 'bank', 'CreditCard', '#1d4ed8', 4),
('ABC', '中国农业银行', 'bank', 'CreditCard', '#16a34a', 5),
('BOC', '中国银行', 'bank', 'CreditCard', '#dc2626', 6),
('CMB', '招商银行', 'bank', 'CreditCard', '#ef4444', 7),
('SPDB', '浦发银行', 'bank', 'CreditCard', '#0ea5e9', 8),
('CIB', '兴业银行', 'bank', 'CreditCard', '#2563eb', 9),
('BOCOM', '交通银行', 'bank', 'CreditCard', '#1d4ed8', 10),
('CMBC', '民生银行', 'bank', 'CreditCard', '#16a34a', 11),
('PAB', '平安银行', 'bank', 'CreditCard', '#f97316', 12),
('CITIC', '中信银行', 'bank', 'CreditCard', '#dc2626', 13),
('GDB', '广发银行', 'bank', 'CreditCard', '#ef4444', 14),
('PSBC', '中国邮政储蓄银行', 'bank', 'CreditCard', '#16a34a', 15)
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `type` = VALUES(`type`),
  `icon` = VALUES(`icon`),
  `color` = VALUES(`color`),
  `sort` = VALUES(`sort`);
