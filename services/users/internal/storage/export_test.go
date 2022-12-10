package storage

var GetStringValue = getStringValue

func NewTestUserEditor(userDb userDb) *UserEditor {
	return &UserEditor{userDb}
}
