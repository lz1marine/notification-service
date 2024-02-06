USE `users`;

INSERT INTO `users` (`id`, `username`, `password`, `email`, `profile`, `is_enabled`) VALUES ('1', 'user1', 'scrambled', 'scrambled', '{"address": "myaddress", "phone": "+35912345", "email": "<YOUR-EMAIL>@gmail.com"}', true);
INSERT INTO `users` (`id`, `username`, `password`, `email`, `profile`, `is_enabled`) VALUES ('2', 'user2', 'scrambled', 'scrambled', '{"address": "myaddress2", "phone": "+35912346", "email": "<YOUR-EMAIL>@gmail.com"}', true);
INSERT INTO `users` (`id`, `username`, `password`, `email`, `profile`, `is_enabled`) VALUES ('3', 'user3', 'scrambled', 'scrambled', '{"address": "myaddress3", "phone": "+3598971234", "email": "<YOUR-EMAIL>@gmail.com"}', true);

USE `notifications`;

INSERT INTO `user_channels` (`id`, `user_id`, `channel_id`, `is_enabled`) VALUES ('1', '1', '1', true);
INSERT INTO `user_channels` (`id`, `user_id`, `channel_id`, `is_enabled`) VALUES ('2', '1', '2', true);
INSERT INTO `user_channels` (`id`, `user_id`, `channel_id`, `is_enabled`) VALUES ('3', '2', '1', true);
INSERT INTO `user_channels` (`id`, `user_id`, `channel_id`, `is_enabled`) VALUES ('4', '3', '1', false);

INSERT INTO `topics` (`id`, `name`) VALUES ('1', 'products');

USE `users`;

INSERT INTO `user_topics` (`id`, `user_id`, `topic_id`) VALUES ('1', '1', '1');
INSERT INTO `user_topics` (`id`, `user_id`, `topic_id`) VALUES ('2', '2', '1');
