package main

import (
	"context"
	"database/sql"
	"demoapi/internal/data"
	"demoapi/internal/jsonlog"
	"encoding/json"
	"fmt"
	"github.com/mxschmitt/playwright-go"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const version = "1.0.0"

type config struct {
	Port int    `json:"port"` //Network Port
	Env  string `json:"env"`  //Current operating environment
	Db   struct {
		Dsn          string `json:"dsn"` //Database connection
		MaxOpenConns int    `json:"maxOpenConns"`
		MaxIdleConns int    `json:"maxIdleConns"`
		MaxIdleTime  string `json:"maxIdleTime"`
	} `json:"db"`
	Limiter struct {
		Rps     float64 `json:"rps"`     //Allowed requests per second
		Burst   int     `json:"burst"`   //Num of  maximum requests in single burst
		Enabled bool    `json:"enabled"` //Is Rate Limiter is on
	} `json:"limiter"`

	// cors struct {
	// 	trustedOrigins []string
	// }
}

type application struct {
	config  config
	logger  *jsonlog.Logger
	browser playwright.Browser
	models  data.Models
	wg      sync.WaitGroup
}

func main() {
	//var cfg config
	conf, err := os.Open("./cmd/configs/config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer conf.Close()

	byteValue, _ := ioutil.ReadAll(conf)

	var configs config
	err = json.Unmarshal(byteValue, &configs)
	if err != nil {
		fmt.Println(err)
		return
	}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	/*db, err := openDB(configs)
	if err != nil {
		logger.PrintFatal(err, nil)
	}*/

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.WebKit.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}

	app := &application{
		config:  configs,
		logger:  logger,
		browser: browser,
		models:  data.NewModels(),
	}

	err = app.Serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Db.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Db.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.Db.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	//if err = db.Ping(); err != nil {
	//	return nil, err
	//}
	return db, nil
}
