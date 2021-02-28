package constants

type UserPermissionConst string

const (
	OnlyMe  UserPermissionConst = "OnlyMe"
	Public  UserPermissionConst = "Public"
	Circles UserPermissionConst = "Circles"
	Custom  UserPermissionConst = "Custom"
)
