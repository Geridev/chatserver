package server

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Valami struct {
	Client []Client
}

type Client struct {
	Id   string
	Name User
	Conn *websocket.Conn
}

type User struct {
	Name string `json:"name"`
}

type Message struct {
	Msg      string `json:"msg"`
	UserName User
}

func (v *Valami) AddClient(client *Client) *Valami {
	v.Client = append(v.Client, *client)
	fmt.Println(client.Id)
	return v
}

func (v *Valami) RemoveClient(client *Client) *Valami {
	for i, c := range v.Client {
		if c.Id == client.Id {
			v.Client = append(v.Client[:i], v.Client[i+1:]...)
		}
	}
	return v
}

func (v *Valami) SendMessage(p []byte, client *Client) *Valami {
	m := &Message{
		Msg:      string(p),
		UserName: client.Name,
	}
	for _, c := range v.Client {
		c.Conn.WriteJSON(m)
	}
	return v
}

func (v *Valami) GetUserName(client *Client) *Valami {
	_, p, err := client.Conn.ReadMessage()
	if err != nil {
		fmt.Println(err)
	}
	client.Name.Name = string(p)
	return v
}
