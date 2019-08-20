package persistence

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

type Persistence interface {
	Write(key string, value string) error
	Read(key string) (string, error)
	Init() error
}

type Dir struct {
	Path string
}

func (dir *Dir) Init() error {
	return os.MkdirAll(dir.Path, os.ModePerm);
}

func (dir *Dir) Write(key string, value string) error {
	return ioutil.WriteFile(dir.getFilePath(key), []byte(value), os.ModePerm);
}

func (dir *Dir) Read(key string) (string, error) {
	file_path := dir.getFilePath(key);
	if _, err := os.Stat(file_path); os.IsNotExist(err) {
		return "", nil
	}
	data, err := ioutil.ReadFile(file_path);
	return string(data), err
}

func (dir *Dir) getFilePath(key string) string {
	file_path := path.Join(dir.Path, key)
	os.MkdirAll(filepath.Dir(file_path), os.ModePerm);
	return file_path;
}

func ReadTime(persistence Persistence, key string) (time.Time, error) {
	value, err := persistence.Read(key)
	if err != nil {
		return time.Now(), err
	}
	if value == "" {
		return time.Unix(0, 0), nil
	}
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return time.Now(), err
	}
	return time.Unix(i, 0), nil

}

func WriteTime(persistence Persistence, key string, value time.Time) error {
	output := fmt.Sprintf("%d", value.Unix());
	return persistence.Write(key, output);
}
