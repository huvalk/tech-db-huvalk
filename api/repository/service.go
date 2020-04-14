package repository

import "github.com/huvalk/tech-db-huvalk/api/models"

func (r *PostgresRepository) ClearAll() (stat int) {
	clear := `TRUNCATE users, forum, post, thread, vote, users_of_forum;`
	_, err := r.db.Exec(clear)
	if err != nil {
		return 500
	}

	return 200
}

func (r *PostgresRepository) GetStatus() (res *models.Status, stat int) {
	res = &models.Status{}
	r.db.QueryRow("SELECT reltuples AS approximate_row_count FROM pg_class WHERE relname = 'forum';").Scan(&res.Forum)
	r.db.QueryRow("SELECT reltuples AS approximate_row_count FROM pg_class WHERE relname = 'post';").Scan(&res.Post)
	r.db.QueryRow("SELECT reltuples AS approximate_row_count FROM pg_class WHERE relname = 'thread';").Scan(&res.Thread)
	r.db.QueryRow(`SELECT reltuples AS approximate_row_count FROM pg_class WHERE relname = 'users';`).Scan(&res.User)
	return res, 200
}
