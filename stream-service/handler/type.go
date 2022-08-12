package handler

import (
	"sync"

	"github.com/pion/webrtc/v2"
)

type Room struct {
	Chanels *webrtc.Track
}

type RoomMap struct {
	Map   map[string]Room
	Mutex sync.RWMutex
}

func (r *RoomMap) Init() {
	r.Map = make(map[string]Room)
}

func (r *RoomMap) Get(roomID string) (Room, bool) {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()
	room, ok := r.Map[roomID]
	return room, ok
}

func (r *RoomMap) CreateRoom(roomID string) Room {
	room := Room{
		Chanels: new(webrtc.Track),
	}
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.Map[roomID] = room
	return room
}

func (r *RoomMap) DeleteRoom(roomID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	delete(r.Map, roomID)
}
