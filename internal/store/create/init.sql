SET SYNCHRONOUS_COMMIT = 'off';

DROP INDEX IF EXISTS idx_users_email_uindex;
DROP INDEX IF EXISTS idx_users_nickname_uindex;
DROP INDEX IF EXISTS idx_forums_slug_uindex;
DROP INDEX IF EXISTS idx_forums_author_unique;;
DROP INDEX IF EXISTS idx_threads_slug;
DROP INDEX IF EXISTS idx_threads_forum;
DROP INDEX IF EXISTS idx_posts_forum;
DROP INDEX IF EXISTS idx_posts_parent;
DROP INDEX IF EXISTS idx_posts_path;
DROP INDEX IF EXISTS idx_posts_thread;
DROP INDEX IF EXISTS idx_posts_thread_id;

DROP TRIGGER IF EXISTS on_vote_insert ON votes;
DROP TRIGGER IF EXISTS on_vote_update ON votes;

DROP FUNCTION IF EXISTS fn_update_thread_votes_ins();
DROP FUNCTION IF EXISTS fn_update_thread_votes_upd();

DROP TABLE IF EXISTS forum_users;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS forums;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS votes;


CREATE TABLE IF NOT EXISTS users (
    id       bigserial not null primary key,
    nickname varchar(50) COLLATE "POSIX" not null unique,
    email    varchar(50)   not null,
    fullname varchar(100)   not null,
    about    text
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_uindex
    ON users (LOWER(email));
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_nickname_uindex
    ON users (LOWER(nickname));

CREATE TABLE IF NOT EXISTS forums (
    id       bigserial not null primary key,
    title    varchar,
    slug     varchar   not null,
    author   varchar   not null,
    posts    int default 0,
    threads  int default 0
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_forums_slug_uindex
    ON forums (LOWER(slug));
CREATE UNIQUE INDEX IF NOT EXISTS idx_forums_author_unique
    ON forums (LOWER(author));

CREATE TABLE IF NOT EXISTS threads (
    id      serial not null primary key,
    title   varchar,
    slug    varchar     not null,
    message varchar,
    votes   int         default 0,
    author  varchar     not null,
    forum   varchar     not null,
    created timestamptz DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_threads_slug
    ON threads (LOWER(slug));
CREATE INDEX IF NOT EXISTS idx_threads_forum
    ON threads (LOWER(forum));

CREATE TABLE IF NOT EXISTS posts (
    id       bigserial not null primary key,
    parent   bigint             DEFAULT NULL,
    message  text,
    path     bigint[]  NOT NULL DEFAULT '{0}',
    author   varchar,
    thread   int REFERENCES threads(id) NOT NULL,
    forum    varchar,
    created  timestamptz        DEFAULT now(),
    isEdited bool               DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_posts_path ON posts USING GIN (path);
CREATE INDEX IF NOT EXISTS idx_posts_thread ON posts (thread);
CREATE INDEX IF NOT EXISTS idx_posts_forum ON posts (forum);
CREATE INDEX IF NOT EXISTS idx_posts_parent ON posts (parent);
CREATE INDEX IF NOT EXISTS idx_posts_thread_id ON posts (thread, id);

CREATE TABLE IF NOT EXISTS votes(
    nickname varchar  REFERENCES users(nickname) NOT NULL,
    thread   int      REFERENCES threads(id) NOT NULL,
    voice    smallint NOT NULL CHECK (voice = 1 OR voice = -1),
    PRIMARY KEY (nickname, thread)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_votes_nickname_thread_unique
    ON votes (LOWER(nickname), thread);

CREATE FUNCTION fn_update_thread_votes_ins()
    RETURNS TRIGGER AS '
    BEGIN
        UPDATE threads
        SET
            votes = votes + NEW.voice
        WHERE id = NEW.thread;
        RETURN NULL;
    END;
' LANGUAGE plpgsql;


CREATE TRIGGER on_vote_insert
    AFTER INSERT ON votes
    FOR EACH ROW EXECUTE PROCEDURE fn_update_thread_votes_ins();

CREATE FUNCTION fn_update_thread_votes_upd()
    RETURNS TRIGGER AS '
    BEGIN
        IF OLD.voice = NEW.voice
        THEN
            RETURN NULL;
        END IF;
        UPDATE threads
        SET
            votes = votes + CASE WHEN NEW.voice = -1
                                     THEN -2
                                 ELSE 2 END
        WHERE id = NEW.thread;
        RETURN NULL;
    END;
' LANGUAGE plpgsql;

CREATE TRIGGER on_vote_update
    AFTER UPDATE ON votes
    FOR EACH ROW EXECUTE PROCEDURE fn_update_thread_votes_upd();


CREATE TABLE IF NOT EXISTS forum_users (
    user_id BIGINT REFERENCES users(id),
    forum_id BIGINT REFERENCES forums(id)
);

CREATE INDEX idx_forum_users_user_id
    ON forum_users(user_id);

CREATE INDEX idx_forum_users_forum_id
    ON forum_users(forum_id);

CREATE INDEX idx_forum_users_user_id_forum_id
    ON forum_users (user_id, forum_id);

CREATE OR REPLACE FUNCTION forum_users_update()
    RETURNS TRIGGER AS '
    BEGIN
        INSERT INTO forum_users (user_id, forum_id) VALUES ((SELECT id FROM users WHERE LOWER(NEW.author) = LOWER(nickname)),
                                                              (SELECT id FROM forums WHERE LOWER(NEW.forum) = LOWER(slug)));
        RETURN NULL;
    END;
' LANGUAGE plpgsql;

CREATE TRIGGER on_post_insert
    AFTER INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE forum_users_update();

CREATE TRIGGER on_thread_insert
    AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE forum_users_update();