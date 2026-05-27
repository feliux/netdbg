package revdns

import (
	"context"
	"testing"
)

// TestGetResolver_Default verifies the default resolver behavior.
func TestGetResolver_Default(t *testing.T) {
	opts := &Options{}
	res := getResolver(opts)
	if res == nil {
		t.Fatal("expected non-nil resolver")
	}
	// Default resolver should not force PreferGo
	if res.PreferGo {
		t.Errorf("expected PreferGo=false for default resolver, got true")
	}
	if res.Dial != nil {
		t.Errorf("expected Dial=nil for default resolver, got non-nil")
	}
}

// TestGetResolver_Custom verifies the custom resolver behavior.
func TestGetResolver_Custom(t *testing.T) {
	opts := &Options{
		ResolverIP: "8.8.8.8",
		Protocol:   "udp",
		Port:       53,
	}
	res := getResolver(opts)
	if res == nil {
		t.Fatal("expected non-nil resolver")
	}
	if !res.PreferGo {
		t.Errorf("expected PreferGo=true for custom resolver, got false")
	}
	if res.Dial == nil {
		t.Fatalf("expected non-nil Dial for custom resolver")
	}

	// Test that the Dial function returns a net.Conn or error (without actually connecting)
	ctx := context.Background()
	network := "udp"
	address := "8.8.8.8:53"
	// The following should not panic, but will likely fail to connect (which is fine for this test)
	_, err := res.Dial(ctx, network, address)
	if err == nil {
		// It's possible to succeed if the network is available, but we mainly care that Dial is callable
		t.Log("dial succeeded (this is fine if network is available)")
	}
}
