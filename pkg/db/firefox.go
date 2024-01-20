package db

import (
	"database/sql"
	"os"
	"os/user"
	"path"
	"resin/pkg/logging"

	_ "github.com/mattn/go-sqlite3"
)

func singleQuery(db *sql.DB, name string) (string, error) {
	var err error
	const query string = `select value from moz_cookies where host = '.hoyolab.com' and name = ?;`
	row, err := db.Query(query, name)
	if err != nil {
		logging.Warn("Failed to query firefox db: %s", err)
		return "", err
	}
	defer row.Close()
	var value *string
	row.Next()
	err = row.Scan(&value)
	if err != nil {
		logging.Warn("Failed to scan firefox db: %s", err)
		return "", err
	}
	if value != nil {
		return *value, nil
	}
	return "nil", nil
}

func queryCookies(file string) (*Cookies, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		logging.Warn("Failed to open firefox db")
		return nil, err
	}
	defer db.Close()
	ltoken, err := singleQuery(db, "ltoken_v2")
	if err != nil {
		return nil, err
	}
	ltuid, err := singleQuery(db, "ltuid_v2")
	if err != nil {
		return nil, err
	}
	return &Cookies{
		Ltoken_v2: ltoken,
		Ltuid_v2:  ltuid,
	}, nil
}

func ReadFirefoxCookies() *Cookies {
	u, err := user.Current()
	if err != nil {
		logging.Warn("Unable to get home dir of user")
		return nil
	}
	profiles := path.Join(u.HomeDir, "AppData\\Roaming\\Mozilla\\Firefox\\Profiles")
	entries, err := os.ReadDir(profiles)
	for _, e := range entries {
		fullpath := path.Join(profiles, e.Name(), "cookies.sqlite")
		if _, err := os.Stat(fullpath); os.IsNotExist(err) {
			continue
		}
		if cookies, err := queryCookies(fullpath); err == nil {
			return cookies
		}
	}
	return nil
}
