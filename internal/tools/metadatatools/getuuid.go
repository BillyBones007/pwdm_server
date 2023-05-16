package metadatatools

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
)

// GetUUIDFromMetadata- getting uuid from incoming context.
func GetUUIDFromMetadata(ctx context.Context) string {
	var uuid string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fmt.Printf("INFO: md: %v+\n", md)
		values := md.Get("uuid")
		if len(values) > 0 {
			uuid = values[0]
		} else {
			return ""
		}
	}
	fmt.Println("INFO: No metadata")
	return uuid
}
