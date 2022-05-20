ALTER TABLE `auth_users`
	ADD COLUMN `mpesa_renewal` TINYINT(1) NOT NULL DEFAULT '0' AFTER `sacco_id`;