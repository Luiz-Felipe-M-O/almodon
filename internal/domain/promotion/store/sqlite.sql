create table if not exists `Promotions` (
	`uuid`    blob primary key,
	`user`    blob not null,
	`expires` datetime not null,

	foreign key (`user`) references `Users`(`uuid`)
);

create index if not exists `Promotions_user` on `Promotions`(`user`);