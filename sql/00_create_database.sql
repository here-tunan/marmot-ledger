-- Marmot Ledger database bootstrap.
-- This database name matches the current backend DSN in internal/infrastructure/mysql.go.

CREATE DATABASE IF NOT EXISTS `marmot_ledger`
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE `marmot_ledger`;
