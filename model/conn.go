package model

//used for each connection process

type SSH struct {
	Key  string
	Host string
	Port string
	User string
}

type DB struct {
	Host      string
	Port      string
	User      string
	Password  string
	DBName    string
	Collation string
}
