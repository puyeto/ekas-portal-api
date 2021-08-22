-- phpMyAdmin SQL Dump
-- version 5.0.2
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1:3307
-- Generation Time: Aug 22, 2021 at 06:06 PM
-- Server version: 10.4.13-MariaDB
-- PHP Version: 7.3.21

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `ekas_portal`
--

-- --------------------------------------------------------

--
-- Table structure for table `saccos`
--

DROP TABLE IF EXISTS `saccos`;
CREATE TABLE IF NOT EXISTS `saccos` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(200) NOT NULL,
  `short_name` varchar(20) NOT NULL,
  `address` varchar(200) NOT NULL,
  `created_on` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;

--
-- Dumping data for table `saccos`
--

INSERT INTO `saccos` (`id`, `name`, `short_name`, `address`, `created_on`) VALUES
(1, '2NK SACCO', '2NK', 'Nairobi', '2021-08-22 18:05:22');

INSERT INTO `roles` (`role_id`, `role_name`, `description`) VALUES ('10004', 'sacco', 'sacco');

ALTER TABLE `vehicle_details` ADD `sacco_id` INT NOT NULL DEFAULT '0' AFTER `vehicle_string_id`;

COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
