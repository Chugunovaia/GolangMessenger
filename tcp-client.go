package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/curve25519"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
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
	return
}
func createSharedKey(their_publicKey [32]byte, my_privateKey [32]byte) [32]byte {
	var pubKey, privKey, secret [32]byte
	copy(pubKey[:], their_publicKey[:])
	copy(privKey[:], my_privateKey[:])
	curve25519.ScalarMult(&secret, &privKey, &pubKey) //общий ключ, он будет одинаковым у двоих пользователей благодаря ScalarBaseMult
	return secret
}
func findName(line string, ind int) (name [32]byte) {
	i := ind + 1
	j := 0

	for i < (len(line)) && (j < 32) && (line[i] != 13) && (line[i] != 10) && (line[i] != 32) && (line[i] != 33) && (line[i] != 63) && (line[i] != 46) && (line[i] != 44) {
		name[j] = line[i]

		j++
		i++
	}

	return
}
func createMes(user string, line string) (out string) {
	buf := make([]byte, 256)
	buf[0] = 64
	i := 0

	for i < len(user) && user[i] != 0 {
		buf[i+1] = user[i]
		i++
	}
	i++
	buf[i] = 32
	i++
	buf[i] = 33
	i++

	for j := 0; j < len(line) && line[j] != 13 && line[j] != 0; j++ {

		buf[i] = line[j]
		i++
	}
	if i < 255 {
		buf[i] = 13
		buf[i+1] = 10
	}

	out = string(buf)

	return
}
func redMes(line string) (name [32]byte, mes []byte) {
	ind := 0
	name = findName(line, ind)

	ind = strings.Index(line, "!")

	mes = make([]byte, 256)
	j := 0

	for i := ind + 1; i < len(line) && line[i] != 13; i++ {
		mes[j] = line[i]
		j++
	}

	return
}
func findKey(mes []byte) (key [32]byte) {
	i := 0

	for i < len(mes) && mes[i] != 10 {
		i++
	}
	i++
	if i < len(mes) {
		j := 0
		for i < len(mes) && j < 32 {
			key[j] = mes[i]
			i++
			j++

		}

	}
	return
}

// encryptString takes a plaintext string and a key, and returns the encrypted data in base64.
func encryptString(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to include it at the beginning of the ciphertext.
	// IV length should be equal to block size.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	// Convert the bytes to a base64 encoded string
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decryptString takes a base64 encoded string and a key, and returns the decrypted plaintext.
func decryptString(encrypted string, key []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt the ciphertext
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return string(plaintext), nil
}

// прием данных из сокета и вывод на печать
func readSock(ch chan string, conn net.Conn, codes map[[32]byte][32]byte, my_priv [32]byte, my_pub [32]byte) {

	if conn == nil {
		panic("Connection is nil")
	}
	buf := make([]byte, 256)
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
			line := string(buf)
			//ind := strings.Index(line, "@")
			if line[0] == 227 {
				ind := strings.Index(line, "@")
				name := findName(line, ind)
				_, ok := codes[name]
				if ok {
					delete(codes, name)
				}

			} else if line[0] == '@' {
				var zero_buf [32]byte
				name, mes := redMes(line)
				var l byte = 0
				mes = bytes.Trim(mes, string(l))

				_, ok := codes[name]
				if ok && codes[name] != zero_buf {
					key := make([]byte, 32)
					op := codes[name]
					copy(key[:], op[:])

					dop_line, err_d := decryptString(string(mes[:]), key)
					if err_d != nil {
						panic(err_d)
					}
					line = "@" + string(name[:]) + " ->>" + dop_line

				} else {

					their_key := findKey(mes)

					shared_key := createSharedKey(their_key, my_priv)

					codes[name] = shared_key

					i := 0
					for buf[i] != 10 {
						i++
					}
					i++
					for i < len(buf) {
						buf[i] = 0
						i++
					}
					line = string(buf)
					if !ok {
						var dop string
						dop = "\n" + string(my_pub[:])
						mes := createMes(string(name[:]), dop)

						ch <- mes

					}

				}

			}
			fmt.Println(line)
		}

	}
}

// ввод данных с консоли и вывод их в канал
func readConsole(ch chan string, codes map[[32]byte][32]byte, my_pub [32]byte) {
	for {
		line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if len(line) > 150 {
			fmt.Println("Error: message is very lagre")
			continue
		}

		ind := strings.Index(line, "@")
		if ind != -1 { //кому-то
			friend := findName(line, ind)
			_, ok := codes[friend]
			if ok {
				//шифрование
				key := make([]byte, 32)
				op := codes[friend]
				copy(key[:], op[:])

				dop_line, err_d := encryptString(line, key)

				if err_d != nil {
					panic(err_d)
				}

				line = createMes(string(friend[:]), dop_line)

			} else { //сообщение в виде @user(пробел=32)<сообщение(10)><публичный ключ>
				var zero_buf [32]byte
				codes[friend] = zero_buf
				buf := []byte(line)
				ln := len(buf)
				buf[ln-1] = 0
				buf[ln-2] = 10
				bbuf := make([]byte, 183)
				copy(bbuf, buf)
				j := 0
				for j < 183 && bbuf[j] != 0 {
					j++
				}
				if j < 183 {
					for i := 0; i < 32; i++ {
						bbuf[j] = my_pub[i]
						j++
					}
				}

				a := string(bbuf[:])

				line = createMes(string(friend[:]), a)

			}
		}

		out := line //[:len(line)-1] // убираем символ возврата каретки

		ch <- out // отправляем данные в канал
	}
}

func main() {
	ch := make(chan string)
	//local_ch := make(chan bool)
	codes := make(map[[32]byte][32]byte, 1024)
	defer close(ch) // закрываем канал при выходе из программы

	conn, _ := net.Dial("tcp", "127.0.0.1:8081")
	if conn == nil {
		panic("Connection is nil")

	}
	fmt.Print("Firstly, enter your username, please: \nDo not use the ' ','@', '!', '?', '.', ','\n>")
	my_pub, my_priv := createKeys()
	go readConsole(ch, codes, my_pub)             //local_ch,
	go readSock(ch, conn, codes, my_priv, my_pub) //local_ch,

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
