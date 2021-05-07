create table users (
  id         int AUTO_INCREMENT primary key ,
  uuid       varchar(64) not null unique,
  name       varchar(255),
  email      varchar(255) not null unique,
  password   varchar(255) not null,
  created_at DATETIME  not null   
);

create table sessions (
  id         int AUTO_INCREMENT primary key,
  uuid       varchar(64) not null unique,
  email      varchar(255),
  user_id    integer references users(id),
  created_at DATETIME not null   
);