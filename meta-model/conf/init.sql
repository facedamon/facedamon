drop table if exists Model_Base_Word_Info;
create table Model_Base_Word_Info(
    Unique_Num varchar(20) not null primary key comment '唯一编号',
    Base_Word_Id varchar(20) not null comment '基词编号',
    Base_Word_Cn varchar(100) not null comment '基词中文含义',
    Base_Word_En varchar(32) not null comment '基词英文含义',
    Abbreviation1 varchar(32) not null comment '基词简写1',
    Abbreviation2 varchar(32) not null comment '基词简写2',
    User_Id varchar(32) not null comment '用户编号',
    Creator varchar(100) not null comment '创建人',
    Create_Time Timestamp not null comment '创建时间',
    Modifier varchar(100)  comment '修改人',
    Modify_Time Timestamp comment '修改时间'
) engine=Innodb default charset=utf8 comment '模型基词信息';

insert into model_base_word_info values('0000000001','B000000001','数据','Data','Data','Data','LPM','lusing','2019-09-01','lusing','2019-09-03');
insert into model_base_word_info values('0000000002','B000000002','日期','Dt','Dt','Dt','LPM','lusing','2019-09-01','lusing','2019-09-03');
insert into model_base_word_info values('0000000003','B000000003','付款','Pay','Pay','Pay','LPM','lusing','2019-09-01','lusing','2019-09-03');
insert into model_base_word_info values('0000000004','B000000004','月份','Month','Month','Month','LPM','lusing','2019-09-01','lusing','2019-09-03');
insert into model_base_word_info values('0000000005','B000000005','币种','Currency','Curr','Curr','LPM','lusing','2019-09-01','lusing','2019-09-03');
insert into model_base_word_info values('0000000006','B000000006','编号','Number','Num','Num','LPM','lusing','2019-09-01','lusing','2019-09-03');
insert into model_base_word_info values('0000000007','B000000007','代码','Code','Cd','Cd','LPM','lusing','2019-09-01','lusing','2019-09-03');