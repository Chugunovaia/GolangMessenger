# GolangMessenger
Мессенджер с шифрованием. (Базовый в branch base cделан на основе туториала + дополнительные изменения по тз https://proglib.io/p/pishem-messendzher-na-go-za-chas-7-prostyh-shagov-ot-eho-servera-k-asinhronnomu-obmenu-soobshcheniyami-2021-09-07).
Если писать личное сообщение: 
1) 1 сообщение не шифруется и передает публичный ключ 1 юзера 2, 
2) после получения 2 юзер передает свой публичный ключ 1 юзеру
3) у них путем умножения своего приватного ключа на публичный ключ другого получается одинаковый общий ключ.(Взяла из статьи: https://habr.com/ru/articles/437686/)
4) Потом происходит шифрование AES 256. (Взяла отсюда: https://bytegoblin.io/blog/aes-encryption-and-decryption-in-golang-php-and-both-with-full-codes.mdx)

Нет GUI.

Для запуска:

Сперва запускаем сервер: go run tcp-server.go (windows)

После этого запускаем клиента: go run tcp-client.go (windows)
