package models

// ChangePassword => Format of body for changing password
type ChangePassword struct {
	Email       string
	OldPassword string
	NewPassword string
}
