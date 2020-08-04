// 创建函数
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

// 创建 users 表
create table users (
   id serial NOT NULL PRIMARY KEY,
   name           varchar(20)    NOT NULL,
   created_at timestamptz not null default now(),
   updated_at timestamptz not null default now()
);

// users 触发器: for each row
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

// 插入数据
INSERT INTO users (name, age) values('mike', 20);

// 用户关系表
create table relationships (
   id serial NOT NULL PRIMARY KEY,
   user_id INT NOT NULL,
   follower_id INT NOT NULL,
   state varchar(20) NOT NULL,
   created_at timestamptz not null default now(),
   updated_at timestamptz not null default now()
);

// relationships 触发器
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON relationships
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

//插入数据
INSERT INTO relationships (user_id, follower_id, state) values(1, 11, 'liked');

//todo index
CREATE INDEX "idx_n" ON "users" ("name");
CREATE INDEX "idx_i" ON "relationships" ("user_id");
CREATE INDEX "idx_u" ON "relationships" ("user_id", "follower_id");
