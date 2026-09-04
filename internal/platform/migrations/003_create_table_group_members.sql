-- +migrate Up

CREATE TABLE group_members (
                               group_id UUID NOT NULL,
                               user_id UUID NOT NULL,

                               joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                               CONSTRAINT group_members_pk
                                   PRIMARY KEY (group_id, user_id),

                               CONSTRAINT group_members_group_fk
                                   FOREIGN KEY (group_id)
                                       REFERENCES groups(id)
                                       ON DELETE CASCADE,

                               CONSTRAINT group_members_user_fk
                                   FOREIGN KEY (user_id)
                                       REFERENCES users(id)
);

CREATE INDEX group_members_user_id_idx
    ON group_members(user_id);


-- +migrate Down

DROP TABLE IF EXISTS group_members;