package handler

import (
	"stream-service/dto"
	"sync"
)

type Room struct {
	Chanels map[string]chan dto.MessageChat
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
		Chanels: make(map[string]chan dto.MessageChat),
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
	room.Chanels[chanelId] = make(chan dto.MessageChat)
}

func (r *RoomMap) DeleteRoom(roomID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	delete(r.Map, roomID)
}

func (r *RoomMap) LeaveRoom(roomID string, chanelId string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	room := r.Map[roomID]
	delete(room.Chanels, chanelId)
}
