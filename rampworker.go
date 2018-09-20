package events

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type work func() (int, error)
type workerSet struct {
	name         string
	max          int
	cur          int
	triggerDelay time.Duration
	work         work
}

var wm sync.Mutex

func NewWorkerSet(name string, max, trig int, f func() (int, error)) *workerSet {
	var ws workerSet
	ws.name = name
	ws.max = max
	ws.triggerDelay = time.Duration(trig)
	ws.work = f
	go ws.mainWorker()
	return &ws
}

func (ws *workerSet) mainWorker() {
	log.Infof("Starting main worker for >%s<", ws.name)
	for {
		//istart := time.Now()
		if w, e := ws.work(); e != nil {
			log.Errorf("[%s] worker@main returns error: %s", ws.name, e)
		} else if w > 0 && ws.cur <= ws.max {
			log.Warnf("Adding worker(s) when ret: %d. CurrThreadCount:%d", w, ws.cur)
			go ws.addWorker()
			if w > 10 {
				go ws.addWorker()
			}
		}
	}
}

func (ws *workerSet) addWorker() {
	wm.Lock()
	ws.cur++
	wm.Unlock()
	for {
		w, e := ws.work()
		if e != nil {
			log.Errorf("[%s] worker returns error: %s", ws.name, e)
		}
		if w == 0 {
			log.Warnf("[%s] Worker going away work ret[0], currThreadCount: %d", ws.name, ws.cur)
			break
		}
	}
	wm.Lock()
	ws.cur--
	wm.Unlock()
}
