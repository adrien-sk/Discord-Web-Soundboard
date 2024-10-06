package websocket

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
	mu   sync.Mutex
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

var buffer = make([][]byte, 0)

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		message := Message{Type: messageType, Body: string(p)}
		fmt.Println(message)

		// Load the sound file.
		err = loadSound()
		if err != nil {
			fmt.Println("Error loading sound: ", err)
			return
		}

		// Create a new Discord session using the provided bot token.
		dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
		if err != nil {
			fmt.Println(message)
			fmt.Println(dg)
			fmt.Println("Error creating Discord session: ", err)
			return
		}

		// Open the websocket and begin listening.
		err = dg.Open()
		if err != nil {
			fmt.Println("Error opening Discord session: ", err)
		}

		err = playSound(dg, "534000212184793098", "534000212184793103")
		if err != nil {
			fmt.Println("Error playing sound:", err)
		}

		// c.Pool.Broadcast <- message
		// fmt.Printf("Message Received: %+v\n", message)
	}
}

// loadSound attempts to load an encoded sound file from disk.
func loadSound() error {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	file, err := os.Open(fmt.Sprintf("%s/%s", "sounds", "test.dca"))
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string) (err error) {

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}
