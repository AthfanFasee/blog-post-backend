CREATE INDEX IF NOT EXISTS posts_title_idx ON posts USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS posts_createdby_idx ON posts(created_by);
CREATE INDEX IF NOT EXISTS comments_postid_idx ON comments(post_id);
CREATE INDEX IF NOT EXISTS tokens_userid_idx ON tokens(user_id);
