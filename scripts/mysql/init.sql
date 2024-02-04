-- When going out of alpha, remove drop database
DROP DATABASE IF EXISTS `users`;
DROP DATABASE IF EXISTS `notifications`;

CREATE DATABASE IF NOT EXISTS `users` DEFAULT CHARACTER SET utf8mb4;

USE `users`;

CREATE TABLE IF NOT EXISTS `users` (
    `id` varchar(255) NOT NULL,
    `username` varchar(255) NOT NULL,
    `password` varchar(255) NOT NULL,
    `profile` varchar(65535) NOT NULL,
    `is_enabled` boolean NOT NULL DEFAULT true,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp,
    PRIMARY KEY (`id`)
)

CREATE DATABASE IF NOT EXISTS `notifications` DEFAULT CHARACTER SET utf8mb4;

USE `notifications`;

CREATE TABLE IF NOT EXISTS `channels` (
    `id` varchar(255) NOT NULL,
    `name` varchar(255) NOT NULL,
    `is_enabled` boolean NOT NULL DEFAULT false,
    PRIMARY KEY (`id`)
)

INSERT INTO `channels` (`id`, `name`, `is_enabled`) VALUES ('1', 'email', true);
INSERT INTO `channels` (`id`, `name`, `is_enabled`) VALUES ('2', 'sms', true);
INSERT INTO `channels` (`id`, `name`, `is_enabled`) VALUES ('3', 'slack', true);

CREATE TABLE IF NOT EXISTS `user_channels` (
    `user_id` varchar(255) NOT NULL,
    `channel_id` varchar(255) NOT NULL,
    `is_enabled` boolean NOT NULL DEFAULT false,
    PRIMARY KEY (`user_id`, `channel_id`),
)

CREATE TABLE IF NOT EXISTS `topics` (
    `id` varchar(255) NOT NULL,
    `name` varchar(255) NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp,
    PRIMARY KEY (`id`)
)

CREATE TABLE IF NOT EXISTS `user_topics` (
    `id` varchar(255) NOT NULL,
    `user_id` varchar(255) NOT NULL,
    `topic_id` varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_topics_users_user_id`) REFERENCES `users.users` (`id`),
    FOREIGN KEY (`user_topics_topics_topic_id`) REFERENCES `topics` (`id`)
)

CREATE TABLE IF NOT EXISTS `messages` (
    `id` varchar(255) NOT NULL,
    `event_id` varchar(255) NOT NULL,
    `title` varchar(255),
    `message` varchar(65535) NOT NULL,
    `template` varchar(255),
    `status` int(11) NOT NULL,
    `channel_id` varchar(255) NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`messages_channels_channel_id`) REFERENCES `channels` (`id`),
    INDEX (`event_id`)
)

CREATE TABLE IF NOT EXISTS `message_topics` (
    `id` varchar(255) NOT NULL,
    `message_id` varchar(255) NOT NULL,
    `topic_id` varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`message_topics_messages_message_id`) REFERENCES `messages` (`id`),
    FOREIGN KEY (`message_topics_topics_topic_id`) REFERENCES `topics` (`id`)
)

CREATE TABLE IF NOT EXISTS `message_users` (
    `id` varchar(255) NOT NULL,
    `message_id` varchar(255) NOT NULL,
    `user_id` varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`message_users_messages_message_id`) REFERENCES `messages` (`id`),
    FOREIGN KEY (`nmessage_users_users_user_id`) REFERENCES `users.users` (`id`)
)