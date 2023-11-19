# Telegram bot для получения превью youtube-ролика

Version [ENG](https://github.com/tecspda/get-img-from-youtube-bot/README.md) | [RUS](https://github.com/tecspda/get-img-from-youtube-bot/README_RU.md)

## Установка
1. Установка golang: [go.dev/doc/install](https://go.dev/doc/install)
2. Клонирование репозитория:
```sh
md youtube-bot && cd youtube-bot
git clone https://github.com/tecspda/get-img-from-youtube-bot.git .
```
3. Инициализация проекта Golang
```sh
go mod init modules
go mod tidy
```

## Компилирование для Unix

```sh
GOOS=linux GOARCH=amd64 go build -o youtubeimgbot
```
В данном примере указан параметр amd64 для компиляции исполняемого файла для Linux x86. Если вы хотите скомпилировать под свою систему, то введите команду

```sh
go build youtubeimgbot
```

## Копирование на production-сервер
Если вы компилировали бот на локальном компьютере:
1. Cкопируйте на сервер три только файла `youtubeimgbot.sh`, `youtubeimgbot`, `.env`.
2. Не забудьте установить им права
```
cd /var/www/bots/get_youtube_img_bot
chmod 0775 youtubeimgbot.*
```
3. Создайте на сервере дирректорию `uploads`. Выдайте права 775.

## Настройка бота
1. [Регистрация бота](https://www.google.com/search?q=botfather+%D1%81%D0%BE%D0%B7%D0%B4%D0%B0%D0%BD%D0%B8%D0%B5+%D0%B1%D0%BE%D1%82%D0%B0)
2. Переименуйте файл `example.env` в `.env` и введите ключ от бота (см. пункт выше)

## Запускаем бота как сервис
### Создание файла сервиса
```sh
sudo nano /etc/systemd/system/youtubeimgbot.service
```
```
[Unit]
Description=Youtube bot Service
After=network.target

[Service]
ExecStart=/var/www/bots/get_youtube_img_bot/youtubeimgbot.sh
Restart=always

[Install]
WantedBy=default.target
```

### Запуск демона на сервере
```sh
systemctl start youtubeimgbot
systemctl status youtubeimgbot
systemctl enabled youtubeimgbot
```

## Готово
Поздравляю. У вас на сервере работает саомстоятельный бот
