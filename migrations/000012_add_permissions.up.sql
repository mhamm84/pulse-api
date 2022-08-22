CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY,
    code text NOT NULL
);

CREATE TABLE IF NOT EXISTS users_permissions (
  user_id bigserial NOT NULL REFERENCES users(id) ON DELETE CASCADE ,
  permission_id bigserial NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, permission_id)
);

INSERT INTO permissions (code)
VALUES
    ('economic:all'),
    ('admin')
