package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type client struct {
	id       string
	conn     net.Conn
	name     string
	room     *room
	commands chan<- command
}

func (c *client) getcommands() {
	count := 0
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')

		if err != nil {
			if count == 3 {
				return
			}
			count += 1
			log.Fatalf("invalid connection recieved")
		} else {
			count = 0
			msg = strings.Trim(msg, "\r\n")
			msg_list := strings.Split(msg, " ")

			var comm_type = map[string]int{
				"nick":  1,
				"join":  2,
				"msg":   3,
				"list":  4,
				"exit":  5,
				"close": 6,
			}

			_, ok := comm_type[strings.ToLower(msg_list[0])]

			if !ok {
				c.err(fmt.Errorf("invalid input command %s", msg_list[0]))

			}
			c.commands <- command{
				argument: msg_list,
				client:   c,
			}

		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("Err : " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte(msg + "\n"))
}
