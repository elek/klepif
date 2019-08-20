package period

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

type LastRun struct {
	now       time.Time
	Time      time.Time
	CacheFile string
}

func (lastRun *LastRun) Init() {
	lastRun.now = time.Now()
	content, err := ioutil.ReadFile(lastRun.CacheFile)
	if err != nil {
		panic(err)
	}
	epoch, err := strconv.ParseInt(string(content), 10, 64)
	if err != nil {
		panic(err)
	}
	lastRun.Time = time.Unix(epoch, 0)
}

func (lastRun *LastRun) Before(t time.Time) bool {
	return lastRun.Time.Before(t)
}

func (lastRun *LastRun) Write() error {
	lastTimeStr := strconv.FormatInt(lastRun.now.Unix(), 10)
	return ioutil.WriteFile(lastRun.CacheFile, []byte(lastTimeStr), 777)
}

func (lastRun *LastRun) Print() {
	fmt.Println(lastRun.Time)
}
