package event

// StockRequest ...
type StockRequest struct {
	ChatRoomName string `json:"chatRoomName"`
	ChatRoomID   uint   `json:"chatRoomId"`
	ChatMessage  string `json:"chatMessage"`
}
