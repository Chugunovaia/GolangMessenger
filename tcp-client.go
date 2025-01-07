package main

import (
	"bufio"
	"crypto/ed25519"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"crypto/curve25519"
)

func createKeys() (my_publicKey [32]byte, my_privateKey [32]byte) { //создаем личные ключи
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		fmt.Println("Error generating keys:", err)
		return
	}
	copy(my_publicKey[:], pub[:])
	copy(my_privateKey[:], priv[:])
	curve25519.ScalarBaseMult(&my_publicKey, &my_privateKey) //домножаем, чтобы потом можно было получить общий ключ
	fmt.Println("KEYS", err)
	return
}

func createSharedKey(their_publicKey [32]byte, my_privateKey [32]byte) [32]byte {
	var pubKey, privKey, secret [32]byte
	copy(pubKey[:], their_publicKey[:])
	copy(privKey[:], my_privateKey[:])
	curve25519.ScalarMult(&secret, &privKey, &pubKey) //общий ключ, он будет одинаковым у двоих пользователей благодаря ScalarBaseMult
	return secret
}

// прием данных из сокета и вывод на печать
func readSock(conn net.Conn, my_privateKey [32]byte, codes map[[256]byte][32]byte) {

	if conn == nil {
		panic("Connection is nil")
	}
	buf := make([]byte, 256)

	//message := make([]byte, 256)
	//their_key := make([]byte, 32)
	eof_count := 0
	for {
		// чистим буфер
		for i := 0; i < 256; i++ {
			buf[i] = 0

		}

		readed_len, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				eof_count++
				time.Sleep(time.Second * 2)
				fmt.Println("EOF")
				if eof_count > 7 {

					fmt.Println("Timeout connection")
					break
				}
				continue
			}
			if strings.Index(err.Error(), "use of closed network connection") > 0 {

				fmt.Println("connection not exist or closed")
				continue
			}
			panic(err.Error())
		}
		if readed_len > 0 {
			//		fmt.Println("MESSAGE")

			/*	if (buf[0] == 227)||(buf[0]==64) {
				var friend [256]byte
				for i := 0; i < 256; i++ {
					friend[i] = 0

				}
				var i = 0
				for buf[i] != 64 {
					i++
				}
				i++
				var j = 0
				for buf[i] != 13 && buf[i] != 32 && buf[i] != 33 && buf[i] != 63 && buf[i] != 46 && buf[i] != 44 {
					friend[j] = buf[i]
					j++
					i++
				}
				_, ok := codes[friend]

				if (buf[0]==227){
					fmt.Println(227)
				if ok {
					delete(codes, friend)
				}
				} else {
					fmt.Println(64)

					if !ok {
						i=0
						//for buf[i]
						//codes[friend]
					} else {
						//дешифровка
					}
				}


			} */
			fmt.Println(string(buf))

		}

	}
}

func readConsole(ch chan string, codes map[[256]byte][32]byte, my_publicKey [32]byte) {
	for {

		line, _ := bufio.NewReader(os.Stdin).ReadString('\n')

		if len(line) > 150 {
			fmt.Println("Error: message is very lagre")
			continue
		}
		result := strings.Index(line, "@")
		//	fmt.Println(result)

		if result != -1 { //кому-то
			var friend [256]byte
			buf := make([]byte, 256)
			for i := 0; i < 256; i++ {
				buf[i] = 0
				friend[i] = 0
			}
			var i = result + 1
			var j = 0
			for line[i] != 13 && line[i] != 32 && line[i] != 33 && line[i] != 63 && line[i] != 46 && line[i] != 44 {
				friend[j] = line[i]
				j++
				i++
			}
			//	fmt.Println(friend)

			_, ok := codes[friend]
			if ok {
				//шифрование
			} else { //сообщение в виде @user(пробел=32)<сообщение(13)><публичный ключ>

				buf[0] = 64
				i := 0
				for friend[i] != 0 {
					buf[i+1] = friend[i]
					i++
				}
				i++
				buf[i] = 32
				i++
				//	fmt.Println(buf)
				for j := 0; line[j] != 10; j++ {
					//	fmt.Println(line[j])
					buf[i] = line[j]
					i++
				}
				//	fmt.Println(buf)
				for j := 0; j < 32; j++ {
					buf[i] = my_publicKey[j]
					i++
				}
				//fmt.Println(buf)
				line = string(buf)
				//fmt.Println(my_publicKey)
				//fmt.Println(buf)
				//fmt.Println(line)
			}

		}
		//a := "ㅁㅁ" //227

		out := line //[:len(line)-1] // убираем символ возврата каретки

		ch <- out // отправляем данные в канал
	}
}

func main() {
	ch := make(chan string)

	defer close(ch) // закрываем канал при выходе из программы

	conn, _ := net.Dial("tcp", "127.0.0.1:8081")
	if conn == nil {
		panic("Connection is nil")

	}
	fmt.Print("Firstly, enter your username in less than 8 ch, please: \nDo not use the ' ', '!', '?', '.', ',', '@'\n>")
	my_conns := make(map[[256]byte][32]byte, 1024)
	var my_privateKey, my_publicKey [32]byte
	my_publicKey, my_privateKey = createKeys() //my_publicKey, my_privateKey)

	go readConsole(ch, my_conns, my_publicKey)
	go readSock(conn, my_privateKey, my_conns)

	for {
		val, ok := <-ch
		if ok { // если есть данные, то их пишем в сокет
			// val_len := len(val)
			out := []byte(val)
			_, err := conn.Write(out)
			if err != nil {
				fmt.Println("Write error:", err.Error())
				break
			}

		} else {
			// данных в канале нет, задержка на 2 секунды
			time.Sleep(time.Second * 2)
		}

	}
	fmt.Println("Finished...")

	conn.Close()
}
