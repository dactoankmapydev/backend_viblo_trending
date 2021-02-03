package crawler

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"ioc/helper"
	"ioc/model"
	"log"
	"runtime"
	"strings"
)

const urlBase = "https://quan-cam.com"

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func getOnePage(pathURL string) ([]model.Post, error) {
	response, err := helper.HttpClient.GetRequestWithRetries(pathURL)
	checkError(err)
	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	checkError(err)
	posts := make([]model.Post, 0)
	doc.Find("div[class=post]").Each(func(i int, s *goquery.Selection) {
		var quancamPost model.Post
		quancamPost.Name = s.Find("h3.post__title > a").Text()
		link, _ := s.Find("h3.post__title > a").Attr("href")
		quancamPost.Link = urlBase + link
		quancamPost.Tag = strings.ToLower(strings.Replace(
			strings.Replace(
				s.Find("span.tagging > a").Text(), "\n", "", -1), "#", " ", -1))

		//post := model.Post {
		//	Name:   quancamPost.Name,
		//	Tag: quancamPost.Tag,
		//}
		//fmt.Println(post)
		fmt.Println("Name", quancamPost.Name)
		fmt.Println("Link", quancamPost.Link)
		fmt.Println("Tag", quancamPost.Tag)
		fmt.Println("\n ")
		posts = append(posts, quancamPost)
	})
	return posts, nil
}

func AllPage() {
	sem := semaphore.NewWeighted(int64(runtime.NumCPU()))
	group, ctx := errgroup.WithContext(context.Background())

	for page := 1; page <= 5; page++ {
		pathURL := fmt.Sprintf("%s/posts?page=%d", urlBase ,page)
		err := sem.Acquire(ctx, 1)
		if err != nil {
			fmt.Printf("Acquire err = %+v\n", err)
			continue
		}
		group.Go(func() error {
			defer sem.Release(1)

			//do work
			_, err := getOnePage(pathURL)
			checkError(err)

			//queue := NewJobQueue(runtime.NumCPU())
			//queue.Start()
			//defer queue.Stop()
			//
			//for _, post := range posts {
			//	queue.Submit(&QuancamProcess{
			//		post:     post,
			//		postRepo: postRepo,
			//	})
			//}

			return nil
		})
	}
	if err := group.Wait(); err != nil {
		fmt.Printf("g.Wait() err = %+v\n", err)
	}
}

//type QuancamProcess struct {
//	posts []model.Post
//	//iocRepo  repository.IocRepo
//}
//
//func (process *QuancamProcess) Process() {}


// import (
// 	"context"
// 	"fmt"
// 	"github.com/PuerkitoBio/goquery"
// 	"golang.org/x/sync/errgroup"
// 	"golang.org/x/sync/semaphore"
// 	"ioc/helper"
// 	"ioc/model"
// 	"log"
// 	"runtime"
// 	"strconv"
// 	"strings"
// )

// const urlBase = "https://quan-cam.com"

// func checkError(err error) {
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func GetListPage() []string {
// 	pageList := make([]string, 0)
// 	page := []int{1}
// 	for len(page) > 0 {
// 		pathURL := fmt.Sprintf("https://quan-cam.com/posts?page=%d", page[0])
// 		fmt.Println("GetListPage")
// 		response, err := helper.HttpClient.GetRequestWithRetries(pathURL)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		defer response.Body.Close()
// 		doc, err := goquery.NewDocumentFromReader(response.Body)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		link, _ := doc.Find("a.next").Attr("href")
// 		if link != "" {
// 			split := strings.Split(link, "=")[1]
// 			nextLink, _ := strconv.Atoi(split)
// 			page[0] = nextLink
// 			url := fmt.Sprintf("https://quan-cam.com/posts?page=%d", nextLink)
// 			pageList = append(pageList, url)
// 		} else {
// 			page = page[:0]
// 		}
// 	}
// 	fmt.Println("list page->", pageList)
// 	return pageList
// }

// func getOnePage(pathURL string) ([]model.Post, error) {
// 	response, err := helper.HttpClient.GetRequestWithRetries(pathURL)
// 	checkError(err)
// 	defer response.Body.Close()
// 	doc, err := goquery.NewDocumentFromReader(response.Body)
// 	checkError(err)

// 	posts := make([]model.Post, 0)
// 	doc.Find("div[class=post]").Each(func(i int, s *goquery.Selection) {
// 		var quancamPost model.Post
// 		quancamPost.Name = s.Find("h3.post__title > a").Text()
// 		link, _ := s.Find("h3.post__title > a").Attr("href")
// 		quancamPost.Link = urlBase + link
// 		quancamPost.Tag = strings.ToLower(strings.Replace(
// 			strings.Replace(
// 				s.Find("span.tagging > a").Text(), "\n", "", -1), "#", " ", -1))
// 		posts = append(posts, quancamPost)

// 		//fmt.Println("Name", quancamPost.Name)
// 		//fmt.Println("Link", quancamPost.Link)
// 		//fmt.Println("Tag", quancamPost.Tag)
// 		//fmt.Println("\n ")
// 	})
// 	return posts, nil
// }

// func AllPage() {
// 	sem := semaphore.NewWeighted(int64(runtime.NumCPU()))
// 	group, ctx := errgroup.WithContext(context.Background())
// 	listPage := GetListPage()

// 	for _, page := range listPage {
// 		err := sem.Acquire(ctx, 1)
// 		if err != nil {
// 			fmt.Printf("Acquire err = %+v\n", err)
// 			continue
// 		}
// 		group.Go(func() error {
// 			defer sem.Release(1)

// 			//do work
// 			_, err := getOnePage(page)
// 			checkError(err)
			
// 			//queue := NewJobQueue(runtime.NumCPU())
// 			//queue.Start()
// 			//defer queue.Stop()
// 			//
// 			//for _, post := range posts {
// 			//	queue.Submit(&QuancamProcess{
// 			//		post:     post,
// 			//		postRepo: postRepo,
// 			//	})
// 			//}

// 			return nil
// 		})
// 	}
// 	if err := group.Wait(); err != nil {
// 		fmt.Printf("g.Wait() err = %+v\n", err)
// 	}
// }

// //type QuancamProcess struct {
// //	posts []model.Post
// //	//iocRepo  repository.IocRepo
// //}
// //
// //func (process *QuancamProcess) Process() {}
