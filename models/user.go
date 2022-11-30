package models

import "time"

/* User - User Struct */
type User struct {
	ID        int64     `json:"id"`
	UserName  string    `json:"userName"`
	Pass      string    `json:"pass"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

/*UserResponse - Postgres User response */
type UserResponse struct {
	UserName string    `json:"username"`
	ReqBytes int       `json:"req_bytes" config:"req_bytes" sql:"req_bytes,omitempty" db:"req_bytes"`
	Time     time.Time `json:"time"`
}

/* GetUserId - Get Id from body Request*/
type GetUserId struct {
	ID int64 `json:"id" uri:"id" binding:"required" `
}

/* UserCreateAndUpdate - Create and Update Body Request */
type UserCreateAndUpdate struct {
	UserName string `json:"userName"`
	Pass     string `json:"pass"`
	Email    string `json:"email"`
	ReqBytes int    `json:"reqBytes"`
}
