# GolangMessenger
Мессенджер с шифрованием. 

Базовый в branch base cделан на основе туториала + дополнительные изменения по тз https://proglib.io/p/pishem-messendzher-na-go-za-chas-7-prostyh-shagov-ot-eho-servera-k-asinhronnomu-obmenu-soobshcheniyami-2021-09-07

Основа для защищенного branch protectred.

Если писать личное сообщение: 
1) 1 сообщение не шифруется и передает публичный ключ 1 юзера 2, 
2) после получения 2 юзер передает свой публичный ключ 1 юзеру
3) у них путем умножения своего приватного ключа на публичный ключ другого получается одинаковый общий ключ.(Взяла из статьи: https://habr.com/ru/articles/437686/)
4) Потом происходит шифрование AES 256. (Взяла отсюда: https://bytegoblin.io/blog/aes-encryption-and-decryption-in-golang-php-and-both-with-full-codes.mdx)

   
Остальные сообщения не шифруются

Для запуска надо открыть приложение GolangMessenger.exe 
Если не запускается, надо скачать всю папку и положить куда-нибудь в GOPATH. (Лично у меня это C:\Program Files\Go\src, нужна папка внутри scr) Для запуска надо использовать команду go build <папка в scr>, чтобы сообрать не 1 файл, а сразу все. И теперь уже запускать полученный .exe
