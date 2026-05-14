create table if not exists `Sessions` (
	`token`   blob     primary key,
	`user`    blob     not null,
	`renewed` int      not null,
	`expires` datetime not null,
	`created` datetime not null,

	foreign key (`user`) references `Users`(`uuid`)
);

create unique index if not exists `Sessions_user` on `Sessions`(`user`);
