CREATE EXTENSION citext;

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       nickname CITEXT COLLATE "ucs_basic" NOT NULL UNIQUE,
                       fullname TEXT NOT NULL,
                       email CITEXT COLLATE "ucs_basic"  NOT NULL UNIQUE,
                       about TEXT
);

CREATE TABLE forum (
                       id SERIAL PRIMARY KEY,
                       slug CITEXT NOT NULL UNIQUE,
                       title TEXT NOT NULL,
                       moderator_id INTEGER REFERENCES users (id) NOT NULL,
                        posts BIGINT CONSTRAINT positive_posts CHECK (posts >= 0) DEFAULT 0,
                        threads INT CONSTRAINT positive_threads CHECK (threads >= 0) DEFAULT 0
);

CREATE OR REPLACE FUNCTION forum_posts() RETURNS TRIGGER AS $forum_posts$
BEGIN
    UPDATE forum SET posts=posts + 1
    FROM thread
    WHERE new.thread_id = thread.id
    AND thread.forum_id = forum.id;
RETURN NULL;
END;
$forum_posts$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION forum_threads() RETURNS TRIGGER AS $forum_threads$
BEGIN
    UPDATE forum SET threads=threads + 1
    WHERE new.forum_id = id;
    RETURN NULL;
END;
$forum_threads$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS users_of_forum (
    user_id INTEGER REFERENCES users(id) NOT NULL,
    forum_id INTEGER REFERENCES forum(id) NOT NULL,
    UNIQUE (user_id, forum_id)
);

CREATE OR REPLACE FUNCTION forum_user_post() RETURNS TRIGGER AS $forum_user_post$
BEGIN
    INSERT INTO users_of_forum
    (forum_id, user_id)
               SELECT f.id, new.author_id
               FROM forum f, thread t
               WHERE new.thread_id = t.id
                AND t.forum_id = f.id
    ON CONFLICT DO NOTHING;
    RETURN NULL;
END;
$forum_user_post$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION forum_user_thread() RETURNS TRIGGER AS $forum_user_thread$
BEGIN
    INSERT INTO users_of_forum
    (user_id, forum_id)
    VALUES (new.author_id, new.forum_id)
    ON CONFLICT DO NOTHING;
    RETURN NULL;
END;
$forum_user_thread$ LANGUAGE plpgsql;

CREATE TABLE thread (
                        id SERIAL PRIMARY KEY,
                        author_id INTEGER REFERENCES users(id) NOT NULL,
                        forum_id INTEGER REFERENCES forum(id) NOT NULL,
                        title TEXT NOT NULL,
                        message TEXT NOT NULL,
                        slug CITEXT NULL UNIQUE,
                        CREATED TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        votes INT CONSTRAINT positive_votes CHECK (votes >= 0) DEFAULT 0
);

CREATE OR REPLACE FUNCTION thread_new_vote() RETURNS TRIGGER AS $thread_new_vote$
BEGIN
    UPDATE thread SET votes=votes + new.voice
    WHERE new.thread_id = id;
    RETURN NULL;
END;
$thread_new_vote$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION thread_change_vote() RETURNS TRIGGER AS $thread_change_vote$
BEGIN
    UPDATE thread SET votes=votes - old.voice + new.voice
    WHERE new.thread_id = id;
    RETURN NULL;
END;
$thread_change_vote$ LANGUAGE plpgsql;

CREATE TRIGGER inc_threads
    AFTER INSERT ON thread
    FOR EACH ROW EXECUTE PROCEDURE forum_threads();

CREATE TRIGGER thread_user
    AFTER INSERT ON thread
    FOR EACH ROW EXECUTE PROCEDURE forum_user_thread();

CREATE TABLE post (
                      id SERIAL PRIMARY KEY,
                      author_id INTEGER REFERENCES users(id) NOT NULL,
                      CREATED TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                      thread_id INTEGER REFERENCES thread(id) NOT NULL,
                      is_edited BOOLEAN NOT NULL DEFAULT FALSE,
                      message TEXT NOT NULL DEFAULT '',
                      parent INTEGER REFERENCES post(id) NULL DEFAULT NULL,
                      parents  INTEGER [],
                      root INTEGER
);

CREATE TRIGGER inc_posts
    AFTER INSERT ON post
    FOR EACH ROW EXECUTE PROCEDURE forum_posts();

CREATE TRIGGER post_user
    AFTER INSERT ON post
    FOR EACH ROW EXECUTE PROCEDURE forum_user_post();

CREATE TABLE vote (
                      id SERIAL PRIMARY KEY,
                      user_id INTEGER REFERENCES users(id) NOT NULL,
                      thread_id INTEGER REFERENCES thread(id) NOT NULL,
                      voice SMALLINT NOT NULL CONSTRAINT voice_in_array CHECK (voice = 1 OR voice = -1),
                    UNIQUE (user_id, thread_id)
);

CREATE TRIGGER new_votes
    AFTER INSERT ON vote
    FOR EACH ROW EXECUTE PROCEDURE thread_new_vote();

CREATE TRIGGER change_votes
    AFTER UPDATE ON vote
    FOR EACH ROW EXECUTE PROCEDURE thread_change_vote();