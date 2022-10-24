package users

type UserEntry struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Teams []string `json:"teams"`
}
