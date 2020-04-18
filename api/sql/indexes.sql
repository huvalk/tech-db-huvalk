CREATE UNIQUE INDEX index_cover_votes ON vote (voice, user_id, thread_id);

CREATE INDEX index_post_thread_parents ON post (thread_id, parents);

CREATE INDEX index_post_root ON post (root);

CREATE UNIQUE INDEX index_post_thread_id ON post (thread_id, id);

CREATE UNIQUE INDEX index_post_thread_id_parent_root ON post (thread_id, id, parent, root);
-- drop index index_post_thread_id_parent_root;
-- заменено index_post_parent_thread_id и index_post_id_parents_root


CREATE UNIQUE INDEX index_post_id_root ON post (id, root);

CREATE UNIQUE INDEX index_user_nickname ON users (nickname, id);
-- drop index index_user_nickname;
-- index_users_full

CREATE UNIQUE INDEX index_forum_slug_id ON forum (slug, id);
-- drop index index_forum_slug_id;
-- index_forum_full

-- ////
--
-- CREATE UNIQUE INDEX index_forum_title_slug_moderator ON forum (slug, id) include (posts, threads, title, moderator_name);
--
-- // заменить коаклис на нул, посмотреть, что получиться
-- CREATE UNIQUE INDEX index_thread_id_slug ON thread (id, COALESCE(thread.slug, ''));
--
-- // покрывающие

CREATE  INDEX index_users_of_forum_forum_id_user_id ON users_of_forum (forum_id, user_id);
CREATE  INDEX index_users_full ON users (nickname, id) include (email, about, fullname);

-- // не используетс
-- CREATE  INDEX index_forum_full on forum (slug, id) include (posts, slug, threads, title, moderator_name);

CREATE  INDEX index_thread_full ON thread (forum_id, created) include (id, slug, title, message, forum_slug, author_name, created, votes);

-- // небольшая прибавка есть
--
-- drop index index_users_of_forum_forum_id_user_id
-- drop index index_users_full
-- -- drop index index_forum_full
-- drop index index_thread_full
--
-- ///

CREATE UNIQUE INDEX index_post_id_parents_root ON post (id) include (parents, root, thread_id);
-- drop index index_post_id_parents_root

CREATE INDEX index_post_parent_thread_id ON post (parent, thread_id) include (id);
-- drop index index_post_parent_thread_id
