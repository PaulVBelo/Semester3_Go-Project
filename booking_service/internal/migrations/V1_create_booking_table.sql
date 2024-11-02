CREATE TABLE booking (
    booking_id          BIGSERIAL                       PRIMARY KEY,
    room_id             BIGINT                          NOT NULL,
    time_from           TIMESTAMP WITHOUT TIME ZONE     NOT NULL,
    time_to             TIMESTAMP WITHOUT TIME ZONE     NOT NULL,
    client_number       VARCHAR(32)                     NOT NULL CHECK (phone_number ~ '^\+?\d+$'),
    booking_status      VARCHAR(16)                     NOT NULL
);

-- Чекать на time_from < time_to реклмендую в проге, но можете и добавить check в БД.