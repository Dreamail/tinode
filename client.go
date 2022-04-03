package tinode

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/tinode/chat/pbx"
)

var (
	addr = "localhost:16060"
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pbx.NewNodeClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c.MessageLoop(ctx)
}

type client struct {
	Conn pbx.NodeClient
	gprcConn grpc.ClientConn
}

func (c *client) Connect() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c.gprcConn = *conn
	c.Conn = pbx.NewNodeClient(conn)
}

func (c *client) Disconnect() {
	c.gprcConn.Close()
}