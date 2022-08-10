package handler

import (
	"sync"

	guuid "github.com/google/uuid"
)

type Room struct {
	UserID  int
	RoomID  string
	Title   string
	Chanels map[string]chan MessageVideo
}

type MessageVideo map[string]interface{}

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

func (r *RoomMap) CreateRoom(userID int, Title string) Room {

	roomID := guuid.New().String()

	room := Room{
		RoomID:  roomID,
		UserID:  userID,
		Title:   Title,
		Chanels: make(map[string]chan MessageVideo),
	}

	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.Map[roomID] = room
	return room
}

func (r *RoomMap) InsertIntoRoom(roomID string, chanelId string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	room := r.Map[roomID]
	room.Chanels[chanelId] = make(chan MessageVideo)
}

func (r *RoomMap) DeleteRoom(roomID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	delete(r.Map, roomID)
}
