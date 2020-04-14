package repository

import (
	"github.com/huvalk/tech-db-huvalk/api/models"
	"log"
	"strconv"
	"time"
)

func (r *PostgresRepository) CreateThread(thread *models.Thread) (res *models.Thread, stat int) {
	var userID int64
	err := r.db.QueryRow("select_user_id_nickname", thread.Author).Scan(&userID, &thread.Author)
	if err != nil {
		return nil, 404
	}

	var forumID int64
	err = r.db.QueryRow("get_forum_id", thread.Forum).Scan(&forumID, &thread.Forum)
	if err != nil {
		return nil, 404
	}

	var threadID uint64
	if thread.Created.IsZero() {
		thread.Created = time.Now()
	}
	// TODO Тут было UTC
	err = r.db.QueryRow("create_thread", thread.Slug, thread.Title, thread.Message, userID, thread.Author,
		forumID, thread.Forum, thread.Created).Scan(&thread.ID)

	if err == nil {
		r.db.QueryRow("get_thread_by_id", threadID).Scan(&thread.Title, &thread.ID, &thread.Slug, &thread.Message,
			&thread.Author, &thread.Forum, &thread.Created)

		return thread, 201
	} else {
		err = r.db.QueryRow("get_thread_by_slug", thread.Slug).Scan(&thread.Title, &thread.ID, &thread.Slug,
			&thread.Message, &thread.Author, &thread.Forum, &thread.Created)

		return thread, 409
	}
}

func (r *PostgresRepository) GetListOfThreads(slug, limit, since, desc string) (models.Threads, int) {
	var timeStamp time.Time
	limitInt, _ := strconv.Atoi(limit)

	var forumID int64
	err := r.db.QueryRow("select_forum_id", slug).Scan(&forumID)
	if err != nil {
		return nil, 404
	}

	getThread := `SELECT t.author_name, t.created, t.forum_slug, t.id, t.message, COALESCE(t.slug, ''), t.title, t.votes
				FROM thread t
				WHERE t.forum_id = $1`
	if desc == "true" {
		if since == "" {
			timeStamp = time.Now().AddDate(10, 0, 0)
		} else {
			timeStamp, _ = time.Parse(time.RFC3339, since)
		}
		getThread += "\nAND t.created <= $2\nORDER BY t.created desc LIMIT $3;"
	} else {
		if since == "" {
			timeStamp = time.Now().AddDate(-200, 0, 0)
		} else {
			timeStamp, _ = time.Parse(time.RFC3339, since)
		}
		getThread += "\nAND t.created >= $2\nORDER BY t.created asc LIMIT $3;"
	}

	rows, err := r.db.Query(getThread, forumID, timeStamp, limitInt)
	res := make(models.Threads, 0)
	defer rows.Close()
	for rows.Next() {
		var threadDB models.Thread
		rows.Scan(&threadDB.Author, &threadDB.Created, &threadDB.Forum, &threadDB.ID, &threadDB.Message,
			&threadDB.Slug, &threadDB.Title, &threadDB.Votes)

		res = append(res, &threadDB)
	}

	return res, 200
}

func (r *PostgresRepository) VoteForThread(slug string, vote *models.Vote) (res *models.Thread, stat int) {
	threadID, errParse := strconv.ParseInt(slug, 10, 64)

	if errParse != nil {
		result, err := r.db.Exec("vote_for_thread_with_slug", vote.Voice, slug, vote.Nickname)
		if err != nil {
			log.Print(err)
		}

    if rowsA := result.RowsAffected(); rowsA == 0 {
			return nil, 404
		}

		res = &models.Thread{}
		r.db.QueryRow("get_thread_with_slug", slug).Scan(&res.Author, &res.Created, &res.Forum, &res.ID, &res.Message,
			&res.Slug, &res.Title, &res.Votes)
	} else {
		result, err := r.db.Exec("vote_for_thread_with_id", vote.Voice, threadID, vote.Nickname)
		if err != nil {
			log.Print(err)
		}

		//TODO Ввести NullString
		if rowsA := result.RowsAffected(); rowsA == 0 {
			return nil, 404
		}

		res = &models.Thread{}
		r.db.QueryRow("get_thread_with_id", threadID).Scan(&res.Author, &res.Created, &res.Forum, &res.ID, &res.Message,
			&res.Slug, &res.Title, &res.Votes)
	}

	return res, 200
}

func (r *PostgresRepository) GetThread(slugThread string) (res *models.Thread, stat int) {
	res = &models.Thread{}
	var err error
	threadID, err := strconv.ParseInt(slugThread, 10, 64)

	if err == nil {
		err = r.db.QueryRow("get_thread_with_id", threadID).
			Scan(&res.Author, &res.Created, &res.Forum,
				&res.ID, &res.Message,
				&res.Slug, &res.Title, &res.Votes)
	} else {
		err = r.db.QueryRow("get_thread_with_slug", slugThread).
			Scan(&res.Author, &res.Created, &res.Forum,
				&res.ID, &res.Message,
				&res.Slug, &res.Title, &res.Votes)
	}

	if err != nil {
		return nil, 404
	}
	return res, 200
}

func (r *PostgresRepository) ChangeThread(slugThread string, thread *models.ThreadUpdate) (res *models.Thread, stat int) {
	res, stat = r.GetThread(slugThread)
	if stat == 404 {
		return nil, 404
	}
	if thread.Title != "" {
		res.Title = thread.Title
	}
	if thread.Message != "" {
		res.Message = thread.Message
	}

	r.db.Exec("change_thread", res.Message, res.Title, res.ID)

	return res, 200
}
