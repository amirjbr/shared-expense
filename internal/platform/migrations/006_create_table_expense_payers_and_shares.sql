-- +migrate Up

CREATE TABLE expense_payers (
                                expense_id UUID NOT NULL,
                                group_id UUID NOT NULL,
                                user_id UUID NOT NULL,
                                amount BIGINT NOT NULL,

                                CONSTRAINT expense_payers_pk
                                    PRIMARY KEY (expense_id, user_id),

                                CONSTRAINT expense_payers_amount_positive
                                    CHECK (amount > 0),

                                CONSTRAINT expense_payers_expense_fk
                                    FOREIGN KEY (expense_id, group_id)
                                        REFERENCES expenses(id, group_id)
                                        ON DELETE CASCADE,

                                CONSTRAINT expense_payers_group_member_fk
                                    FOREIGN KEY (group_id, user_id)
                                        REFERENCES group_members(group_id, user_id)
);

CREATE INDEX expense_payers_group_user_idx
    ON expense_payers(group_id, user_id);


CREATE TABLE expense_shares (
                                expense_id UUID NOT NULL,
                                group_id UUID NOT NULL,
                                user_id UUID NOT NULL,
                                amount BIGINT NOT NULL,

                                CONSTRAINT expense_shares_pk
                                    PRIMARY KEY (expense_id, user_id),

                                CONSTRAINT expense_shares_amount_positive
                                    CHECK (amount > 0),

                                CONSTRAINT expense_shares_expense_fk
                                    FOREIGN KEY (expense_id, group_id)
                                        REFERENCES expenses(id, group_id)
                                        ON DELETE CASCADE,

                                CONSTRAINT expense_shares_group_member_fk
                                    FOREIGN KEY (group_id, user_id)
                                        REFERENCES group_members(group_id, user_id)
);

CREATE INDEX expense_shares_group_user_idx
    ON expense_shares(group_id, user_id);


-- +migrate Down

DROP TABLE IF EXISTS expense_shares;
DROP TABLE IF EXISTS expense_payers;