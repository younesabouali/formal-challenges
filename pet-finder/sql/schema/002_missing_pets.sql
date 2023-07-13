-- +goose Up
CREATE TYPE missing_pets_status as ENUM ('missing','found');
CREATE TABLE missing_pets (
    id UUID NOT NULL PRIMARY KEY,
    createdAt TIMESTAMP NOT NULL,
    updatedAt TIMESTAMP NOT NULL,

    pet_name text NOT NULL,
    description text ,
    image_url text ,

    status missing_pets_status NOT NULL,
    lost_in  POINT NOT NULL,
    lost_at   TIMESTAMP NOT NULL,
    
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE



);

-- +goose Down
DROP TABLE missing_pets;
DROP TYPE missing_pets_status;
