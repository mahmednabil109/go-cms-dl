package main

import (
	"fmt"
	"os"

	"github.com/mahmednabil109/go-cms-dl/cms"
	"github.com/mahmednabil109/go-cms-dl/dl"
	"github.com/mahmednabil109/go-cms-dl/utils"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment(
		zap.WithCaller(false),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Err while creating the logger: %v", err)
		os.Exit(-1)
	}
	defer logger.Sync()

	cfg, err := Parse(logger)
	if err != nil {
		logger.Panic("", zap.Error(err))
	}

	dotfile, err := utils.NewDotfile(cfg.DotFilePath)
	if err != nil {
		logger.Fatal("failed to create dotfile", zap.Error(err))
	}

	url, user, password := cfg.CMS_URL, dotfile.Get("MAIL"), dotfile.Get("PASS")
	cmsParser := cms.NewParser()
	dl.DownloadCourses(
		user,
		password,
		url,
		cmsParser,
		logger,
	)
}
