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
