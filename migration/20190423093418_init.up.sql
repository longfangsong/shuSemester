create table Semester
(
    id        bigserial    not null,
    name      varchar(128) not null,
    dateRange daterange    not null,
    EXCLUDE USING gist (dateRange WITH &&)
);

create unique index Semester_dateRange_uindex
    on Semester (dateRange);

create unique index Semester_id_uindex
    on Semester (id);

create unique index Semester_name_uindex
    on Semester (name);

alter table Semester
    add constraint Semester_pk
        primary key (id);

create table Holiday
(
    id        bigserial    not null
        constraint holiday_pk
            primary key,
    name      varchar(128) not null,
    belongTo  bigint
        constraint holiday_semester_id_fk
            references semester
            on update cascade on delete cascade,
    dateRange daterange    not null
);

create unique index holiday_daterange_uindex
    on Holiday (dateRange);

create unique index holiday_id_uindex
    on Holiday (id);

create table Shift
(
    id          bigserial not null,
    fromHoliday bigint    not null
        constraint Shift_holiday_id_fk
            references Holiday
            on update cascade on delete cascade,
    restDate    date      not null,
    workDate    date      not null
);

create unique index Shift_id_uindex
    on Shift (id);

create unique index Shift_restDate_uindex
    on Shift (restDate);

create unique index Shift_workDate_uindex
    on Shift (workDate);

alter table Shift
    add constraint Shift_pk
        primary key (id);


create table Token
(
    id        bigserial    not null,
    tokenHash varchar(256) not null
);

create unique index Token_id_uindex
    on Token (id);

create unique index Token_tokenHash_uindex
    on Token (tokenHash);

alter table Token
    add constraint Token_pk
        primary key (id);

