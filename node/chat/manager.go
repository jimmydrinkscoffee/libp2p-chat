package chat

type Manager interface {
	Join(string, string) error
	HasJoined(string) bool
}

type manager struct{}

var _ Manager = (*manager)(nil)

func NewManager() Manager {
	return &manager{}
}

func (m *manager) Join(roomName string, nickName string) error {
	if m.HasJoined(roomName) {
		return nil
	}

	return nil
}

func (m *manager) HasJoined(roomName string) bool {
	return false
}
