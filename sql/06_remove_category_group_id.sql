USE `marmot_ledger`;
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 1;

-- 1. 先移除外键约束（如果存在）
ALTER TABLE `category` DROP FOREIGN KEY IF EXISTS `fk_category_category_group`;

-- 2. 删除冗余字段
ALTER TABLE `category` DROP COLUMN IF EXISTS `category_group_id`;

-- 3. 备注：
-- 系统级 category_group 表保留，用于模板和统计口径
-- 分类与分组的关联完全通过 family_category_group_member 多对多表实现
