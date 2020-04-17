package repository

import (
	"github.com/huvalk/tech-db-huvalk/api/models"
)

func (r *PostgresRepository) CreateUser(user *models.User) (res models.Users, stat int) {
	_, err := r.db.Exec("create_user", user.Nickname, user.Fullname, user.Email, user.About)

	if err == nil {
		return nil, 201
	}

	rows, err := r.db.Query("create_user_select", user.Nickname, user.Email)
	defer rows.Close()

	res = make(models.Users, 0)
	for rows.Next() {
		var userDB models.User
		rows.Scan(&userDB.Nickname, &userDB.Fullname, &userDB.Email, &userDB.About)

		res = append(res, &userDB)
	}

	return res, 409
}

func (r *PostgresRepository) GetUser(user string) (res *models.User, stat int) {
	res = &models.User{}
	err := r.db.QueryRow("get_user", user).Scan(&res.Nickname, &res.Fullname, &res.Email, &res.About)
	if err != nil {
		return nil, 404
	}

	return res, 200
}

func (r *PostgresRepository) ChangeUser(userNickname string, user *models.UserUpdate) (res *models.User, stat int) {
	var userID int64
	err := r.db.QueryRow("select_user_id_nickname", userNickname).Scan(&userID, &userNickname)
	if err != nil {
		return nil, 404
	}

	_, err = r.db.Exec("change_user", user.Fullname, user.Email, user.About, userID)
	if err != nil {
		return nil, 409
	}

	res = &models.User{}
	r.db.QueryRow("get_user_by_id", userID).Scan(&res.Nickname, &res.Fullname, &res.Email, &res.About)

	return res, 200
}

func (r *PostgresRepository) GetListOfUsers(slug, limit, since, desc string) (models.Users, int) {
	//TODO упростить
	var (
		forumID int64
	)
	err := r.db.QueryRow("select_forum_id", slug).Scan(&forumID)
	if err != nil {
		return nil, 404
	}

	if desc == "true" {
		if since != "" {
			since = " AND u.nickname < '" + since + "'"
		}
		desc = "desc"
	} else {
		if since != "" {
			since = " AND u.nickname > '" + since + "'"
		}
		desc = "asc"
	}

	if limit != "" {
		limit = " LIMIT " + limit
	}

	getThread := "SELECT u.nickname, u.fullname, u.email, u.about " +
		"FROM users_of_forum f  " +
		"JOIN  users u on f.user_id = u.id " +
		"WHERE f.forum_id=$1 " + since +
		"ORDER BY u.nickname " + desc + limit

	rows, err := r.db.Query(getThread, forumID)
	defer rows.Close()
	res := make(models.Users, 0)
	for rows.Next() {
		var userDB models.User
		err = rows.Scan(&userDB.Nickname, &userDB.Fullname, &userDB.Email, &userDB.About)

		res = append(res, &userDB)
	}

	return res, 200
}
