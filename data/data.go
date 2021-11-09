package data

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var dbCon *sql.DB
var (
	DataUserNotFound     = errors.New("User Not Found")
	DataPasswordNotMatch = errors.New("Password Not Match")
)

func init() {
	var err error
	log.SetFlags(log.Lshortfile | log.Ltime)

	err = godotenv.Load(".env")

	dburl := fmt.Sprintf("host=shorturl_db port=5432 user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	if err != nil {
		log.Fatal("Could not load the .env File ", err.Error())
	}

	dbCon, err = sql.Open("postgres", dburl)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = dbCon.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
}

// USER DATABAES OPERRATIONS
//Insert User Data
func (u *User) InsertUser() error {

	stm, err := dbCon.Prepare("INSERT INTO users(name,password,email) VALUES($1,$2,$3)")

	pass := hashPassword(u.Password)
	res, err := stm.Exec(u.Name, pass, u.Email)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	lastId, _ := res.LastInsertId()
	rowcount, _ := res.RowsAffected()

	if err != nil {
		log.Println(err.Error())
	} else {
		log.Printf("INSERT USESRS ID:=%d , affected:=%d", lastId, rowcount)
	}

	return nil
}

//Get User Data by Email
func AuthUser(email, pass string) (int, error) {
	User, err := GetUser(email)
	if err != nil {
		return 0, err
	}
	if checkPasswordHash(pass, User.Password) {
		return User.ID, nil
	}
	return -1, DataPasswordNotMatch
}
func GetUser(email string) (User, error) {
	u := User{}
	row := dbCon.QueryRow("SELECT * FROM users WHERE email= $1;", email)
	err := row.Scan(&u.ID, &u.Name, &u.Password, &u.Email)
	if err != nil {
		log.Println(err.Error())
		return User{}, DataUserNotFound
	}
	u.UrlServices = GetUrlsByUserID(u.ID)
	return u, nil
}
func GetUserById(id int) User {
	u := User{}
	row := dbCon.QueryRow("SELECT * FROM users WHERE user_id= $1;", id)
	err := row.Scan(&u.ID, &u.Name, &u.Password, &u.Email)
	if err != nil {
		log.Println(err.Error())
		return User{}
	}
	u.UrlServices = GetUrlsByUserID(u.ID)
	return u
}

//-----------------------------------------------------------------------------------
// URL_SERVICE DATABAES OPERRATIONS
func InsertSevice(url string, userid int) {
	stm, err := dbCon.Prepare("INSERT INTO urlservice(url,code,user_id) VALUES($1,$2,$3)")
	if err != nil {
		log.Println(err.Error())
	}

	code := createServiceCode(url)
	res, err := stm.Exec(url, code, userid)
	if err != nil {
		log.Println(err.Error())
	} else {
		lastId, _ := res.LastInsertId()
		rowcount, _ := res.RowsAffected()
		log.Printf("ID:=%d , affected:=%d", lastId, rowcount)
	}
}

func GetUrl(code string) UrlService {
	url := UrlService{}
	row := dbCon.QueryRow("SELECT urlservice_id,url,code FROM urlservice WHERE code= $1", code)
	err := row.Scan(&url.ID, &url.Url, &url.Code)
	if err != nil {
		log.Println(err.Error())
		return UrlService{}
	}
	return url
}

func GetUrlsByUserID(userid int) []UrlService {

	var urls []UrlService
	rows, err := dbCon.Query("SELECT urlservice_id,url,code FROM urlservice WHERE user_id= $1", userid)
	if err != nil {
		log.Println(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		temp := UrlService{}
		err := rows.Scan(&temp.ID, &temp.Url, &temp.Code)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		urls = append(urls, temp)
	}
	return urls
}

func DeleteServiceByCode(code string) {
	rs, err := dbCon.Exec("DELETE FROM urlservice WHERE code=$1", code)
	if err != nil {
		log.Println(err.Error())
	}
	r, _ := rs.LastInsertId()
	log.Println("DELETE service Rows Affected = ", r)
}

//------------------------------------------------------------------------
//UTILITY FUNCTIONS
func hashPassword(password string) string {
	data, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Println(err.Error())
	}
	return string(data)
}
func checkPasswordHash(pass, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

func createServiceCode(name string) string {

	ha := sha512.New()
	ha.Write([]byte(name))
	code := hex.EncodeToString(ha.Sum(nil))

	return code[:8]
}

//------------------------------------------------------------------------------
