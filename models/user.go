package models

type User struct {
	Email    string
	Password string
	Role 	 string
}

func (u *User) PutPID(pid string)           { u.Email = pid }
func (u *User) PutEmail(email string)       { u.Email = email }
func (u *User) PutPassword(password string) { u.Password = password }
func (u *User) PutRole(role string)         { u.Role = role }
func (u *User) PutArbitrary(values map[string]string) {
	if n, ok := values["role"]; ok {
		u.Role = n
	}
}

func (u User) GetPID() string                 { return u.Email }
func (u User) GetPassword() (password string) { return u.Password }
func (u User) GetEmail() string               { return u.Email }
func (u User) GetRole() string                { return u.Role }
func (u User) GetArbitrary() map[string]string {
	return map[string]string{
		"role": u.Role,
	}
}

