package db

import (
	"crypto/rand"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int 		 
	UserName  string `binding:"required,min=5,max=30"`
	Password  string `binding:"required,min=6,max=30"`
	Email     string `binding:"required,min=5,max=100"`
	CreatedAt time.Time 
	UpdatedAt time.Time 
	Salt      []byte `json:"-"`
	HashedPassword []byte `json:"-"`
}





func (u *User)AddUser() (User, error){
	salt, err := GenerateSalt()
	if err != nil {
    return User{}, err
  }
	toHash := append([]byte(u.Password), salt...)
	hashedPassword, err := bcrypt.GenerateFromPassword(toHash, bcrypt.DefaultCost)
	if err != nil {
    return User{}, err
  }
	u.Salt = salt
  u.HashedPassword = hashedPassword
	cmd := `insert into users (
		username, 
		email, 
		hashedpassword,
		salt) values (?, ?, ?, ?)`
	_, err = Db.Exec(cmd, u.UserName, u.Email, u.HashedPassword, u.Salt)
	if err != nil {
		fmt.Printf("ユーザー追加時にエラーが起きました: %v\n", err)
		return User{},err
	}

	cmd = `select id, username, email, hashedpassword, salt from users
	where email = ?`
	user := User{}
	err = Db.QueryRow(cmd, u.Email).Scan(
		&user.Id,
		&user.UserName,
		&user.Email,
		&user.HashedPassword,
		&user.Salt,)
	if err != nil {
		fmt.Printf("スキャン時にエラーが起きました: %v\n", err)
		return User{}, err
	}	

	return user, err
}

func (u *User) Signin() (User, error) {
	cmd := `select id, username, email, hashedpassword, salt from users
	where email = ?`
	user := User{}
	err := Db.QueryRow(cmd, u.Email).Scan(
		&user.Id,
		&user.UserName,
		&user.Email,
		&user.HashedPassword,
		&user.Salt,)
	if err != nil {
		fmt.Printf("スキャン時にエラーが起きました: %v\n", err)
		return User{}, err
	}
	if err := bcrypt.CompareHashAndPassword(user.HashedPassword, append([]byte(u.Password), user.Salt...)); err != nil {
		fmt.Printf("パスワードが一致しません: %v\n", err)
		return User{}, err
	}
	return user, err
}


func FetchUser(id int64) (User, error) {
	cmd := `select id, username, email, hashedpassword, salt from users
	where id = ?`
	user := User{}
	err := Db.QueryRow(cmd, id).Scan(
		&user.Id,
		&user.UserName,
		&user.Email,
		&user.HashedPassword,
		&user.Salt,)
	if err != nil {
		fmt.Printf("スキャン時にエラーが起きました: %v\n", err)
		return User{}, err
	}
	return user, err
}

func GenerateSalt() ([]byte, error) {
  salt := make([]byte, 16)
  if _, err := rand.Read(salt); err != nil {
		fmt.Printf("salt作成時にエラーが起きました: %v\n", err)
    return nil, err
  }
  return salt, nil
}






