package db

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID              int    `json:"id"`
	UUID            string `json:"UUID"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	IsShop          bool   `json:"isShop"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Repassword      string `json:"repassword"`
	Photo           string `json:"photo"`
	LanguageCode    string `json:"languageCode"`
	CreatedAt       string
	Email           string `json:"email"`
	Description     string `json:"description"`
	UsernameOrEmail string `json:"usernameOrEmail"`
}

var ErrExistsUsernameOrEmail = errors.New("Already have username or email on other account")

func UserByEmailOrUsername(usernameOrEmail string) (user User, err error) {
	if strings.Contains(usernameOrEmail, "@") {
		user, err = UserByEmail(usernameOrEmail)
	} else {
		user, err = UserByUsername(usernameOrEmail)
	}
	return
}

func (user *User) Create() (err error) {
	if existsUsernameNotID(user.Username, user.ID) || existsEmailNotID(user.Email, user.ID) {
		return ErrExistsUsernameOrEmail
	}
	st, err := DB.Prepare("INSERT INTO USERS(UUID, USERNAME, EMAIL, PASSWORD, FIRST_NAME, LAST_NAME, IS_SHOP, CREATED_AT) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING ID, UUID, CREATED_AT")
	if err != nil {
		return
	}
	defer st.Close()
	err = st.QueryRow(
		CreateUUID(), user.Username, user.Email, Encrypt(user.Password), user.FirstName, user.LastName, user.IsShop,
		time.Now(),
	).Scan(
		&user.ID, &user.UUID, &user.CreatedAt,
	)
	return
}

func UserByID(user_id int) (user User, err error) {
	err = DB.QueryRow("SELECT * FROM USERS WHERE ID = $1", user_id).Scan(
		&user.ID, &user.UUID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.IsShop, &user.Photo, &user.LanguageCode,
		&user.Description, &user.CreatedAt,
	)
	return
}

func (user *User) CreateSession() (session Session, err error) {
	stmt, err := DB.Prepare("INSERT INTO SESSIONS (UUID, EMAIL, USER_ID, CREATED_AT) VALUES ($1, $2, $3, $4) RETURNING ID, UUID, EMAIL, USER_ID, CREATED_AT")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(CreateUUID(), user.Email, user.ID, time.Now()).Scan(
		&session.ID, &session.UUID, &session.Email, &session.User_ID, &session.Created_at,
	)
	return
}

func (user *User) Session() (session Session, err error) {
	err = DB.QueryRow(
		"SELECT ID, UUID, EMAIL, USER_ID, CREATED_AT FROM SESSIONS WHERE USER_ID = ?", user.ID,
	).Scan(&session.ID, &session.UUID, &session.Email, &session.User_ID, &session.Created_at)
	return
}

func UserByUsername(username string) (user User, err error) {
	err = DB.QueryRow("SELECT * FROM USERS WHERE USERNAME=$1", username).Scan(
		&user.ID, &user.UUID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.IsShop, &user.Photo, &user.LanguageCode,
		&user.Description, &user.CreatedAt,
	)
	return
}

func UserByEmail(email string) (user User, err error) {
	err = DB.QueryRow("SELECT * FROM USERS WHERE EMAIL = $1", email).Scan(
		&user.ID, &user.UUID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.IsShop, &user.Photo, &user.LanguageCode,
		&user.Description, &user.CreatedAt,
	)
	return
}

func existsUsernameNotID(username string, id int) bool {
	var existsUsername bool
	DB.QueryRow(
		"SELECT EXISTS (SELECT EMAIL FROM USERS WHERE USERNAME=$1 AND ID != $2)", username, id,
	).Scan(&existsUsername)
	return existsUsername
}
func existsEmailNotID(email string, id int) bool {
	var existsEmail bool
	DB.QueryRow(
		"SELECT EXISTS (SELECT EMAIL FROM USERS WHERE EMAIL=$1 AND ID != $2)", email, id,
	).Scan(&existsEmail)
	return existsEmail
}
