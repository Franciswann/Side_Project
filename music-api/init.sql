-- Create the musics table with unique constraint
CREATE TABLE IF NOT EXISTS musics (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    UNIQUE(title, artist)
);

-- Create index for better performance
CREATE INDEX IF NOT EXISTS idx_musics_title ON musics(title);
CREATE INDEX IF NOT EXISTS idx_musics_artist ON musics(artist);
