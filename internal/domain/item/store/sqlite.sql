create table if not exists `Items_true` (
    `uuid`      blob     not null,
    `version`   int      not null,
    `material`  blob     not null,
    `amount`    real     not null,
    `unit_cost` int      not null,
    `expires`   datetime not null,
    `created`   datetime not null,

    primary key (`uuid`, `version`),
    foreign key (`material`) references `Materials`(`uuid`)
);

create table if not exists `Items` (
    `uuid`    blob     primary key,
    `version` int      not null,
    `created` datetime not null,

    foreign key (`uuid`, `version`) references `Items_true`(`uuid`, `version`)
);

create view if not exists `Items_View` as
    select
        i.`uuid`,
        i.`version`,
        m.`name`,
        m.`ecampus`,
        m.`catmat`,
        m.`siads`,
        t.`material`,
        t.`amount`,
        t.`unit_cost`,
        m.`unit`,
        t.`expires`,
        m.`min`,
        i.`created`,
        t.`created` as 'updated'
    from
        `Items` i
        join `Items_true` t on i.`uuid` = t.`uuid` and i.`version` = t.`version`
        join `Materials` m on t.`material` = m.`uuid`;

create view if not exists `Items_History_View` as
    select
        `version`,
        `uuid`,
        `material`,
        `amount`,
        `unit_cost`,
        `expires`,
        `created`
    from
        `Items_true`
    order by
        `version` desc;
