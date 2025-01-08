package main

import (
	"fmt"
	"net"
	"strings"
)

func sentmessage(message string, username string) []byte {
	/*builded := &strings.Builder{}
	out_buf := make([]byte, 512)
	//builded.WriteString(username)
	builded.WriteString("->>")
	builded.WriteString(message)
	builded.WriteString("\n")
	//out_buf := make([]byte, 512)
	extra_buf := []byte(builded.String())
	out_buf = []byte(username)
	j := 0
	for i := 0; i < 256; i++ {
		if out_buf[i] == 0 || out_buf[i] == 13 {
			if extra_buf[j] != 0 {
				out_buf[i] = extra_buf[j]
				j++
			} else {
				break

			}
		}
	}*/
	builded := &strings.Builder{}
	//out_buf := make([]byte, 256)
	plus_buf := make([]byte, 256)
	//builded.WriteString(username)
	builded.WriteString(" ->>")
	builded.WriteString(message)
	builded.WriteString("\n")
	//out_buf := make([]byte, 512)

	extra_buf := []byte(builded.String())
	//fmt.Println(extra_buf)
	for j := 0; j < 256; j++ {
		plus_buf[j] = 0
		//	out_buf[j] = 0

	}
	out_buf := []byte(username[:(len(username) - 1)])
	//fmt.Println(out_buf)
	//	fmt.Println(username)
	j := 0
	plus_buf[0] = 64
	i := 1
	//fmt.Println(len(out_buf))
	//fmt.Println(cap(out_buf))
	//fmt.Println(out_buf)
	for j < len(out_buf) && out_buf[j] != 0 {
		plus_buf[i] = out_buf[j]
		//	fmt.Println(plus_buf)
		i++
		j++
	}

	//	fmt.Println(plus_buf)
	//var h [32]byte
	//h[0] = 13
	//fmt.Println("!!!")
	//fmt.Println(extra_buf)

	j = 0
	for ii := 0; ii < 256; ii++ {

		if plus_buf[ii] == 0 {
			for j < len(extra_buf) && ii < 256 {
				plus_buf[ii] = extra_buf[j]
				//	fmt.Println(plus_buf)
				j++
				ii++
				//fmt.Println(plus_buf[ii], ii)
			}
			break
		}
	}
	//	fmt.Println(plus_buf)
	//	fmt.Println(string(plus_buf))
	return plus_buf
	//return out_buf
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
			//friend_buf[i] = 0
		}

		//fmt.Println(username)

		_, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Close ")

			}
			fmt.Println(err)
			_, ok := conns[username]
			if ok {
				delete(conns, username)
			}
			break
		}

		message = string(buf)

		if un_flag == false {
			//username = message

			b := make([]byte, 32)
			//var b [256]byte
			for i := 0; message[i] != 13 && message[i] != 0 && i < 32; i++ {
				b[i] = message[i]
			}
			username = string(b)
			_, ok := conns[username]
			if ok {
				conn.Write([]byte("This username already exists. Please, enter different one:\n>"))

			} else if strings.Contains(username, "@") || strings.Contains(username, " ") || strings.Contains(username, ",") || strings.Contains(username, ".") || strings.Contains(username, "!") || strings.Contains(username, "?") {
				//fmt.Println(buf)
				conn.Write([]byte("Used the forbidden symbol:\n>"))
			} else {
				conn.Write([]byte("Username accepted: " + username + "\n"))

				//fmt.Println(b)
				fmt.Println(username)
				conns[username] = conn
				conn.Write([]byte("To send a message to a specific user, enter their nickname as @username and the message. \nTo send a message to all users, simply leave the nickname out.\n>"))

				un_flag = true
			}
		} else {
			result := strings.Index(message, "@")
			if result == -1 {

				//out_buf := make([]byte, 512)
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
						//	friend_buf[j] = 13
						break

					}
				}
				friend = string(friend_buf)

				_, ok := conns[friend]

				if ok {
					//out_buf := make([]byte, 512)
					out_buf := sentmessage(message, username)
					//out_buf := []byte(builded.String())
					//fmt.Println(string(out_buf))
					//fmt.Println(out_buf)
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
	// создаем пул соединений
	//conns := make(map[int]net.Conn, 1024) //максимум 1024 соединеия
	count_of_conns := 0
	my_conns := make(map[string]net.Conn, 1024)
	// Устанавливаем прослушивание порта
	ln, _ := net.Listen("tcp", ":8081")
	// Запускаем цикл обработки соединений
	for {
		// Принимаем входящее соединение
		if count_of_conns < 1024 {
			conn, _ := ln.Accept()
			// сохраняем соединение в пул
			//conns[count_of_conns] = conn
			// запускаем функцию process(conn)   как горутину
			go nick(conn, my_conns, false)
			//count_of_conns--
			count_of_conns = len(my_conns)
			//fmt.Println("Number of connections:", count_of_conns)
			//	go process(conns, count_of_conns)
		}

	}
}
