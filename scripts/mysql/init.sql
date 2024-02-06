-- When going out of alpha, remove drop database
DROP DATABASE IF EXISTS `users`;
DROP DATABASE IF EXISTS `notifications`;

CREATE DATABASE IF NOT EXISTS `users` DEFAULT CHARACTER SET utf8mb4;

USE `users`;

CREATE TABLE IF NOT EXISTS `users` (
    `id` varchar(255) NOT NULL,
    `username` varchar(255) NOT NULL,
    `password` varchar(255) NOT NULL,
    `email` varchar(255) NOT NULL,
    `profile` varchar(2096) NOT NULL,
    `is_enabled` boolean NOT NULL DEFAULT true,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp,
    PRIMARY KEY (`id`)
);

CREATE DATABASE IF NOT EXISTS `notifications` DEFAULT CHARACTER SET utf8mb4;

USE `notifications`;

CREATE TABLE IF NOT EXISTS `channels` (
    `id` varchar(255) NOT NULL,
    `name` varchar(255) NOT NULL,
    `is_enabled` boolean NOT NULL DEFAULT false,
    PRIMARY KEY (`id`)
);

INSERT INTO `channels` (`id`, `name`, `is_enabled`) VALUES ('1', 'email', true);
INSERT INTO `channels` (`id`, `name`, `is_enabled`) VALUES ('2', 'sms', false);
INSERT INTO `channels` (`id`, `name`, `is_enabled`) VALUES ('3', 'slack', false);

CREATE TABLE IF NOT EXISTS `user_channels` (
    `id` varchar(255) NOT NULL,
    `user_id` varchar(255) NOT NULL,
    `channel_id` varchar(255) NOT NULL,
    `is_enabled` boolean NOT NULL DEFAULT false,
    PRIMARY KEY (`user_id`, `channel_id`),
    FOREIGN KEY (`user_id`) REFERENCES `users`.`users` (`id`),
    FOREIGN KEY (`channel_id`) REFERENCES `channels` (`id`)
);

CREATE TABLE IF NOT EXISTS `topics` (
    `id` varchar(255) NOT NULL,
    `name` varchar(255) NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `messages` (
    `id` varchar(255) NOT NULL,
    `event_id` varchar(255) NOT NULL,
    `subject` varchar(255),
    `message` varchar(8192) NOT NULL,
    `template_id` varchar(255),
    `status` int(11) NOT NULL,
    `channel_id` varchar(255) NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp,
    `version` int(11) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`channel_id`) REFERENCES `channels` (`id`),
    INDEX `ix_event_id` (`event_id`),
    UNIQUE `uq_event_id` (`event_id`)
);

USE `users`;

CREATE TABLE IF NOT EXISTS `user_topics` (
    `id` varchar(255) NOT NULL,
    `user_id` varchar(255) NOT NULL,
    `topic_id` varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
    FOREIGN KEY (`topic_id`) REFERENCES `notifications`.`topics` (`id`)
);
