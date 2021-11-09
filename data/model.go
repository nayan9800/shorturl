package data

type User struct {
	ID          int
	Name        string
	Email       string
	Password    string
	UrlServices []UrlService
}
type UrlService struct {
	ID   int
	Url  string
	Code string
}
