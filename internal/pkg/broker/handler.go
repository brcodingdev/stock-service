package broker

// Handler contract to handle messages from broker
type Handler interface {
	// HandleStockRequest handles events received from broker
	HandleStockRequest(body string) string
}
