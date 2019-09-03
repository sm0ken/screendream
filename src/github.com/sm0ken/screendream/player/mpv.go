package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"time"
)

// plays a file at a given location
// not sure about how we want to pass file specifiers here
func Play(arg1 string) {
	app := "mpv"

	arg0 := "--input-ipc-server=/tmp/mpvsocket"

	cmd := exec.Command(app, arg0, arg1)
	log.Printf("Playing %s in mpv, IPC enabled via /tmp/mpvsocket", arg1)

	out, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
		fmt.Print(out)
		return
	}
	log.Println("Successfully finished playing %s", arg1)
}

func connect() net.Conn {
	c, err := net.Dial("unix", "/tmp/mpvsocket")
	if err != nil {
		log.Fatal(err)
		panic(err) // TODO: better error handling
	}

	return c
}

func setProperty(p, v string) string {
	return fmt.Sprintf(`{"command": ["set_property", "%s", %s]}`, p, v) + "\n"
}

func getProperty(p string) string {
	return fmt.Sprintf(`{"command": ["get_property", "%s"]}`, p) + "\n"
}

// exits mpv
func Quit() {
	c := connect()
	defer c.Close()

	_, err := c.Write([]byte(`{"command": ["quit"]}` + "\n")) // TODO check amount of bytes sent
	if err != nil {
		log.Fatal(err)
	}
}

// Currently broken, TODO either use static variable or get current state
func Toggle() {
	c := connect()
	defer c.Close()

	cmd := setProperty("pause", "true")

	_, _ = c.Write([]byte(cmd))
}

func Pause() {
	c := connect()
	defer c.Close()

	cmd := setProperty("pause", "true")

	_, _ = c.Write([]byte(cmd))
}

func Resume() {
	c := connect()
	defer c.Close()

	cmd := setProperty("pause", "false")

	_, _ = c.Write([]byte(cmd))
}

func Jump(pos int) {
	c := connect()
	defer c.Close()

	cmd := fmt.Sprintf(`{"command": ["seek", "%d", "absolute"]}`, pos) + "\n"

	_, _ = c.Write([]byte(cmd))
}

// TODO: fix, reading the message does not work yet
func Position() int {
	c := connect()
	defer c.Close()

	cmd := getProperty("playback-time")

	_, _ = c.Write([]byte(cmd))

	buff := make([]byte, 0)

	i, _ := c.Read(buff)

	fmt.Println(buff)

	return i
}

func main() {
	go Play("/Users/zero/feet.mp4")

	time.Sleep(5000 * time.Millisecond)
	Position()

	Jump(360)

	time.Sleep(2000 * time.Millisecond)
	Position()
}

// TODO: command for subtitles/audio tracks, anything else?
