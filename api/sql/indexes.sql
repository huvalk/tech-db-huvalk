CREATE UNIQUE INDEX index_cover_votes ON vote (voice, user_id, thread_id);

CREATE INDEX index_post_thread_parents ON post (thread_id, parents);

CREATE INDEX index_post_root ON post (root);

CREATE UNIQUE INDEX index_post_thread_id ON post (thread_id, id);

CREATE UNIQUE INDEX index_post_thread_id_parent_root ON post (thread_id, id, parent, root);

CREATE UNIQUE INDEX index_post_id_root ON post (id, root);

-- CREATE UNIQUE INDEX index_user_nickname ON users (nickname, id);
-- index_users_full

-- CREATE UNIQUE INDEX index_forum_slug_id ON forum (slug, id);
-- index_forum_full

CREATE INDEX index_thread_sort ON thread (forum_id, created);

-- ////
--
-- CREATE UNIQUE INDEX index_forum_title_slug_moderator ON forum (slug, id) include (posts, threads, title, moderator_name);
--
-- // заменить коаклис на нул, посмотреть, что получиться
-- CREATE UNIQUE INDEX index_thread_id_slug ON thread (id, COALESCE(thread.slug, ''));
--
-- // можно добавить, неплохо используется
-- CREATE UNIQUE INDEX index_users_of_forum_forum_id_user_id ON users_of_forum (forum_id, user_id);
--
-- //

-- // покрывающие

CREATE UNIQUE INDEX index_users_of_forum_forum_id_user_id ON users_of_forum (forum_id, user_nickname);
CREATE INDEX index_users_full ON users (nickname, id, email, about, fullname);
CREATE INDEX index_forum_full on forum (slug, id, title, moderator_name, threads, posts);
CREATE UNIQUE INDEX index_thread_full ON thread (forum_id, created, id, slug, title, message, forum_slug, author_name, created, votes);