package main

import (
	"admin-server/models"
	"github.com/robfig/cron"
	"log"
	"time"
)

func main() {
	log.Println("Starting...")

	c := cron.New()

	var tag models.Tag
	c.AddFunc("* * * * * *", func() {
		log.Println("Run tag.CleanAll...")
		tag.CleanAll()
	})
	var article models.Article
	c.AddFunc("* * * * * *", func() {
		log.Println("Run article.CleanAll...")
		article.CleanAll()
	})

	c.Start()

	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}
