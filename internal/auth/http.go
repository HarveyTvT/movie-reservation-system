package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func AuthAnnotator(ctx context.Context, r *http.Request) metadata.MD {
	var (
		auth string
	)

	if v := r.Header.Get("Authorization"); v != "" {
		auth = v
	} else if v := r.CookiesNamed("authorization"); len(v) > 0 {
		auth = v[0].Value
	} else {
		return nil
	}

	if strings.HasPrefix(auth, "Bearer ") {
		auth = strings.TrimPrefix(auth, "Bearer ")
	}

	jwt, err := Parse(auth)
	if err == nil {
		return metadata.Pairs("username", jwt.Payload.Username)
	}

	return nil
}

func ForwardResponseOption(ctx context.Context, w http.ResponseWriter, response proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	if vs := md.HeaderMD["set-cookie"]; len(vs) > 0 {
		for _, v := range vs {
			w.Header().Add("Set-Cookie", v)
		}
		delete(md.HeaderMD, "set-cookie")
	}
	return nil
}

func UsernameFromContext(ctx context.Context) string {
	if v := metadata.ValueFromIncomingContext(ctx, "username"); len(v) > 0 {
		return v[0]
	}

	return ""
}
