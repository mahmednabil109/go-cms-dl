package cms

import "io"

type Parser interface {
	Parse(r io.Reader) error
	GetCourses() ([]string, error)
	GetTitle() (string, error)
	GetWeeks() ([]*Week, error)
}

type Course struct {
	Link string
}

type CmsFile struct {
	Name string
	Link string
}

type Week struct {
	Name  string
	Files []CmsFile
}
