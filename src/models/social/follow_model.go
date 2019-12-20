package models

type FollowModel struct {
	RightUser RelMetaModel `json:"right"`
	CircleIds []string     `json:"circleIds"`
}
