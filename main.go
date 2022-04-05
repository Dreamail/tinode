package main

import (
	"github.com/Dreamil/tinode_go/client"
	"log"
)

func main() {
	c := client.NewClient("10.168.1.4:16060")

	err := c.Connect()
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer func(c *client.Client) {
		err := c.Disconnect()
		if err != nil {
			log.Fatalln(err)
		}
	}(c)

	ok, err := c.Login("fan2tao", "fan2chen")
	if err != nil {
		log.Fatalln(err)
		return
	}
	if ok {
		log.Println("Login success")
	}
}
