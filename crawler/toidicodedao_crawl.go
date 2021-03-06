package crawler

import (
	"context"
	"devread/custom_error"
	"devread/helper"
	"devread/model"
	"devread/repository"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"runtime"
	"strings"
	"time"
)

func ToidicodedaoPost(postRepo repository.PostRepo) {
	c := colly.NewCollector()
	c.SetRequestTimeout(30 * time.Second)

	posts := []model.Post{}
	var toidicodedaoPost model.Post

	c.OnHTML("footer[class=entry-meta]", func(e *colly.HTMLElement) {
		if toidicodedaoPost.Name == "" || toidicodedaoPost.Link == "" {
			return
		}
		toidicodedaoPost.Tag = strings.ToLower(e.ChildText("span.tag-links > a:last-child"))
		posts = append(posts, toidicodedaoPost)
	})

	c.OnHTML(".site-content .entry-title", func(e *colly.HTMLElement) {
		toidicodedaoPost.Name = e.Text
		toidicodedaoPost.Link = e.ChildAttr("h1.entry-title > a", "href")
		c.Visit(e.Request.AbsoluteURL(toidicodedaoPost.Link))
		if toidicodedaoPost.Name == "" || toidicodedaoPost.Link == "" {
			return
		}
		toidicodedaoPost.PostID = helper.Hash(toidicodedaoPost.Name, toidicodedaoPost.Link)
		posts = append(posts, toidicodedaoPost)
	})

	c.OnScraped(func(r *colly.Response) {
		queue := helper.NewJobQueue(runtime.NumCPU())
		queue.Start()
		defer queue.Stop()

		for _, post := range posts {
			queue.Submit(&ToidicodedaoProcess{
				post:     post,
				postRepo: postRepo,
			})
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	for i := 1; i < 32; i++ {
		fullURL := fmt.Sprintf("https://toidicodedao.com/category/chuyen-coding/page/%d", i)
		c.Visit(fullURL)
		fmt.Println(fullURL)
	}
}

type ToidicodedaoProcess struct {
	post     model.Post
	postRepo repository.PostRepo
}

func (process *ToidicodedaoProcess) Process() {
	// select post by id
	cacheRepo, err := process.postRepo.SelectById(context.Background(), process.post.PostID)
	if err == custom_error.PostNotFound {
		// insert post to database
		fmt.Println("Add: ", process.post.Name)
		_, err = process.postRepo.Save(context.Background(), process.post)
		if err != nil {
			log.Println(err)
		}
		return
	}

	// update post
	if process.post.PostID != cacheRepo.PostID {
		fmt.Println("Updated: ", process.post.Name)
		_, err = process.postRepo.Update(context.Background(), process.post)
		if err != nil {
			log.Println(err)
		}
	}
}
