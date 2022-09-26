ALTER TABLE posts ADD CONSTRAINT posts_readtime_check CHECK (read_time >= 0);
