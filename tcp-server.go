package main

import (
	"fmt"
	"net"
	"strings"
)

func sentmessage(message string, username string) []byte {
	builded := &strings.Builder{}
	plus_buf := make([]byte, 256)
	builded.WriteString(" ->>")
	builded.WriteString(message)
	builded.WriteString("\n")

	extra_buf := []byte(builded.String())
	for j := 0; j < 256; j++ {
		plus_buf[j] = 0

	}
	out_buf := []byte(username[:(len(username) - 1)])
	j := 0
	plus_buf[0] = 64
	i := 1
	for j < len(out_buf) && out_buf[j] != 0 {
		plus_buf[i] = out_buf[j]
		i++
		j++
	}

	j = 0
	for ii := 0; ii < 256; ii++ {

		if plus_buf[ii] == 0 {
			for j < len(extra_buf) && ii < 256 {
				plus_buf[ii] = extra_buf[j]

				j++
				ii++

			}
			break
		}
	}

	return plus_buf

}

func nick(conn net.Conn, conns map[string]net.Conn, un_flag bool) {
	buf := make([]byte, 256)

	fmt.Print("Accept cnn:")
	var message string = ""
	var username = ""
	defer conn.Close()
	for {

		for i := 0; i < 255; i++ {
			buf[i] = 0
		}

		_, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Close ")

			}
			fmt.Println(err)
			_, ok := conns[username]
			if ok {
				delete(conns, username)
				for _, connection := range conns {
					out_buf := []byte("ㅁUser @" + username + " finished their work")
					_, err2 := connection.Write(out_buf)
					if err2 != nil {
						fmt.Println("Error:", err2.Error())
						break
					}

				}

			}
			break
		}

		message = string(buf)

		if un_flag == false {

			b := make([]byte, 32)

			for i := 0; message[i] != 13 && message[i] != 0 && i < 32; i++ {
				b[i] = message[i]
			}
			username = string(b)
			_, ok := conns[username]
			if ok {
				conn.Write([]byte("This username already exists. Please, enter different one:\n>"))

			} else if strings.Contains(username, "@") || strings.Contains(username, " ") || strings.Contains(username, ",") || strings.Contains(username, ".") || strings.Contains(username, "!") || strings.Contains(username, "?") {

				conn.Write([]byte("Used the forbidden symbol:\n>"))
			} else {
				conn.Write([]byte("Username accepted: " + username + "\n"))

				fmt.Println(username)
				conns[username] = conn
				conn.Write([]byte("To send a message to a specific user, enter their nickname as @username and the message. \nTo send a message to all users, simply leave the nickname out.\n>"))

				un_flag = true
			}
		} else {
			result := strings.Index(message, "@")
			if result == -1 {

				for _, connection := range conns {
					out_buf := []byte(message)
					_, err2 := connection.Write(out_buf)
					if err2 != nil {
						fmt.Println("Error:", err2.Error())
						break
					}

				}
			} else {
				friend_buf := make([]byte, 32)
				var friend = ""
				for i := result; i < 256; i++ {
					if buf[i] == 64 {
						i++
						j := 0
						for buf[i] != 13 && buf[i] != 32 && buf[i] != 33 && buf[i] != 63 && buf[i] != 46 && buf[i] != 44 {
							friend_buf[j] = buf[i]
							j++
							i++
						}

						break

					}
				}
				friend = string(friend_buf)

				_, ok := conns[friend]

				if ok {

					out_buf := sentmessage(message, username)

					_, err2 := conns[friend].Write(out_buf)
					if err2 != nil {
						fmt.Println("Error:", err2.Error())
						break
					}
				} else {
					conn.Write([]byte("ㅁUser @" + friend + " does not exist or finished it's work. \n>"))
				}
			}
		}

	}

}

func main() {

	fmt.Println("Start server...")

	count_of_conns := 0
	my_conns := make(map[string]net.Conn, 1024)

	ln, _ := net.Listen("tcp", ":8081")

	for {

		if count_of_conns < 1024 {
			conn, _ := ln.Accept()

			go nick(conn, my_conns, false)

			count_of_conns = len(my_conns)

		}

	}
}
