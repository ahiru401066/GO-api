CREATE TABLE IF NOT EXISTS `access_logs` (
  `id` INT AUTO_INCREMENT NOT NULL,
  `postal_code` VARCHAR(8) NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);
