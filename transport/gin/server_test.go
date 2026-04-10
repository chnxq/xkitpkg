package gin

import (
	"context"
	"math/rand"
	"strconv"
	"testing"

	kHttp "github.com/chnxq/xkitpkg/transport/http"
	"github.com/chnxq/xkitpkg/transport/http/binding"
	pb "github.com/chnxq/xkitpkg/transport/internal/testdata/helloworld"
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	ctx := context.Background()

	srv := NewServer(
		WithAddress(":8800"),
	)

	srv.Use(gin.Recovery())
	srv.Use(gin.Logger())

	srv.GET("/login/*param", func(c *gin.Context) {
		if len(c.Params.ByName("param")) > 1 {
			c.AbortWithStatus(404)
			return
		}
		c.String(200, "Hello World!")
	})

	srv.GET("/hello", func(c *gin.Context) {
		var out pb.HelloReply
		out.Message = strconv.FormatInt(int64(rand.Intn(100)), 10)
		c.JSON(200, &out)
	})

	if err := srv.Start(ctx); err != nil {
		panic(err)
	}

	defer func() {
		if err := srv.Stop(ctx); err != nil {
			t.Errorf("expected nil got %v", err)
		}
	}()
}

func TestClient(t *testing.T) {
	ctx := context.Background()

	cli, err := kHttp.NewClient(ctx,
		kHttp.WithEndpoint("127.0.0.1:8800"),
	)
	assert.Nil(t, err)
	assert.NotNil(t, cli)

	resp, err := GetHygrothermograph(ctx, cli, nil, kHttp.EmptyCallOption{})
	assert.Nil(t, err)
	t.Log(resp)
}

func GetHygrothermograph(ctx context.Context, cli *kHttp.Client, in *pb.HelloRequest, opts ...kHttp.CallOption) (*pb.HelloReply, error) {
	var out pb.HelloReply

	pattern := "/hello"
	path := binding.EncodeURL(pattern, in, true)

	opts = append(opts, kHttp.Operation("/GetHello"))
	opts = append(opts, kHttp.PathTemplate(pattern))

	err := cli.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
