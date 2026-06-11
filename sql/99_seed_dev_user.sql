-- Development-only user seed.
-- Do not execute this file in production.
-- The current login API compares account/password directly, so this seed keeps local smoke testing simple.

USE `marmot_ledger`;
SET NAMES utf8mb4;

INSERT INTO `user` (`account`, `password`, `name`, `desc`, `avatar`, `extra`, `is_deleted`)
VALUES ('admin', 'admin', '管理员', 'Marmot Ledger development user', '', '', 0)
ON DUPLICATE KEY UPDATE
  `password` = VALUES(`password`),
  `name` = VALUES(`name`),
  `desc` = VALUES(`desc`),
  `avatar` = VALUES(`avatar`),
  `extra` = VALUES(`extra`),
  `is_deleted` = VALUES(`is_deleted`);
