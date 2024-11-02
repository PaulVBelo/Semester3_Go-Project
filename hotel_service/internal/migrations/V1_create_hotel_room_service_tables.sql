CREATE TABLE hotel (
    hotel_id            BIGSERIAL           PRIMARY KEY,
    hotel_name          VARCHAR(64)         NOT NULL,
    hotel_address       VARCHAR(256)        NOT NULL,
    phone_number        VARCHAR(32)         NOT NULL CHECK (phone_number ~ '^\+?\d+$')
);

CREATE UNIQUE INDEX idx_name_adress ON hotel (hotel_name, hotel_address);

CREATE TABLE room (
    room_id             BIGSERIAL           PRIMARY KEY,
    room_name           VARCHAR(128)        NOT NULL,
    hotel_id            BIGINT              REFERENCES  hotel.hotel_id,
    price               NUMERIC(18, 2)      NOT NULL CHECK (price > 0)
);

CREATE TABLE amenity (
    amenity_id          BIGSERIAL           PRIMARY KEY,
    amenity_name       VARCHAR(128)         NOT NULL,
    hotel_id            BIGINT              REFERENCES  hotel.hotel_id      
);

CREATE UNIQUE INDEX idx_name_hotel ON hotel (hotel_name, hotel_address);

CREATE TABLE room_x_amenity (
    room_id             BIGINT              REFERENCES room.room_id,
    amenity_id          BIGINT              REFERENCES amenity.amenity_id
);