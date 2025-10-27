package sip

import (
	"fmt"
	"sync"
	"time"

	"github.com/phughe11/catiSip/backend/config"
)

// Client represents a SIP client connection to FreeSWITCH
type Client struct {
	config *config.SIPConfig
	calls  map[string]*Call
	mu     sync.RWMutex
}

// Call represents an active SIP call
type Call struct {
	ID          string    `json:"id"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Status      string    `json:"status"`
	StartTime   time.Time `json:"start_time"`
	AnswerTime  time.Time `json:"answer_time,omitempty"`
	EndTime     time.Time `json:"end_time,omitempty"`
}

// NewClient creates a new SIP client
func NewClient(cfg config.SIPConfig) (*Client, error) {
	client := &Client{
		config: &cfg,
		calls:  make(map[string]*Call),
	}

	// In a real implementation, this would establish connection to FreeSWITCH
	// using ESL (Event Socket Library) or similar
	return client, nil
}

// MakeCall initiates a new call
func (c *Client) MakeCall(from, to string) (*Call, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	callID := fmt.Sprintf("call-%d", time.Now().Unix())
	call := &Call{
		ID:        callID,
		From:      from,
		To:        to,
		Status:    "dialing",
		StartTime: time.Now(),
	}

	c.calls[callID] = call

	// In real implementation, this would use FreeSWITCH ESL to originate call
	// For now, simulate call progression
	go c.simulateCall(callID)

	return call, nil
}

// HangupCall terminates an active call
func (c *Client) HangupCall(callID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	call, exists := c.calls[callID]
	if !exists {
		return fmt.Errorf("call not found: %s", callID)
	}

	call.Status = "ended"
	call.EndTime = time.Now()

	// In real implementation, send hangup command to FreeSWITCH
	return nil
}

// GetCallStatus returns the status of a call
func (c *Client) GetCallStatus(callID string) (*Call, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	call, exists := c.calls[callID]
	if !exists {
		return nil, fmt.Errorf("call not found: %s", callID)
	}

	return call, nil
}

// ListCalls returns all active calls
func (c *Client) ListCalls() []*Call {
	c.mu.RLock()
	defer c.mu.RUnlock()

	calls := make([]*Call, 0, len(c.calls))
	for _, call := range c.calls {
		calls = append(calls, call)
	}
	return calls
}

// Close closes the SIP client connection
func (c *Client) Close() error {
	// Clean up resources
	return nil
}

// simulateCall simulates call progression (for demonstration)
func (c *Client) simulateCall(callID string) {
	time.Sleep(2 * time.Second)

	c.mu.Lock()
	if call, exists := c.calls[callID]; exists {
		call.Status = "ringing"
	}
	c.mu.Unlock()

	time.Sleep(3 * time.Second)

	c.mu.Lock()
	if call, exists := c.calls[callID]; exists {
		call.Status = "answered"
		call.AnswerTime = time.Now()
	}
	c.mu.Unlock()
}
