CREATE TABLE reviews (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    place_id text NOT NULL,
    user_id uuid NOT NULL,
    rating double precision NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT reviews_rating_range CHECK (rating >= 1 AND rating <= 5)
);

CREATE INDEX idx_reviews_place_id ON reviews (place_id);
