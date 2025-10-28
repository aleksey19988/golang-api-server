package repository

import (
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
)

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	DSN      string
}

func NewDB(host, port, user, password, DBName, SSLMode string) *Database {
	return &Database{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   DBName,
		SSLMode:  SSLMode,
	}
}

func (db *Database) ValidateConfig() error {
	if db.Host == "" {
		return errors.New("host is required")
	}

	if db.Port == "" {
		return errors.New("port is required")
	}

	if db.User == "" {
		return errors.New("user is required")
	}

	if db.Password == "" {
		return errors.New("password is required")
	}

	if db.DBName == "" {
		return errors.New("DBName is required")
	}

	if db.SSLMode == "" {
		return errors.New("SSLMode is required")
	}

	return nil
}

func (db *Database) Connect() (*gorm.DB, error) {
	dsn := db.BuildDsn()
	gormDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return gormDb, nil
}

func (db *Database) BuildDsn() string {
	b := strings.Builder{}
	b.WriteString("host=")
	b.WriteString(db.Host)
	b.WriteString(" ")

	b.WriteString("port=")
	b.WriteString(db.Port)
	b.WriteString(" ")

	b.WriteString("user=")
	b.WriteString(db.User)
	b.WriteString(" ")

	b.WriteString("password=")
	b.WriteString(db.Password)
	b.WriteString(" ")

	b.WriteString("dbname=")
	b.WriteString(db.DBName)
	b.WriteString(" ")

	b.WriteString("sslmode=")
	b.WriteString(db.SSLMode)

	db.DSN = b.String()
	return db.DSN
}
