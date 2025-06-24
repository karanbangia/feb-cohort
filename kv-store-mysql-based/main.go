package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()

	dsn := "dbeaver:dbeaver@tcp(127.0.0.1:3306)/mysql"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	r.GET("/get/:key", func(c *gin.Context) {
		key := c.Param("key")
		var value string
		currentTime := time.Now().Unix()
		err = db.QueryRow("SELECT `VALUE` FROM KV_STORE WHERE `KEY` = ? AND COALESCE(EXPIRED_AT,  '4566546456456') > ?", key, currentTime).Scan(&value)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
	})

	r.PUT("/put", func(c *gin.Context) {
		var kv struct {
			Key   string `json:"key"`
			Value string `json:"value"`
			Ttl   *int   `json:"ttl"`
		}

		if err = c.ShouldBindJSON(&kv); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if kv.Ttl == nil {

			_, err = db.Exec("REPLACE INTO KV_STORE (`KEY`, `VALUE`) VALUES (?, ?)", kv.Key, kv.Value)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			//ttl in milliseconds, convert to epoch
			expiredAt := time.Now().Add(time.Duration(*kv.Ttl) * time.Second).Unix()
			_, err = db.Exec("REPLACE INTO KV_STORE (`KEY`, `VALUE`, `EXPIRED_AT`) VALUES (?, ?, ?)", kv.Key, kv.Value, expiredAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	r.DELETE("/delete/:key", func(c *gin.Context) {
		key := c.Param("key")
		_, err = db.Exec("UPDATE KV_STORE SET EXPIRED_AT=-1 WHERE `KEY` = ?", key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	err = r.Run()
	if err != nil {
		panic(err)
	} // listen and serve on 0.0.0.0:8080
}
