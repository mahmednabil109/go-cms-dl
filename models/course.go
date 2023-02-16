package models

type CmsFile struct {
	Name string
	Link string
}

type week []CmsFile

type Course struct {
	Name  string
	Link  string
	Weeks []week
}
