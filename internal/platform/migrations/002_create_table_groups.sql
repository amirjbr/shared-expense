-- +migrate Up

CREATE TABLE groups (
                        id UUID PRIMARY KEY,
                        name VARCHAR(100) NOT NULL,
                        owner_user_id UUID NOT NULL,

                        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                        CONSTRAINT groups_name_not_empty
                            CHECK (LENGTH(TRIM(name)) > 0),

                        CONSTRAINT groups_owner_user_fk
                            FOREIGN KEY (owner_user_id)
                                REFERENCES users(id)
);

CREATE INDEX groups_owner_user_id_idx
    ON groups(owner_user_id);


-- +migrate Down

DROP TABLE IF EXISTS groups;