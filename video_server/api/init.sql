create database video;

use video;

create table user (
    id int primary key auto_increment,
    name varchar(10) not null comment '用户名'
)engine=innodb default charset=utf8