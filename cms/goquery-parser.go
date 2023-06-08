package cms

import (
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const COURSE_BASE_URL = "https://cms.guc.edu.eg/apps/student/CourseViewStn"

type Gquery struct {
	doc *goquery.Document
}

func NewParser() *Gquery {
	return &Gquery{}
}

func (g *Gquery) Parse(r io.Reader) error {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}
	g.doc = doc
	return nil
}

func (g *Gquery) GetCourses() ([]string, error) {
	courses := make([]string, 0, 5)
	if g.doc == nil {
		return courses, nil
	}
	g.doc.Find("input[value=\"View Course\"]").Each(func(i int, s *goquery.Selection) {
		siblings := s.Parent().Parent().Children()

		sid := siblings.Get(siblings.Length() - 1).FirstChild.Data
		cid := siblings.Get(siblings.Length() - 2).FirstChild.Data
		courses = append(courses,
			COURSE_BASE_URL+"?sid="+sid+"&id="+cid,
		)
	})
	return courses, nil
}

func (g *Gquery) GetTitle() (string, error) {
	title := ""
	if g.doc == nil {
		return title, nil
	}

	g.doc.Find("span#ContentPlaceHolderright_ContentPlaceHoldercontent_LabelCourseName").
		Each(func(i int, s *goquery.Selection) {
			// TODO(mahmednabil109): trim the name and replace any special characters
			title = strings.Replace(s.Text(), "/", "", -1)
		})
	return title, nil
}

func (g *Gquery) GetWeeks() ([]*Week, error) {
	weeks := make([]*Week, 0, 12)
	if g.doc == nil {
		return weeks, nil
	}

	g.doc.Find(".weeksdata").Each(func(i int, s *goquery.Selection) {
		week := Week{
			Files: make([]CmsFile, 0),
		}

		s.Find("strong").FilterFunction(func(_ int, s *goquery.Selection) bool {
			return strings.HasPrefix(s.Text(), "Description")
		}).Each(func(i int, s *goquery.Selection) {
			week.Name = s.Parent().Parent().Find("p").Text()
			week.Name = strings.Replace(week.Name, "/", "", -1)
		})

		s.Find(".card-body strong").Each(func(_ int, s *goquery.Selection) {
			week.Files = append(week.Files, CmsFile{Name: s.Text()})
		})

		s.Find("#download").Each(func(i int, s *goquery.Selection) {
			week.Files[i].Link, _ = s.Attr("href")
		})

		weeks = append(weeks, &week)
	})
	return weeks, nil
}
