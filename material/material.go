package material

import (
	"fmt"
	"io"
	"os"
	"path"
	"sync"
	"time"
)

type IManager interface {
	GetCourseId(name string) (int, error)
	GetWeekId(courseId int, name string) (int, error)
	FileExists(courseId int, name string) bool
	SaveFile(weekId int, name string, data io.Reader) error
}

type Manager struct {
	metadb        *Metadb
	downloadsPath string
	fileNames     map[string]struct{}
	coursePaths   map[int]string
	weekPaths     map[int]string
	mux           sync.Mutex
}

func NewManager(downloadsPath, dbPath string) (*Manager, error) {
	db, err := NewDB(dbPath)
	if err != nil {
		return nil, err
	}

	// hold hashtable of file names
	files, err := db.GetAllFileNames()
	if err != nil {
		db.Close()
		return nil, err
	}
	manager := Manager{
		metadb:        db,
		downloadsPath: downloadsPath,
		fileNames:     make(map[string]struct{}),
		coursePaths:   make(map[int]string),
		weekPaths:     make(map[int]string),
	}
	for _, name := range files {
		manager.fileNames[name] = struct{}{}
	}

	return &manager, nil
}

func (m *Manager) GetCourseId(name string) (int, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	id, err := m.metadb.GetCourseIdByName(name)
	if err != nil {
		return -1, err
	}
	if id == NOTFOUND {
		id, err = m.metadb.InsertCourse(name, time.Now())
		if err != nil {
			return -1, err
		}
	}

	coursePath := path.Join(m.downloadsPath, name)
	err = os.Mkdir(
		coursePath,
		0777,
	)
	if err != nil && !os.IsExist(err) {
		return -1, nil
	}
	m.coursePaths[id] = coursePath
	return id, nil
}

func (m *Manager) GetWeekId(courseId int, name string) (int, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	id, err := m.metadb.GetWeekIdByName(name)
	if err != nil {
		return -1, err
	}
	if id == NOTFOUND {
		id, err = m.metadb.InsertWeek(courseId, name, time.Now())
		if err != nil {
			return -1, err
		}
	}

	coursePath, ok := m.coursePaths[courseId]
	if !ok {
		fmt.Printf("course path not existing %v\n", courseId)
	}
	weekPath := path.Join(coursePath, name)
	err = os.Mkdir(weekPath, 0777)
	if err != nil && !os.IsExist(err) {
		return -1, err
	}
	m.weekPaths[id] = weekPath
	return id, nil
}

func (m *Manager) FileExists(courseId int, name string) bool {
	// THIS IS UGLY :/
	_, ok := m.fileNames[fmt.Sprintf("%d-%s", courseId, name)]
	return ok
}

func (m *Manager) SaveFile(weekId int, name string, data io.Reader) error {
	weekPath, ok := m.weekPaths[weekId]
	if !ok {
		fmt.Printf("week path is not found %v\n", weekId)
	}
	filePath := path.Join(weekPath, name)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, data)
	if err != nil {
		return err
	}

	_, err = m.metadb.InsertFile(weekId, name, filePath)
	return err
}
