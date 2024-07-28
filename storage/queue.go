package storage

import "sync"

var SongQueue Queue

type Queue struct {
	Songs []Song
	sync.Mutex
}

func (q *Queue) Add(song Song) {
	q.Lock()
	defer q.Unlock()
	q.Songs = append(q.Songs, song)
}

func (q *Queue) Peek() *Song {
	q.Lock()
	defer q.Unlock()
	if len(q.Songs) == 0 {
		return nil
	}
	return &q.Songs[0]
}

func (q *Queue) Remove() {
	q.Lock()
	defer q.Unlock()
	if len(q.Songs) > 0 {
		q.Songs = q.Songs[1:]
	}
}

func (q *Queue) Clear() {
	q.Lock()
	defer q.Unlock()
	q.Songs = []Song{}
}

func GetSongQueue() *Queue {
	return &SongQueue
}
