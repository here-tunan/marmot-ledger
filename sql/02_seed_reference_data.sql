-- Marmot Ledger reference data seed.
-- This file is idempotent and safe to rerun.

USE `marmot_ledger`;
SET NAMES utf8mb4;

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

INSERT INTO `account_template` (`provider_code`, `name`, `type`, `icon`, `sort`, `enabled`)
VALUES
  ('WECHAT', '微信', 'wallet', '', 10, 1),
  ('ALIPAY', '支付宝', 'wallet', '', 20, 1),
  ('CMB', '招商银行', 'bank', '', 30, 1),
  ('ICBC', '工商银行', 'bank', '', 40, 1),
  ('CCB', '建设银行', 'bank', '', 50, 1),
  ('ABC', '农业银行', 'bank', '', 60, 1),
  ('BOC', '中国银行', 'bank', '', 70, 1),
  ('CASH', '现金', 'cash', '', 80, 1),
  ('IBKR', 'IBKR', 'investment', '', 90, 1),
  ('FUTU', '富途', 'investment', '', 100, 1),
  ('OTHER', '其他', 'other', '', 999, 1)
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `type` = VALUES(`type`),
  `icon` = VALUES(`icon`),
  `sort` = VALUES(`sort`),
  `enabled` = VALUES(`enabled`);

INSERT INTO `category_group` (`group_code`, `name`, `type`, `icon`, `color`, `sort`, `enabled`)
VALUES
  ('FOOD', '餐饮', 'expense', '', '#10b981', 10, 1),
  ('GROCERY', '买菜日用', 'expense', '', '#10b981', 20, 1),
  ('TRANSPORT', '交通', 'expense', '', '#10b981', 30, 1),
  ('HOUSING', '住房', 'expense', '', '#10b981', 40, 1),
  ('UTILITIES', '水电燃气', 'expense', '', '#10b981', 50, 1),
  ('SHOPPING', '购物', 'expense', '', '#10b981', 60, 1),
  ('ENTERTAINMENT', '娱乐', 'expense', '', '#10b981', 70, 1),
  ('HEALTHCARE', '医疗健康', 'expense', '', '#10b981', 80, 1),
  ('EDUCATION', '教育学习', 'expense', '', '#10b981', 90, 1),
  ('TRAVEL', '旅行', 'expense', '', '#10b981', 100, 1),
  ('OTHER_EXPENSE', '其他支出', 'expense', '', '#10b981', 999, 1),
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
  ('WECHAT', '微信', 'wallet', 'WECHAT', 'income,expense', '', 10, 1, '微信收款/微信支付，由事件类型区分'),
  ('ALIPAY', '支付宝', 'wallet', 'ALIPAY', 'income,expense', '', 20, 1, '支付宝收款/支付宝支付，由事件类型区分'),
  ('BANK_TRANSFER', '银行转账', 'bank', NULL, 'income,expense,transfer', '', 30, 1, ''),
  ('BANK_CARD', '银行卡刷卡', 'card', NULL, 'expense', '', 40, 1, ''),
  ('CREDIT_CARD', '信用卡', 'card', NULL, 'expense', '', 50, 1, ''),
  ('UNIONPAY', '云闪付', 'card', 'UNIONPAY', 'income,expense', '', 60, 1, ''),
  ('CASH', '现金', 'cash', 'CASH', 'income,expense', '', 70, 1, ''),
  ('AUTO_DEBIT', '自动扣款', 'system', NULL, 'expense', '', 80, 1, ''),
  ('REFUND_ORIGINAL', '原路退回', 'system', NULL, 'refund', '', 90, 1, ''),
  ('POS', 'POS刷卡', 'card', NULL, 'expense', '', 100, 1, ''),
  ('QR_CODE', '收付款码', 'wallet', NULL, 'income,expense', '', 110, 1, ''),
  ('OTHER', '其他', 'other', 'OTHER', 'income,expense', '', 999, 1, '')
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `channel_type` = VALUES(`channel_type`),
  `provider_code` = VALUES(`provider_code`),
  `supported_event_types` = VALUES(`supported_event_types`),
  `icon` = VALUES(`icon`),
  `sort` = VALUES(`sort`),
  `enabled` = VALUES(`enabled`),
  `remark` = VALUES(`remark`);
