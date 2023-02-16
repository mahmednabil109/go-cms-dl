package dl

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/Azure/go-ntlmssp"
	"github.com/mahmednabil109/go-cms-dl/cms"
	"github.com/mahmednabil109/go-cms-dl/utils"
	"go.uber.org/zap"
)

// TODO(mahmednabil109): refactor method parameters
func DownloadCourses(mail, password, url string, cmsParser cms.Parser, logger *zap.Logger) {
	client := &http.Client{
		Transport: ntlmssp.Negotiator{},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Fatal("failed to create url", zap.Error(err))
	}
	req.SetBasicAuth(mail, password)
	res, err := client.Do(req)
	if err != nil {
		logger.Fatal("failled to get page", zap.Error(err))
	}
	defer res.Body.Close()

	err = cmsParser.Parse(res.Body)
	if err != nil {
		logger.Fatal("failed to create doc parser", zap.Error(err))
	}

	courses, _ := cmsParser.GetCourses()
	// TODO(mahmednabil): use `config.cocurrent`
	pool := utils.NewPool(1)

	var wait sync.WaitGroup
	for _, courseLink := range courses {
		wait.Add(1)
		courseLink := courseLink
		pool.Submit(func() error {
			defer wait.Done()
			downloadCourse(courseLink, mail, password, logger)
			return nil
		})
	}
	wait.Wait()
}

func downloadCourse(url, mail, password string, logger *zap.Logger) {
	res, err := http.Get(url)
	if err != nil {
		logger.Fatal("failed to load course page", zap.Error(err))
	}
	defer res.Body.Close()
	fmt.Println(url, res.StatusCode)

	gq := cms.NewParser()
	err = gq.Parse(res.Body)
	if err != nil {
		logger.Fatal("failed to create doc parser", zap.Error(err))
	}
	name, _ := gq.GetTitle()
	weeks, _ := gq.GetWeeks()

	for _, week := range weeks {
		fmt.Println("\t", name, week.Name)
		for _, file := range week.Files {
			fmt.Println("\t\t", name, file.Name)
		}
	}
}
