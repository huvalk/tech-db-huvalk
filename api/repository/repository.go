package repository

import (
	"github.com/jackc/pgx"
)

type PostgresRepository struct {
	db *pgx.ConnPool
}

func NewPostgresRepository(db *pgx.ConnPool) *PostgresRepository {
	// CreateUser
	db.Prepare("create_user", "INSERT INTO users (nickname, fullname, email, about) values ($1, $2, $3, $4);")
	db.Prepare("create_user_select", "SELECT nickname, fullname, email, about FROM users WHERE nickname=$1 OR email=$2;")
	// GetUser
	db.Prepare("get_user", "SELECT nickname, fullname, email, about FROM users WHERE nickname=$1")
	// ChangeUser
	db.Prepare("get_user_by_id", "SELECT nickname, fullname, email, about FROM users WHERE id=$1")
	db.Prepare("change_user",
		`UPDATE users SET 
                 fullname=coalesce(NULLIF($1, ''), fullname), 
                 email=coalesce(NULLIF($2, ''), email),
                 about=coalesce(NULLIF($3, ''), about) 
				WHERE id=$4`)
	// GetListOfUsers
	db.Prepare("select_forum_id",
		`SELECT f.id
					FROM forum f
					WHERE f.slug=$1;`)

	// CreateForum
	db.Prepare("select_user_id", "SELECT id FROM users WHERE nickname=$1;")
	db.Prepare("select_user_id_nickname", "SELECT id, nickname FROM users WHERE nickname=$1;")
	db.Prepare("create_forum", "INSERT INTO forum (slug, title, moderator_id, moderator_name) VALUES ($1, $2, $3, $4);")
	db.Prepare("create_forum_select",
		`SELECT title, slug, moderator_name
					FROM forum
					WHERE slug=$1;`)
	// GetForum
	db.Prepare("get_forum",
		`SELECT posts, slug, threads, title, moderator_name
				FROM forum 
				WHERE forum.slug = $1;`)

	// CreateThread
	db.Prepare("get_forum_id",
		`SELECT id, slug FROM forum WHERE slug=$1;`)
	db.Prepare("create_thread",
		`INSERT INTO thread (slug, title, message, author_id, author_name, forum_id, forum_slug, created) 
					values (NULLIF($1, ''), $2, $3, $4, $5, $6, $7, $8) RETURNING id;`)
	db.Prepare("get_thread_by_id",
		`SELECT thread.title, thread.id, COALESCE(thread.slug, ''), thread.message, thread.author_name, 
				thread.forum_slug, thread.created
				FROM thread
				WHERE thread.id=$1;`)
	db.Prepare("get_thread_by_slug",
		`SELECT thread.title, thread.id, COALESCE(thread.slug, ''), thread.message, thread.author_name, 
				thread.forum_slug, thread.created
				FROM thread
				WHERE thread.slug=$1;`)

	// GetThread
	db.Prepare("get_thread_with_id",
		`SELECT t.author_name, t.created, t.forum_slug, t.id, t.message, COALESCE(t.slug, ''), t.title, t.votes
				FROM thread  t
				WHERE t.id = $1;`)
	db.Prepare("get_thread_with_slug",
		`SELECT t.author_name, t.created, t.forum_slug, t.id, t.message, COALESCE(t.slug, ''), t.title, t.votes
				FROM thread  t
				WHERE t.slug = $1;`)
	db.Prepare("vote_for_thread_with_id",
		`INSERT INTO vote (voice, thread_id, user_id)
						   SELECT $1, thread.id, users.id
							FROM thread, users
						   WHERE thread.id=$2
							AND users.nickname=$3
       				ON CONFLICT (user_id, thread_id) DO UPDATE SET voice=$1`)
	db.Prepare("vote_for_thread_with_slug",
		`INSERT INTO vote (voice, thread_id, user_id)
						   SELECT $1, thread.id, users.id
							FROM thread, users
						   WHERE thread.slug=$2
							AND users.nickname=$3
       				ON CONFLICT (user_id, thread_id) DO UPDATE SET voice=$1`)


	// ChangeThread
	db.Prepare("change_thread",
		`UPDATE thread SET 
                 message=$1, 
                 title=$2
				WHERE id=$3;`)

	// CreatePosts
	// TODO придется разбить?
	db.Prepare("get_thread_id",
		`SELECT COALESCE(t.slug, ''), t.forum_slug, t.forum_id
					FROM thread t
					WHERE t.id=$1;`)
	db.Prepare("get_thread_slug",
		`SELECT t.id, t.forum_slug, t.forum_id
					FROM thread t
					WHERE t.slug=$1;`)
	db.Prepare("get_nextval",
		`SELECT
					array_agg(nextval('post_id_seq'))
					FROM generate_series(1, $1)`)

	//TODO индекс на post.id, post.thread_id
	db.Prepare("get_parent",
		`SELECT post.parents
					FROM post
					WHERE post.id=$1
					AND post.thread_id=$2;`)
	db.Prepare("create_post",
		`INSERT INTO post (author_id, author_name, forum_slug, forum_id, message, parent, thread_id, created, parents, root, id) 
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`)
	// GetListOfPost
	db.Prepare("get_thread_id_and_slug_by_id",
		`SELECT COALESCE(thread.slug, '')
					FROM thread
					WHERE thread.id=$1;`)
	db.Prepare("get_thread_id_and_slug_by_slug",
		`SELECT thread.id
					FROM thread
					WHERE thread.slug=$1;`)
	// ChangePost
	db.Prepare("change_post_select",
		`SELECT p.author_name, p.created, p.forum_slug, p.id, p.is_edited, p.message, p.Parent, p.thread_id
				FROM post p
				WHERE p.id = $1;`)
	db.Prepare("change_post",
		`UPDATE post p SET
                 message=$1,
                  is_edited=true
                 WHERE p.id=$2;`)
	db.Prepare("get_post",
		`SELECT p.author_name, p.created, p.forum_slug, p.id, p.is_edited, p.message, p.Parent, p.thread_id
				FROM post p
				WHERE p.id = $1;`)

	// VoteForThread
	db.Prepare("vote_for_thread",
		`INSERT INTO vote (voice, thread_id, user_id)
						   SELECT $1, thread.id, users.id
							FROM thread, users
						   WHERE thread.slug=$2
							AND users.nickname=$3
       				ON CONFLICT (user_id, thread_id) DO UPDATE SET voice=$1`)


	return &PostgresRepository{
		db: db,
	}
}
