create table "users" (
  id bigserial primary key,
  created_at timestamptz default now(),
  email varchar not null unique,
  password varchar not null,
  username varchar,
  role varchar,
  verified boolean default false,
  payload jsonb
);

create index idx_user_email on "users" (
  email asc
);

create index idx_user_username on "users" (
  username desc
);

create table verify_accounts (
  id bigserial primary key,
  created_at timestamptz default now(),
  token varchar not null,
  user_id bigint references users(id) on delete cascade on update cascade
);

create table forget_passwords (
  id bigserial primary key,
  created_at timestamptz default now(),
  token varchar not null,
  user_id bigint references users(id) on delete cascade on update cascade
);
