package material

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
)

const NOTFOUND = -1

const init_db = `
	CREATE TABLE IF NOT EXISTS courses (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(25)
	);
	CREATE TABLE IF NOT EXISTS weeks (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		courseId INTEGER,
		name VARCHAR(25),
		FOREIGN KEY(courseId) REFERENCES courses(id)
	);
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		weekId INTEGER,
		name VARCHAR(25),
		path TEXT,
		FOREIGN KEY(weekId) REFERENCES weeks(id)
	);
`

type Metadb struct {
	mux sync.Mutex
	db  *sql.DB
}

func NewDB(file string) (*Metadb, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(init_db); err != nil {
		return nil, err
	}

	return &Metadb{
		db: db,
	}, nil
}

func (m *Metadb) GetCourseIdByName(name string) (int, error) {
	row, err := m.db.Query("SELECT id FROM courses WHERE name=?", name)
	if err != nil {
		return -1, err
	}

	var id int
	if row.Next() {
		err = row.Scan(&id)
	}
	return id, err
}

func (m *Metadb) InsertCourse(name string, date time.Time) (int, error) {
	res, err := m.db.Exec("INSERT INTO courses VALUES(?)", name)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	return int(id), err
}

func (m *Metadb) GetWeekIdByName(name string) (int, error) {
	row, err := m.db.Query("SELECT id FROM weeks WHERE name=?", name)
	if err != nil {
		return -1, err
	}

	var id int
	if row.Next() {
		err = row.Scan(&id)
	}
	return id, err
}

func (m *Metadb) InsertWeek(courseId int, name string, date time.Time) (int, error) {
	res, err := m.db.Exec("INSERT INTO weeks(name, courseId) VALUES(?, ?)", name, courseId)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	return int(id), err
}

func (m *Metadb) GetFileIdByName(name string, date time.Time) (int, error) {
	row, err := m.db.Query("SELECT id FROM files WHERE name=?", name)
	if err != nil {
		return -1, err
	}

	var id int
	if row.Next() {
		err = row.Scan(&id)
	}
	return id, err
}

func (m *Metadb) InsertFile(weekId int, name, path string) (int, error) {
	res, err := m.db.Exec("INSERT INTO files(name, weekId, path) VALUES(?, ?, ?);", name, weekId, path)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	return int(id), err
}

func (m *Metadb) GetAllFileNames() ([]string, error) {
	var files []string
	rows, err := m.db.Query(`
		SELECT f.name, c.id FROM files f
			INNER JOIN weeks w ON f.weekId = w.id 
			INNER JOIN courses c ON w.courseId = c.id;
	`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			name     string
			courseId int
		)
		if err := rows.Scan(&name, &courseId); err != nil {
			return nil, err
		}
		files = append(files, fmt.Sprintf("%d-%s", courseId, name))
	}

	if err := rows.Err(); err != nil {
		return files, err
	}

	return files, nil
}

func (m *Metadb) Close() error {
	return m.db.Close()
}
