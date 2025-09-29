package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/google/uuid"
)

var comm_type = map[string]int{
	"nick":  1,
	"join":  2,
	"msg":   3,
	"list":  4,
	"exit":  5,
	"close": 6,
}

type command struct {
	argument []string
	client   *client
}

type server struct {
	commands chan command
	rooms    map[string]*room
}

func servers() *server {

	return &server{
		commands: make(chan command),
		rooms:    make(map[string]*room),
	}
}

func (s *server) newclient(conn net.Conn) {
	c := &client{
		id:       uuid.NewString(),
		conn:     conn,
		name:     "anonymous",
		commands: s.commands,
		room:     nil,
	}
	c.getcommands()
}

func (s *server) fetchcommands() {
	for cmd := range s.commands {

		comm := cmd.argument

		key := comm_type[strings.ToLower(comm[0])]

		switch key {
		case 1:
			s.setNickName(cmd)
		case 2:
			s.joinroom(cmd)
		case 3:
			s.messagerooms(cmd)

		case 4:
			s.listroom(cmd)

		case 5:
			s.exitroom(cmd)

		case 6:
			s.quit(cmd)

		}
	}
}

func (s *server) setNickName(cmd command) {

	name_args := cmd.argument[1:]

	if cmd.client.name != "anonymous" {
		if cmd.client.name == strings.TrimSpace(strings.Join(name_args, " ")) {
			cmd.client.msg("The username already exists please retry with other name")
		} else {

			cmd.client.name = strings.TrimSpace(strings.Join(name_args, " "))
			val := fmt.Sprintf("The username changed to %s", cmd.client.name)
			cmd.client.msg(val)
		}

	} else {
		cmd.client.name = strings.TrimSpace(strings.Join(name_args, " "))
		val := fmt.Sprintf("The username changed to %s", cmd.client.name)
		cmd.client.msg(val)
	}

}

func (s *server) joinroom(cmd command) {
	roomname := strings.TrimSpace(strings.Join(cmd.argument[1:], " "))
	fmt.Printf("joining the room")

	if len(roomname) > 50 {
		cmd.client.msg("The name of the room should be less than or equal to 50 characters")
		return
	}

	if cmd.client.room != nil {
		cmd.client.msg("Please exit the current room before joining the other room")
		return
	}

	_, ok := s.rooms[roomname]

	if !ok {

		s.rooms[roomname] = &room{
			name: roomname,
			members: map[string]*client{
				cmd.client.id: cmd.client,
			},
		}

		cmd.client.room = s.rooms[roomname]
		
		
		vals := fmt.Sprintf("%s has joined the room", cmd.client.name)
		cmd.client.msg(vals)
		
	} else {
		s.rooms[roomname].members[cmd.client.id] = cmd.client

		cmd.client.room = s.rooms[roomname]

		vals := fmt.Sprintf("You have successfully joined the room %s", roomname)

		fmt.Printf("Joined the room")
		cmd.client.msg(vals)
		vals = fmt.Sprintf("%s has joined the room", cmd.client.name)
		s.broadcastmessage(roomname, cmd, vals)

	}
}

func (s *server) broadcastmessage(roomname string, cmd command, msg string) {
	for _, client_data := range s.rooms[roomname].members {
		if client_data.id != cmd.client.id {
			val := msg
			client_data.msg(val)
		}
	}
}

func (s *server) messagerooms(cmd command) {

	cl_room := cmd.client.room

	msg_argument := strings.TrimSpace(strings.Join(cmd.argument[1:], " "))

	vals := fmt.Sprintf("Me -> %s", msg_argument)

	cmd.client.msg(vals)

	vals = fmt.Sprintf("%s -> %s", cmd.client.name, msg_argument)
	s.broadcastmessage(cl_room.name, cmd, vals)

}

func (s *server) listroom(cmd command) {

	var room_list []string

	for room_name, _ := range s.rooms {

		room_list = append(room_list, room_name)
	}

	room_string := strings.Join(room_list, "\n")

	cmd.client.msg(room_string)
}

func (s *server) exitroom(cmd command) {

	if cmd.client.room == nil {
		cmd.client.msg("you are not in any room, join the room to exit it")
		return
	}

	delete(s.rooms[cmd.client.room.name].members, cmd.client.id)
	cmd.client.msg("you have exited the room")
	vals := fmt.Sprintf("%s has exited the room", cmd.client.name)

	s.broadcastmessage(cmd.client.room.name, cmd, vals)
	cmd.client.room = nil

}

func (s *server) quit(cmd command) {
	cmd.client.msg("you have closed the connection")
	cmd.client.conn.Close()
}
