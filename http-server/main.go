package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/TikTokTechImmersion/assignment_demo_2023/http-server/kitex_gen/rpc"
	"github.com/TikTokTechImmersion/assignment_demo_2023/http-server/kitex_gen/rpc/imservice"
	"github.com/TikTokTechImmersion/assignment_demo_2023/http-server/proto_gen/api"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var cli imservice.Client

func main() {
	r, err := etcd.NewEtcdResolver([]string{"etcd:2379"})
	if err != nil {
		log.Fatal(err)
	}
	cli = imservice.MustNewClient("demo.rpc.server",
		client.WithResolver(r),
		client.WithRPCTimeout(1*time.Second),
		client.WithHostPorts("rpc-server:8888"),
	)

	h := server.Default(server.WithHostPorts("0.0.0.0:8080"))

	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"message": "pong"})
	})

	h.POST("/api/send", sendMessage)
	h.GET("/api/pull", pullMessage)

	h.Spin()
}

func sendMessage(ctx context.Context, c *app.RequestContext) {
	resp, err := cli.Send(ctx, &rpc.SendRequest{
		Message: &rpc.Message{
			Chat:     c.Query("chat"),
			Text:     c.Query("text"),
			Sender:   c.Query("sender"),
			SendTime: time.Now().Unix(),
		},
	})
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
	} else if resp.Code != 0 {
		c.String(consts.StatusInternalServerError, resp.Msg)
	} else {
		c.Status(consts.StatusOK)
	}
}

func pullMessage(ctx context.Context, c *app.RequestContext) {
	cursor, err := strconv.Atoi(c.Query("cursor"))
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.String(consts.StatusBadRequest, "Failed to parse request body: %v", err)
		return
	}
	reverse := c.Query("reverse") == "true"
	resp, err := cli.Pull(ctx, &rpc.PullRequest{
		Chat:    c.Query("chat"),
		Cursor:  int64(cursor),
		Limit:   int32(limit),
		Reverse: &reverse,
	})
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	} else if resp.Code != 0 {
		c.String(consts.StatusInternalServerError, resp.Msg)
		return
	}
	messages := make([]*api.Message, 0, len(resp.Messages))
	for _, msg := range resp.Messages {
		messages = append(messages, &api.Message{
			Chat:     msg.Chat,
			Text:     msg.Text,
			Sender:   msg.Sender,
			SendTime: msg.SendTime,
		})
	}
	c.JSON(consts.StatusOK, &api.PullResponse{
		Messages:   messages,
		HasMore:    resp.GetHasMore(),
		NextCursor: resp.GetNextCursor(),
	})
}
