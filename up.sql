drop table users;
drop table sessions;
drop table practicecontents;
drop table tags;


create table users (
    id  serial PRIMARY KEY ,
	uuid varchar(64) NOT NULL UNIQUE,
	name varchar(255),
	email varchar(255) NOT NULL UNIQUE,
	password varchar(255) NOT NULL,
	created_at timestamp
);


create table sessions (
    id serial PRIMARY KEY ,
	uuid varchar(64) NOT NULL UNIQUE,
	name varchar(255) NOT NULL,
	email varchar(255),
	user_id INTEGER,
	created_at timestamp
);


create table practicecontents (
    id serial PRIMARY KEY ,
	user_id INTEGER NOT NULL,
	prefecture varchar(255) NOT NULL,
	place varchar(255) NOT NULL,
	strat_time varchar(255)  NOT NULL,
	end_time varchar(255)  NOT NULL,
	scale INTEGER ,
	tags varchar(255),
	describe TEXT,
	uuid varchar(64) NOT NULL UNIQUE,
	created_at timestamp

);

create table tags (
    id INTEGER ,
    tag varchar(255)
);
