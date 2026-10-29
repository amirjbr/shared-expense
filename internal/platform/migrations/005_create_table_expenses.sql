-- +migrate Up

CREATE TABLE expenses (
                          id UUID PRIMARY KEY,
                          group_id UUID NOT NULL,
                          title VARCHAR(150) NOT NULL,
                          created_by_user_id UUID NOT NULL,
                          total_amount BIGINT NOT NULL,
                          split_type VARCHAR(20) NOT NULL,
                          occurred_at TIMESTAMPTZ NOT NULL,

                          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                          CONSTRAINT expenses_title_not_empty
                              CHECK (LENGTH(TRIM(title)) > 0),

                          CONSTRAINT expenses_total_amount_positive
                              CHECK (total_amount > 0),

                          CONSTRAINT expenses_split_type_valid
                              CHECK (
                                  split_type IN (
                                                 'equal',
                                                 'custom'
                                      )
                                  ),

                          CONSTRAINT expenses_group_fk
                              FOREIGN KEY (group_id)
                                  REFERENCES groups(id)
                                  ON DELETE CASCADE,

                          CONSTRAINT expenses_creator_is_group_member_fk
                              FOREIGN KEY (group_id, created_by_user_id)
                                  REFERENCES group_members(group_id, user_id),

                          CONSTRAINT expenses_id_group_unique
                              UNIQUE (id, group_id)
);

CREATE INDEX expenses_group_occurred_at_idx
    ON expenses(group_id, occurred_at DESC);

CREATE INDEX expenses_group_creator_idx
    ON expenses(group_id, created_by_user_id);


-- +migrate Down

DROP TABLE IF EXISTS expenses;