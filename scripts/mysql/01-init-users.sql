-- When going out of alpha, remove drop database
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

CREATE TABLE IF NOT EXISTS `user_topics` (
    `id` varchar(255) NOT NULL,
    `user_id` varchar(255) NOT NULL,
    `topic_id` varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
    FOREIGN KEY (`topic_id`) REFERENCES `notifications`.`topics` (`id`)
);
