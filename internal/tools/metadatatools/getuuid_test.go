package metadatatools

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestGetUUIDFromMetadata(t *testing.T) {
	// Test case 1: When metadata contains UUID
	md := metadata.New(map[string]string{"uuid": "123"})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	result := GetUUIDFromMetadata(ctx)
	if result != "123" {
		t.Errorf("Expected: %s, but got: %s", "123", result)
	}
	// Test case 2: When metadata does not contain UUID
	md = metadata.New(map[string]string{"foo": "bar"})
	ctx = metadata.NewIncomingContext(context.Background(), md)
	result = GetUUIDFromMetadata(ctx)
	if result != "" {
		t.Errorf("Expected: %s, but got: %s", "", result)
	}

	// Test case 3: When context is empty
	result = GetUUIDFromMetadata(context.Background())
	if result != "" {
		t.Errorf("Expected: %s, but got: %s", "", result)
	}
}
