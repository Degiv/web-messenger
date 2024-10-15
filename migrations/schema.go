package migrations

const schema = `
create table users (
    id serial primary key,
    username varchar not null
);

create table messages (
    id serial primary key,
    user_id int references users(id),
    text varchar not null 
);

create table conferences (
    id serial primary key,
    name varchar not null,
    last_message_id int references messages(id) on delete set null 
);

create table UsersConferencesRelation (
    user_id int,
    conference_id int,
    foreign key (user_id) references users(id) on delete cascade,
    foreign key (conference_id) references conferences(id) on delete cascade,
    unique (user_id, conference_id)
)
`
