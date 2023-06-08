package dl

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/Azure/go-ntlmssp"
	"github.com/mahmednabil109/go-cms-dl/cms"
	"github.com/mahmednabil109/go-cms-dl/material"
	"github.com/mahmednabil109/go-cms-dl/utils"
	"go.uber.org/zap"
)

const base_url = "https://cms.guc.edu.eg"

// TODO(mahmednabil109): refactor method parameters
func DownloadCourses(
	mail, password, url string,
	cmsParser cms.Parser,
	mManager material.IManager,
	logger *zap.Logger,
) {
	client := &http.Client{
		Transport: ntlmssp.Negotiator{},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Fatal("failed to create url", zap.Error(err))
	}
	req.SetBasicAuth(mail, password)
	res, err := client.Do(req)
	fmt.Println(res.StatusCode)

	if err != nil || res.StatusCode != http.StatusOK {
		logger.Fatal("failled to get page", zap.Error(err))
	}
	defer res.Body.Close()

	err = cmsParser.Parse(res.Body)
	if err != nil {
		logger.Fatal("failed to create doc parser", zap.Error(err))
	}

	courses, _ := cmsParser.GetCourses()
	// TODO(mahmednabil): use `config.concurrent`
	pool := utils.NewPool(1)

	var wait sync.WaitGroup
	for _, courseLink := range courses {
		wait.Add(1)
		courseLink := courseLink
		pool.Submit(func() error {
			defer wait.Done()
			downloadCourse(courseLink, mail, password, mManager, logger)
			return nil
		})
	}
	wait.Wait()
}

func downloadCourse(url, mail, password string, mManager material.IManager, logger *zap.Logger) {
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
	name, err := gq.GetTitle()
	if err != nil {
		logger.Fatal("failed to get course title", zap.Error(err))
	}
	courseId, err := mManager.GetCourseId(name)
	if err != nil {
		logger.Fatal("failed to get course ID:", zap.Error(err))
	}

	weeks, err := gq.GetWeeks()
	if err != nil {
		logger.Fatal("failed to get weeks content", zap.Error(err))
	}

	// TODO(mahmednabil109): again read from the app config
	pool := utils.NewPool(1)
	var files_wg sync.WaitGroup
	for _, week := range weeks {
		fmt.Println("\t", name, week.Name)
		weekId, err := mManager.GetWeekId(courseId, week.Name)
		if err != nil {
			logger.Fatal("failed to get week ID:", zap.Error(err))
		}
		for _, file := range week.Files {
			file := file
			if mManager.FileExists(courseId, file.Name) {
				continue
			}

			files_wg.Add(1)
			fmt.Println("\t\t", name, file.Name)
			pool.Submit(func() error {
				defer files_wg.Done()
				err := downloadFile(weekId, file, mManager)
				if err != nil {
					logger.Fatal(
						"failed to download file",
						zap.String("name", file.Name), zap.Error(err),
					)
				}
				return nil
			})
		}
	}
	files_wg.Wait()
}

func downloadFile(weekId int, file cms.CmsFile, mManager material.IManager) error {
	res, err := http.Get(base_url + file.Link)
	urlComponents := strings.Split(file.Link, ".")
	if len(urlComponents) < 2 {
		return errors.New("can't extract the file extension")
	}
	ext := urlComponents[len(urlComponents)-1]
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return mManager.SaveFile(weekId, file.Name, ext, res.Body)
}
