package sip

import (
	"testing"

	"github.com/phughe11/catiSip/backend/config"
)

func TestNewClient(t *testing.T) {
	cfg := config.SIPConfig{
		Host:     "localhost",
		Port:     5060,
		Username: "1000",
		Password: "1234",
		Domain:   "localhost",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create SIP client: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.config.Host != "localhost" {
		t.Errorf("Expected host 'localhost', got '%s'", client.config.Host)
	}
}

func TestMakeCall(t *testing.T) {
	cfg := config.SIPConfig{
		Host:     "localhost",
		Port:     5060,
		Username: "1000",
		Password: "1234",
		Domain:   "localhost",
	}

	client, _ := NewClient(cfg)
	call, err := client.MakeCall("1000", "1001")

	if err != nil {
		t.Fatalf("Failed to make call: %v", err)
	}

	if call.From != "1000" || call.To != "1001" {
		t.Errorf("Expected from=1000 to=1001, got from=%s to=%s", call.From, call.To)
	}

	if call.Status != "dialing" {
		t.Errorf("Expected status 'dialing', got '%s'", call.Status)
	}

	if call.ID == "" {
		t.Error("Expected non-empty call ID")
	}
}

func TestGetCallStatus(t *testing.T) {
	cfg := config.SIPConfig{
		Host:     "localhost",
		Port:     5060,
		Username: "1000",
		Password: "1234",
		Domain:   "localhost",
	}

	client, _ := NewClient(cfg)
	call, _ := client.MakeCall("1000", "1001")

	status, err := client.GetCallStatus(call.ID)
	if err != nil {
		t.Fatalf("Failed to get call status: %v", err)
	}

	if status.ID != call.ID {
		t.Errorf("Expected call ID '%s', got '%s'", call.ID, status.ID)
	}
}

func TestGetCallStatusNotFound(t *testing.T) {
	cfg := config.SIPConfig{
		Host:     "localhost",
		Port:     5060,
		Username: "1000",
		Password: "1234",
		Domain:   "localhost",
	}

	client, _ := NewClient(cfg)
	_, err := client.GetCallStatus("non-existent-id")

	if err == nil {
		t.Error("Expected error for non-existent call ID")
	}
}

func TestHangupCall(t *testing.T) {
	cfg := config.SIPConfig{
		Host:     "localhost",
		Port:     5060,
		Username: "1000",
		Password: "1234",
		Domain:   "localhost",
	}

	client, _ := NewClient(cfg)
	call, _ := client.MakeCall("1000", "1001")

	err := client.HangupCall(call.ID)
	if err != nil {
		t.Fatalf("Failed to hangup call: %v", err)
	}

	status, _ := client.GetCallStatus(call.ID)
	if status.Status != "ended" {
		t.Errorf("Expected status 'ended', got '%s'", status.Status)
	}
}

func TestListCalls(t *testing.T) {
	cfg := config.SIPConfig{
		Host:     "localhost",
		Port:     5060,
		Username: "1000",
		Password: "1234",
		Domain:   "localhost",
	}

	client, _ := NewClient(cfg)
	
	call1, _ := client.MakeCall("1000", "1001")
	call2, _ := client.MakeCall("1000", "1002")

	calls := client.ListCalls()
	
	// Check that both calls are in the list
	found1, found2 := false, false
	for _, call := range calls {
		if call.ID == call1.ID {
			found1 = true
		}
		if call.ID == call2.ID {
			found2 = true
		}
	}
	
	if !found1 || !found2 {
		t.Errorf("Expected both calls to be in the list. Found call1: %v, Found call2: %v", found1, found2)
	}
}
