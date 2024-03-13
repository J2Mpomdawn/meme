package model

//used for DB connection process
type DB struct {
	Host      string
	Port      string
	User      string
	Password  string
	DBName    string
	Collation string
}

//used for SSH connection process
type SSH struct {
	Key  string
	Host string
	Port string
	User string
}
