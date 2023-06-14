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
	resp.Code = 0
	resp.Msg = newMessage.Text
	return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	resp := rpc.NewPullResponse()
	resp.Code = 1
	s.db.First(&resp.Msg)
	return resp, nil
}
