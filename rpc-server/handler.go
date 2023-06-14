package main

import (
	"context"

	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
	"gorm.io/gorm"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct {
	db *gorm.DB
}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	newMessage := req.GetMessage()
	s.db.Create(&newMessage)
	resp := rpc.NewSendResponse()
	resp.Code, resp.Msg = 0, "success"
	resp.Msg = newMessage.Text
	return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	resp := rpc.NewPullResponse()
	resp.Code, resp.Msg = 0, "success"
	var order string
	if *req.Reverse {
		order = "send_time desc"
	} else {
		order = "send_time"
	}
	s.db.Where("chat = ?", req.Chat).Where("send_time > ?",
		req.Cursor).Order(order).Limit(int(req.Limit)).Find(&resp.Messages)
	return resp, nil
}
