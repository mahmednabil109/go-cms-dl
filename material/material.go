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
	SaveFile(weekId int, dirPath, name string, data io.Reader) error
}

type Manager struct {
	metadb        *Metadb
	downloadsPath string
	fileNames     map[string]struct{}
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

	err = os.Mkdir(name, 0777)
	if err != nil && !os.IsExist(err) {
		return -1, nil
	}
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

	err = os.Mkdir(name, 0777)
	if err != nil && !os.IsExist(err) {
		return -1, err
	}

	return id, nil
}

func (m *Manager) FileExists(courseId int, name string) bool {
	// THIS IS UGLY :/
	_, ok := m.fileNames[fmt.Sprintf("%d-%s", courseId, name)]
	return ok
}

func (m *Manager) SaveFile(weekId int, dirPath, name string, data io.Reader) error {
	filePath := path.Join(dirPath, name)
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
