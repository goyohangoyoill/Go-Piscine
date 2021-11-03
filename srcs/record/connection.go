package record

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var DB *sql.DB

func preparation() (info [5]string) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./secret")
	viper.AddConfigPath("../secret")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Viper Loading Failed")
	}
	info[0] = (viper.Get("DB_KIND")).(string)
	info[1] = os.Getenv((viper.Get("DB_HOST")).(string))
	info[2] = (viper.Get("DB_NAME")).(string)
	info[3] = (viper.Get("DB_USER")).(string)
	info[4] = (viper.Get("DB_PASS")).(string)
	return
}

func Connection() error {
	info := preparation()
	for _, v := range info {
		if v == "" {
			return fmt.Errorf("Viper Invalid Value")
		}
	}
	literal := info[3] + ":" + info[4] + "@tcp(" + info[1] + ")/" + info[2] + "?charset=utf8&parseTime=True&loc=Local"
	if pool, err := sql.Open(info[0], literal); err != nil {
		return err
	} else {
		DB = pool
	}
	_, _ = DB.Exec(`CREATE TABLE IF NOT EXISTS people (
		id int NOT NULL AUTO_INCREMENT PRIMARY KEY UNIQUE,
		name varchar(30) NOT NULL UNIQUE,
		course int DEFAULT 0,
		point int DEFAULT 5,
		password varchar(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at DATETIME
		) ;`)
	_, _ = DB.Exec(`CREATE TABLE IF NOT EXISTS evaluation (
		id int NOT NULL AUTO_INCREMENT PRIMARY KEY UNIQUE,
		interviewee_id int NOT NULL,
		interviewer_id int NOT NULL,
		course int NOT NULL,
		score int NOT NULL,
		pass int DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at date,
		FOREIGN KEY (interviewee_id) REFERENCES people(id),
		FOREIGN KEY (interviewer_id) REFERENCES people(id)
		) ;`)
	return nil
}
