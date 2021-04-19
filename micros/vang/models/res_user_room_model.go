package models

type ResUserRoomModel struct {
	Rooms   map[string]interface{} `json:"rooms" bson:"rooms"`
	RoomIds []string               `json:"roomIds" bson:"roomIds"`
}
