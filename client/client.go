package client

import (
	"context"
	"errors"
	"github.com/tinode/chat/pbx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/rand"
	"strconv"
)

type Client struct {
	Addr string

	Conn        pbx.NodeClient
	MessageLoop pbx.Node_MessageLoopClient

	gprcConn *grpc.ClientConn

	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewClient(addr string) *Client {
	return &Client{
		Addr: addr,
	}
}

func (c *Client) Connect() error {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	c.gprcConn = conn
	c.Conn = pbx.NewNodeClient(conn)

	c.ctx, c.ctxCancel = context.WithCancel(context.Background())

	c.MessageLoop, err = c.Conn.MessageLoop(c.ctx)
	if err != nil {
		return err
	}

	err = c.MessageLoop.Send(&pbx.ClientMsg{
		Message: &pbx.ClientMsg_Hi{
			Hi: &pbx.ClientHi{
				UserAgent:  "tinode-bridge",
				Ver:        "0.18.3",
				Background: true,
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Disconnect() error {
	c.ctxCancel()
	return c.gprcConn.Close()
}

func (c *Client) Login(user, passwd string) (bool, error) {
	msgId := strconv.FormatInt(int64(rand.Int()), 10)

	err := c.MessageLoop.Send(&pbx.ClientMsg{
		Message: &pbx.ClientMsg_Login{
			Login: &pbx.ClientLogin{
				Id:     msgId,
				Scheme: "basic",
				Secret: []byte(user + ":" + passwd),
			},
		},
	})
	if err != nil {
		return false, err
	}

	for {
		serverMsg, err := c.MessageLoop.Recv()
		if err != nil {
			return false, err
		}
		if ctrlMsg := serverMsg.GetCtrl(); ctrlMsg.Id == msgId {
			if ctrlMsg.Text == "ok" {
				return true, nil
			} else {
				log.Println(ctrlMsg.Text)
				return false, errors.New("login failed")
			}
		}
	}
}
