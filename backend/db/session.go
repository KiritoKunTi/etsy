package db

type Session struct {
	ID         int
	UUID       string
	Email      string
	User_ID    int
	Created_at string
}

func (session *Session) Check() (valid bool, err error) {
	err = DB.QueryRow(
		"SELECT ID, UUID, EMAIL, USER_ID, CREATED_AT FROM SESSIONS WHERE UUID=$1", session.UUID,
	).Scan(&session.ID, &session.UUID, &session.Email, &session.User_ID, &session.Created_at)
	if err != nil {
		valid = false
		return
	}
	if session.ID != 0 {
		valid = true
	}
	return
}

func (session *Session) DeleteByUUID() (err error) {
	statement := "delete from sessions where uuid=$1"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(session.UUID)
	return
}

func (session *Session) User() (user User, err error) {
	user, err = UserByID(session.User_ID)
	return
}

func SessionDeleteAll() (err error) {
	_, err = DB.Exec("DELETE FROM SESSIONS")
	return
}

