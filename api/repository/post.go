package repository

import (
	"github.com/huvalk/tech-db-huvalk/api/models"
	"github.com/lib/pq"
	"strconv"
	"time"
)

func (r *PostgresRepository) CreatePosts(threadSlug string, posts models.Posts) (res models.Posts, stat int) {
	var (
		userID      int64
		slugForum   string
		forumId     int64
		threadID    int64
		currentTime time.Time
		err         error
	)

	lenPosts := len(posts)

	if lenPosts != 0 {
		if posts[0].Created.IsZero() {
			currentTime = time.Now().Truncate(time.Millisecond)
		} else {
			currentTime = posts[0].Created
		}
	}

	threadID, err = strconv.ParseInt(threadSlug, 10, 64)
	if err == nil {
		err = r.db.QueryRow("get_thread_id", threadID).Scan(&threadSlug, &slugForum, &forumId)
	} else {
		err = r.db.QueryRow("get_thread_slug", threadSlug).Scan(&threadID, &slugForum, &forumId)
	}
	if err != nil {
		return nil, 404
	} else if lenPosts == 0 {
		return posts, 201
	}

	newId := make(pq.Int64Array, lenPosts)
	err = r.db.QueryRow("get_nextval", lenPosts).Scan(&newId)

	for i, post := range posts {
		parents := make(pq.Int64Array, 0)
		post.ID = newId[i]

		if post.Parent != 0 {
			err = r.db.QueryRow("get_parent", post.Parent, threadID).Scan(&parents)
			if err != nil {
				return nil, 409
			}
		}

		err = r.db.QueryRow("select_user_id_nickname", post.Author).Scan(&userID, &post.Author)
		if err != nil {
			return nil, 404
		}

		parents = append(parents, post.ID)
		_, err = r.db.Exec("create_post", userID, post.Author, slugForum, forumId, post.Message, post.Parent, threadID,
			currentTime, parents, parents[0], post.ID)

		post.Forum = slugForum
		post.Created = currentTime
		post.Thread = int32(threadID)
	}

	return posts, 201
}

func (r *PostgresRepository) GetListOfPosts(slugThread, limit, since, sort, desc string) (models.Posts, int) {
	if limit != "" {
		limit = "LIMIT " + limit
	}
	threadID, err := strconv.ParseInt(slugThread, 10, 64)
	if err == nil {
		err = r.db.QueryRow("get_thread_id_and_slug_by_id", threadID).Scan(&slugThread)
	} else {

		err = r.db.QueryRow("get_thread_id_and_slug_by_slug", slugThread).Scan(&threadID)
	}
	if err != nil {
		return nil, 404
	}

	switch sort {
	case "tree":
		if desc == "true" {
			if since != "" {
				//desc = "\nAND array_remove(parents, 0) < (SELECT array_remove(parents, 0) FROM post WHERE id = " + since + ")"
				desc = "\nAND parents < (SELECT parents FROM post WHERE id = " + since + ")"
			} else {
				desc = " "
			}
			desc += "\nORDER BY parents desc "
		} else {
			if since != "" {
				//desc = "\nAND array_remove(parents, 0) > (SELECT array_remove(parents, 0) FROM post WHERE id = " + since + ")"
				desc = "\nAND parents > (SELECT parents FROM post WHERE id = " + since + ")"
			} else {
				desc = " "
			}
			desc += "\nORDER BY parents asc "
		}

		getThread := `	SELECT p.author_name, p.created, p.forum_slug, p.id, p.is_edited, p.message, COALESCE(p.Parent, 0), p.thread_id
						FROM post p
						WHERE p.thread_id = $1` + desc + limit

		rows, _ := r.db.Query(getThread, threadID)
		res := models.Posts{}
		defer rows.Close()
		for rows.Next() {
			var postDB models.Post
			err = rows.Scan(&postDB.Author, &postDB.Created, &postDB.Forum, &postDB.ID, &postDB.IsEdited,
				&postDB.Message, &postDB.Parent, &postDB.Thread)

			res = append(res, &postDB)
		}

		return res, 200
	case "parent_tree":

		if desc == "true" {
			if since != "" {
				since = "\nAND root < (SELECT root FROM post WHERE id = " + since + ")"
			} else {
				since = " "
			}
			desc = "desc"
		} else {
			if since != "" {
				since = "\nAND root > (SELECT root FROM post WHERE id = " + since + ")"
			} else {
				since = " "
			}
			desc = "asc"
		}

		getThread := `	SELECT p.author_name, p.created, p.forum_slug, p.id, p.is_edited, p.message, p.Parent, p.thread_id
						FROM post p
						JOIN (
								SELECT id
								FROM post
								WHERE parent = 0 AND thread_id = $1 ` + since + `
								ORDER BY id ` + desc + `
								` + limit + `) sub ON p.root=sub.id
						ORDER BY p.root ` + desc + `, p.parents asc`
						//ORDER BY p.root ` + desc + `, array_remove(p.parents, p.root) asc`

		rows, _ := r.db.Query(getThread, threadID)
		defer rows.Close()

		res := models.Posts{}
		for rows.Next() {
			var postDB models.Post
			err = rows.Scan(&postDB.Author, &postDB.Created, &postDB.Forum, &postDB.ID, &postDB.IsEdited,
				&postDB.Message, &postDB.Parent, &postDB.Thread)

			res = append(res, &postDB)
		}

		return res, 200
	default:
		getThread := `SELECT p.author_name, p.created, p.forum_slug, p.id, p.is_edited, p.message, p.Parent, p.thread_id
				FROM post p
				WHERE p.thread_id = $1`

		if desc == "true" {
			if since != "" {
				getThread += "\nAND p.id < " + since
			}
			getThread += "\nORDER BY p.created desc, p.id desc "
		} else {
			if since != "" {
				getThread += "\nAND p.id > " + since
			}
			getThread += "\nORDER BY p.created asc, p.id asc "
		}
		getThread += limit

		rows, _ := r.db.Query(getThread, threadID)
		res := make(models.Posts, 0)
		defer rows.Close()

		for rows.Next() {
			var postDB models.Post
			err = rows.Scan(&postDB.Author, &postDB.Created, &postDB.Forum, &postDB.ID, &postDB.IsEdited,
				&postDB.Message, &postDB.Parent, &postDB.Thread)

			res = append(res, &postDB)
		}

		return res, 200
	}
}

func (r *PostgresRepository) PostDetails(slug string, related []string) (res models.PostFull, stat int) {
	postID, _ := strconv.ParseInt(slug, 10, 64)
	res.Post = &models.Post{}
	//TODO Ускорить
	err := r.db.QueryRow("get_post", postID).Scan(&res.Post.Author, &res.Post.Created, &res.Post.Forum, &res.Post.ID,
		&res.Post.IsEdited, &res.Post.Message, &res.Post.Parent, &res.Post.Thread)
	if err != nil {
		return res, 404
	}

	for _, rel := range related {
		switch rel {
		case "user":
			result, _ := r.GetUser(res.Post.Author)
			res.Author = result
		case "forum":
			result, _ := r.GetForum(res.Post.Forum)
			res.Forum = result
		case "thread":
			result, _ := r.GetThread(strconv.Itoa(int(res.Post.Thread)))
			res.Thread = result
		}
	}
	return res, 200
}

func (r *PostgresRepository) ChangePost(slug string, post *models.PostUpdate) (res *models.Post, stat int) {
	postID, _ := strconv.ParseInt(slug, 10, 64)
	res = &models.Post{}
	err := r.db.QueryRow("change_post_select", postID).Scan(&res.Author, &res.Created, &res.Forum, &res.ID,
		&res.IsEdited, &res.Message, &res.Parent, &res.Thread)
	if err != nil {
		return res, 404
	}

	if post.Message != "" && post.Message != res.Message {
		_, _ = r.db.Exec("change_post", post.Message, postID)
		res.Message = post.Message
		res.IsEdited = true
	}

	return res, 200
}
