# Telegram Bot. Getting YouTube image previews

Version [ENG](https://github.com/tecspda/get-img-from-youtube-bot/README.md) | [RUS](https://github.com/tecspda/get-img-from-youtube-bot/README_RU.md)

## Install
1. Install golang: [go.dev/doc/install](https://go.dev/doc/install)
2. Clone repo:
```sh
md youtube-bot && cd youtube-bot
git clone https://github.com/tecspda/get-img-from-youtube-bot.git .
```
3. Initialize the Golang project
```sh
go mod init modules
go mod tidy
```

## Compilation for Unix

```sh
GOOS=linux GOARCH=amd64 go build -o youtubeimgbot
```
In this example, the amd64 parameter is specified for compiling the executable file for Linux x86. If you want to compile for your system, enter the command

```sh
go build youtubeimgbot
```

## Copy to production-server
If you compiled the bot on your local computer:
1. Copy only three files to the server: `youtubeimgbot.sh`, `youtubeimgbot`, `.env`.
2. Don't forget to set the permissions:
```
cd /var/www/bots/get_youtube_img_bot
chmod 0775 youtubeimgbot.*
```
3. Create a directory uploads on the server. Grant 775 permissions.

## Bot Configuration
1. [Регистрация бота](https://www.google.com/search?q=botfather+create+bot)
2. Rename the example.env file to .env and enter the bot key (see the above point)

## Running the Bot as a Service
### Creating the Service File
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

### Starting the Daemon on the Server
```sh
systemctl start youtubeimgbot
systemctl status youtubeimgbot
systemctl enabled youtubeimgbot
```

## Done
Congratulations! Your standalone bot is now running on the server.
