-- +goose Up
CREATE TYPE event_names as ENUM('sighted','commentary','area_checked','found','rewarded');
CREATE TABLE pets_records(
  id UUID NOT NULL PRIMARY KEY,
  createdAt TIMESTAMP NOT NULL,
  updatedAt TIMESTAMP NOT NULL,

  pet_id UUID NOT NULL REFERENCES missing_pets(id) ON DELETE CASCADE,
  user_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  event_name event_names NOT NULL,
  image_url text,
  description text ,

    event_location  GEOMETRY(Point, 4326),
    area  GEOGRAPHY(Polygon, 4326),
    event_time   TIMESTAMP 

);

-- +goose Down
DROP TABLE pets_records;
DROP TYPE event_names;
