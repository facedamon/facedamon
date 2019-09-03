
create database if not exists `negroni`;

use `negroni`;

create table if not exists `users` (
  `id` int not null auto_increment primary key,
  `username` varchar(60),
  `password` varchar(60),
  `email` varchar(60)
)