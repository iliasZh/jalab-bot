CREATE TABLE users
(
    id               BIGINT                NOT NULL
        CONSTRAINT users_pk
            PRIMARY KEY,
    username         VARCHAR(32),
    first_name       VARCHAR(64)           NOT NULL
        CONSTRAINT check_first_name
            CHECK ((first_name)::TEXT <> ''::TEXT),
    created_at       TIMESTAMP             NOT NULL,
    updated_at       TIMESTAMP             NOT NULL,
    disable_mentions BOOLEAN DEFAULT FALSE NOT NULL
);

ALTER TABLE users
    OWNER TO admin;

CREATE TABLE jalabs
(
    id            SERIAL
        CONSTRAINT jalabs_pk
            PRIMARY KEY,
    group_chat_id BIGINT    NOT NULL,
    user_id       BIGINT    NOT NULL
        CONSTRAINT jalabs_users_id_fk
            REFERENCES users,
    created_at    TIMESTAMP NOT NULL,
    updated_at    TIMESTAMP NOT NULL
);

ALTER TABLE jalabs
    OWNER TO admin;

CREATE UNIQUE INDEX jalabs_group_chat_id_user_id_uindex
    ON jalabs (group_chat_id, user_id);

CREATE TABLE todays_jalabs
(
    id            SERIAL
        CONSTRAINT todays_jalabs_pk
            PRIMARY KEY,
    user_id       BIGINT NOT NULL,
    group_chat_id BIGINT NOT NULL,
    created_at    DATE   NOT NULL,
    CONSTRAINT todays_jalabs_jalabs_group_chat_id_user_id_fk
        FOREIGN KEY (group_chat_id, user_id) REFERENCES jalabs (group_chat_id, user_id)
);

ALTER TABLE todays_jalabs
    OWNER TO admin;

CREATE UNIQUE INDEX todays_jalabs_group_chat_id_created_at_uindex
    ON todays_jalabs (group_chat_id, created_at);

CREATE UNIQUE INDEX users_username_uindex
    ON users (username);

CREATE TABLE yaxshis
(
    id            BIGSERIAL
        CONSTRAINT yaxshis_pk
            PRIMARY KEY,
    count         BIGINT    NOT NULL,
    group_chat_id BIGINT    NOT NULL,
    user_id       BIGINT    NOT NULL,
    created_at    TIMESTAMP NOT NULL,
    updated_at    TIMESTAMP NOT NULL,
    CONSTRAINT yaxshis_jalabs_group_chat_id_user_id_fk
        FOREIGN KEY (group_chat_id, user_id) REFERENCES jalabs (group_chat_id, user_id)
);

ALTER TABLE yaxshis
    OWNER TO admin;

CREATE UNIQUE INDEX yaxshis_group_chat_id_user_id_uindex
    ON yaxshis (group_chat_id, user_id);

