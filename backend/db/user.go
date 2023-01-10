package db

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID              int    `json:"id,omitempty"`
	UUID            string `json:"UUID,omitempty"`
	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	IsShop          bool   `json:"is_shop"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	Repassword      string `json:"repassword,omitempty"`
	Photo           string `json:"photo,omitempty"`
	LanguageCode    string `json:"languageCode,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	Email           string `json:"email,omitempty"`
	Description     string `json:"description,omitempty"`
	UsernameOrEmail string `json:"username_or_email,omitempty"`
	OldPassword     string `json:"old_password,omitempty"`
	IsActive        bool   `json:"-"`
}

var ErrExistsUsernameOrEmail = errors.New("Already have username or email on other account")

func (user *User) HideInfo() {
	user.UUID = ""
	user.Password = ""
	user.Repassword = ""
	user.OldPassword = ""
	user.Email = ""
}

func (user *User) Update() (err error) {
	stmt, err := DB.Prepare("UPDATE USERS SET FIRST_NAME=$1, lAST_NAME=$2, USERNAME=$3, PASSWORD=$4, EMAIL=$5, DESCRIPTION=$6 WHERE ID=$7")
	if err != nil {
		return
	}
	_, err = stmt.Exec(
		user.FirstName, user.LastName, user.Username, Encrypt(user.Password), user.Email, user.Description, user.ID,
	)
	user.HideInfo()
	return
}

func (user *User) UpdatePhoto() (err error) {
	stmt, err := DB.Prepare("UPDATE USERS SET PHOTO=$1 WHERE ID=$2")
	if err != nil {
		return
	}
	_, err = stmt.Exec(user.Photo, user.ID)
	return
}

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
	err = DB.QueryRow("SELECT * FROM USERS WHERE ID = $1 AND IS_ACTIVE=TRUE", user_id).Scan(
		&user.ID, &user.UUID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.IsShop, &user.Photo, &user.LanguageCode,
		&user.Description, &user.IsActive, &user.CreatedAt,
	)
	return
}

func UserByIDForPublic(user_id int) (user User, err error) {
	err = DB.QueryRow(
		"SELECT USERNAME, FIRST_NAME, LAST_NAME, IS_SHOP, PHOTO, DESCRIPTION FROM USERS WHERE ID = $1 AND IS_ACTIVE=TRUE",
		user_id,
	).Scan(
		&user.Username, &user.FirstName, &user.LastName,
		&user.IsShop, &user.Photo, &user.Description,
	)
	user.ID = user_id
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
	err = DB.QueryRow("SELECT * FROM USERS WHERE USERNAME=$1 AND IS_ACTIVE=TRUE", username).Scan(
		&user.ID, &user.UUID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.IsShop, &user.Photo, &user.LanguageCode,
		&user.Description, &user.IsActive, &user.CreatedAt,
	)
	return
}

func UserByEmail(email string) (user User, err error) {
	err = DB.QueryRow("SELECT * FROM USERS WHERE EMAIL = $1 AND IS_ACTIVE=TRUE", email).Scan(
		&user.ID, &user.UUID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.IsShop, &user.Photo, &user.LanguageCode,
		&user.Description, &user.IsActive, &user.CreatedAt,
	)
	return
}

func existsUsernameNotID(username string, id int) bool {
	var existsUsername bool
	DB.QueryRow(
		"SELECT EXISTS (SELECT EMAIL FROM USERS WHERE USERNAME=$1 AND ID != $2 AND IS_ACTIVE_TRUE)", username, id,
	).Scan(&existsUsername)
	return existsUsername
}
func existsEmailNotID(email string, id int) bool {
	var existsEmail bool
	DB.QueryRow(
		"SELECT EXISTS (SELECT EMAIL FROM USERS WHERE EMAIL=$1 AND ID != $2 AND IS_ACTIVE=TRUE)", email, id,
	).Scan(&existsEmail)
	return existsEmail
}

func (user *User) Deactivate() {
	user.IsActive = false
	DB.Exec("UPDATE USERS SET IS_ACTIVE=FALSE WHERE ID=$1", user.ID)
}
