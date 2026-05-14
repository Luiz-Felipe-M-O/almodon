create table if not exists `Users` (
	`uuid`     blob primary key,
	`siape`    text unique not null,
	`name`     text not null,
	`email`    text not null,
	`password` blob not null,
	`role`     text not null,
	`created`  datetime not null,
	`updated`  datetime not null
);

create unique index if not exists `Users_siape` on `Users`(`siape`);

create view if not exists `Users_View` as
	select
        u.`uuid`,
        u.`siape`,
        u.`name`,
        u.`email`,
        u.`password`,
        iif(p.`uuid` is null, u.`role`, 'promoted-admin') as 'role',
        s.`token` is not null as 'logged',
        u.`created`,
        u.`updated`
	from
        Users u
	    left join `Sessions` s on s.`user` = u.`uuid`
	    left join `Promotions` p on p.`user` = u.`uuid`;
