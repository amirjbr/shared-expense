-- +migrate Up

CREATE TABLE settlements (
                             id UUID PRIMARY KEY,
                             group_id UUID NOT NULL,

                             payer_user_id UUID NOT NULL,
                             receiver_user_id UUID NOT NULL,

                             amount BIGINT NOT NULL,
                             settled_at TIMESTAMPTZ NOT NULL,

                             created_by_user_id UUID NOT NULL,

                             created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                             updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                             CONSTRAINT settlements_group_fk
                                 FOREIGN KEY (group_id)
                                     REFERENCES groups(id)
                                     ON DELETE CASCADE,

                             CONSTRAINT settlements_amount_positive
                                 CHECK (amount > 0),

                             CONSTRAINT settlements_different_parties
                                 CHECK (payer_user_id <> receiver_user_id),

                             CONSTRAINT settlements_creator_is_party
                                 CHECK (
                                     created_by_user_id = payer_user_id
                                         OR created_by_user_id = receiver_user_id
                                     ),

                             CONSTRAINT settlements_payer_is_group_member_fk
                                 FOREIGN KEY (group_id, payer_user_id)
                                     REFERENCES group_members(group_id, user_id),

                             CONSTRAINT settlements_receiver_is_group_member_fk
                                 FOREIGN KEY (group_id, receiver_user_id)
                                     REFERENCES group_members(group_id, user_id),

                             CONSTRAINT settlements_creator_is_group_member_fk
                                 FOREIGN KEY (group_id, created_by_user_id)
                                     REFERENCES group_members(group_id, user_id)
);

CREATE INDEX settlements_group_settled_at_idx
    ON settlements(group_id, settled_at DESC);

CREATE INDEX settlements_group_payer_idx
    ON settlements(group_id, payer_user_id);

CREATE INDEX settlements_group_receiver_idx
    ON settlements(group_id, receiver_user_id);


-- +migrate Down

DROP TABLE IF EXISTS settlements;