package provider

import (
	"context"
	"github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ReadSchema(client *authzed.Client, ctx context.Context) (string, error) {
	authResp, err := client.ReadSchema(ctx, &v1.ReadSchemaRequest{})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return "", nil
		}
		return "", err
	}
	return authResp.SchemaText, nil
}
