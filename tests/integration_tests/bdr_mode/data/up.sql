use `bdr_mode`;

begin;
insert into `t1` values (1, '1'), (3, '3'), (5, '5'), (7, '7'), (9, '9');
commit;

begin;
update `t1` set `name` = '11' where `id` = 1;
delete from `t1` where `id` = 3;
update `t1` set `name` = '55' where `id` = 5;
delete from `t1` where `id` = 7;
commit;

begin;
insert into `t1` values (2, '2'), (4, '4'), (6, '6'), (8, '8'), (10, '10');
rollback;

insert into `t1` values (100, '100'), (300, '300'), (500, '500'), (700, '700'), (900, '900');
