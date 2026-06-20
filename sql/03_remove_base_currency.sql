-- Remove base-currency reporting fields from ordinary bill statistics model.
-- Marmot Ledger now keeps bill events in their original currency and displays
-- statistics grouped by financial_event.currency without FX conversion.

USE `marmot_ledger`;
SET NAMES utf8mb4;

ALTER TABLE `financial_event`
  DROP FOREIGN KEY `fk_financial_event_base_currency`;

ALTER TABLE `financial_event`
  DROP INDEX `idx_financial_event_base_currency`;

ALTER TABLE `financial_event`
  DROP COLUMN `exchange_rate`,
  DROP COLUMN `base_amount`,
  DROP COLUMN `base_currency`;

ALTER TABLE `family`
  DROP FOREIGN KEY `fk_family_base_currency`;

ALTER TABLE `family`
  DROP INDEX `idx_family_base_currency`;

ALTER TABLE `family`
  DROP COLUMN `base_currency`;
