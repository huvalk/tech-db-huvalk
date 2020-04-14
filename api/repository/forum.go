package repository

import (
	"github.com/huvalk/tech-db-huvalk/api/models"
)

func (r *PostgresRepository) CreateForum(forum *models.Forum) (*models.Forum, int) {
	var userID int64
	err := r.db.QueryRow("select_user_id_nickname", forum.User).Scan(&userID, &forum.User)
	if err != nil {
		return nil, 404
	}

	_, err = r.db.Exec("create_forum", forum.Slug, forum.Title, userID, forum.User)
	if err == nil {
		return nil, 201
	}

	r.db.QueryRow("create_forum_select", forum.Slug).Scan(&forum.Title, &forum.Slug, &forum.User)

	return forum, 409
}

func (r *PostgresRepository) GetForum(slug string) (res *models.Forum, stat int) {
	res = &models.Forum{}

	err := r.db.QueryRow("get_forum", slug).
		Scan(&res.Posts, &res.Slug, &res.Threads, &res.Title, &res.User)
	if err != nil {
		return nil, 404
	}
	return res, 200
}
