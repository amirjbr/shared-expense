-- +migrate Up

CREATE TABLE group_invitations (
                                   id UUID PRIMARY KEY,

                                   group_id UUID NOT NULL,
                                   invited_user_id UUID NOT NULL,
                                   invited_by_user_id UUID NOT NULL,

                                   status VARCHAR(20) NOT NULL DEFAULT 'pending',

                                   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                   responded_at TIMESTAMPTZ,

                                   CONSTRAINT group_invitations_group_fk
                                       FOREIGN KEY (group_id)
                                           REFERENCES groups(id)
                                           ON DELETE CASCADE,

                                   CONSTRAINT group_invitations_invited_user_fk
                                       FOREIGN KEY (invited_user_id)
                                           REFERENCES users(id),

                                   CONSTRAINT group_invitations_inviter_member_fk
                                       FOREIGN KEY (group_id, invited_by_user_id)
                                           REFERENCES group_members(group_id, user_id),

                                   CONSTRAINT group_invitations_status_check
                                       CHECK (
                                           status IN (
                                                      'pending',
                                                      'accepted',
                                                      'rejected'
                                               )
                                           ),

                                   CONSTRAINT group_invitations_no_self_invite
                                       CHECK (invited_user_id <> invited_by_user_id)
);

CREATE UNIQUE INDEX group_invitations_unique_pending_idx
    ON group_invitations(group_id, invited_user_id)
    WHERE status = 'pending';

CREATE INDEX group_invitations_invited_user_status_idx
    ON group_invitations(invited_user_id, status);


-- +migrate Down

DROP TABLE IF EXISTS group_invitations;