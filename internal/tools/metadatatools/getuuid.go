package metadatatools

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// GetUUIDFromMetadata- getting uuid from incoming context.
func GetUUIDFromMetadata(ctx context.Context) string {
	var uuid string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("uuid")
		if len(values) > 0 {
			uuid = values[0]
		} else {
			return ""
		}
	}
	return uuid
}
