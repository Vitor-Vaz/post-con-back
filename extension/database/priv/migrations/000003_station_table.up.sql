CREATE TABLE station (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    place_id text NOT NULL,
    name text NOT NULL,
    address text,
    latitude double precision,
    longitude double precision,
    total_score double precision NOT NULL DEFAULT 0,
    summary text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT station_place_id_unique UNIQUE (place_id)
);

CREATE INDEX idx_station_place_id ON station (place_id);
