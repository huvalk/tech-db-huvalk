CREATE UNIQUE INDEX index_cover_votes ON vote (voice, user_id, thread_id);

CREATE INDEX index_post_thread_parents ON post (thread_id, parents);

CREATE INDEX index_post_root ON post (root);

CREATE UNIQUE INDEX index_post_thread_id ON post (thread_id, id);

CREATE UNIQUE INDEX index_post_thread_id_parent_root ON post (thread_id, id, parent, root);

CREATE UNIQUE INDEX index_post_id_root ON post (id, root);

CREATE UNIQUE INDEX index_user_nickname ON users (nickname, id);

CREATE UNIQUE INDEX index_forum_slug_id ON forum (slug, id);

CREATE INDEX index_thread_forum_id ON thread (forum_id);

CREATE INDEX index_thread_forum_slug ON thread (forum_slug);

CREATE INDEX index_thread_sort ON thread (forum_id, created);