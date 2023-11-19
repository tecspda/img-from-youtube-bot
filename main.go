package main

import (
	"fmt"
	"io"
	"log"
	"modules/app/models"
	"net/http"
	"os"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var uploadPath string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	tgToken := os.Getenv("TG_TOKEN")
	uploadPath = os.Getenv("UPLOAD_PATH")

	db, err := models.NewDatabase()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	db.CreateTable()

	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Println(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	var chatID int64
	for update := range updates {
		if update.Message == nil { // игнорируем обновления, не являющиеся сообщениями
			continue
		}

		chatID = update.Message.Chat.ID
		keyboard := models.GetKb(chatID, db)

		text := update.Message.Text
		text_lower := strings.ToLower(text)
		msg := tgbotapi.NewMessage(chatID, "")
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = keyboard
		msg.DisableWebPagePreview = true

		if update.Message.IsCommand() || update.CallbackQuery != nil {
			switch update.Message.Command() {
			case "start", "help":
				msg.Text = "<b>СПРАВКА</b>\n\nВставьте url любого Youtube-ролика и получите картинку его превью."
				if update.Message.Command() == "start" {
					db.SaveChatId(update, bot)
				}
			default:
				msg.Text = "Неизвестная команда. Для примера введите <b>Погода Владивосток</b>, чтобы узнать погоду."
				db.SaveError(update, text)
			}
			bot.Send(msg)
		} else if strings.Contains(text_lower, "start") || strings.Contains(text_lower, "помощь") {
			msg.Text = "<b>СПРАВКА</b>\nОтправьте этому боту ссылку любого Youtube-ролика и получите картинку его превью."
			bot.Send(msg)
			continue
		} else if strings.Contains(text_lower, "наши проекты") {
			msg.Text = "<b>НАШИ ПРОЕКТЫ</b>\n\n"
			msg.Text += `👉 <a href="https://ne-propusti.ru">НЕ-ПРОПУСТИ.РУ</a> - отслеживание цен на Велдберис каждый час с уведомлением об изменении.` + "\n"
			msg.Text += "👉 @ShkolaPozitiva - канал ШколаПозитива. Лучший контент о конспирологии и мире.\n"
			msg.Text += "👉 @Theweather2023bot - погодный бот. Погода в любом городе мира.\n"
			msg.Text += "👉 @Moondays2024_bot - бот с описанием лунных дней, календарь новолуний и полнолуний.\n"
			msg.Text += "👉 @GetYoutubeImgBot - бот для получения превью-картинок из Youtube.\n"
			msg.Text += "\n<b>СПРАВКА</b>\nОтправьте этому боту ссылку любого Youtube-ролика и получите картинку его превью."
			bot.Send(msg)
			continue
		} else {
			videoID, err := extractYouTubeVideoID(text)
			if err != nil {
				log.Println(err)
				db.SaveError(update, text)
				msg.Text = "<b>ОШИБКА</b>. Не удалось извлечь картинку. Попробуйте еще."
				bot.Send(msg)
				continue
			}

			videoURL := fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", videoID)
			imageFileName, err := downloadImage(videoURL, videoID)
			if err != nil {
				msg.Text = "<b>ОШИБКА</b>. Не удалось загрузить картинку."
				bot.Send(msg)
				log.Println(err)
			} else {
				photo := tgbotapi.NewPhotoUpload(chatID, imageFileName)
				photo.Caption = fmt.Sprintf("\nСсылка на изображение в оригинальном качестве:\n%s\n", videoURL)
				bot.Send(msg)

				_, err = bot.Send(photo)
				if err != nil {
					log.Println(err)
				}

				// Удалить скачанное изображение после отправки
				os.Remove(imageFileName)
			}
		}
	}
}

func downloadImage(url string, ytID string) (string, error) {
	// Отправить HTTP-запрос
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Проверить статус ответа
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Неверный статус ответа: %d", response.StatusCode)
	}

	// Создать файл для сохранения изображения
	imageFileName := generateUniqueFileName("jpg")
	imagePath := uploadPath + "/" + imageFileName
	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("Не удалось создать файл")
	}
	defer file.Close()

	// Записать тело ответа в файл
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	// fmt.Printf("Изображение успешно сохранено в %s\n", destination)
	return imagePath, nil
}

func generateUniqueFileName(extension string) string {
	uuid := uuid.New()
	return fmt.Sprintf("%s.%s", strings.Replace(uuid.String(), "-", "", -1), extension)
}

func extractYouTubeVideoID(url string) (string, error) {
	// Паттерн регулярного выражения для извлечения идентификатора видео
	pattern := `(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/(?:[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})`
	re := regexp.MustCompile(pattern)

	// Найти совпадение в URL
	matches := re.FindStringSubmatch(url)

	// Если есть совпадение, вернуть идентификатор видео
	if len(matches) >= 2 {
		return matches[1], nil
	} else {
		urlFinded, err := extractYouTubeVideoID2(url)
		if err == nil {
			return urlFinded, nil
		}
	}

	// В противном случае вернуть ошибку
	return "", fmt.Errorf("Не удалось извлечь идентификатор видео из URL")
}

func extractYouTubeVideoID2(url string) (string, error) {
	// Паттерн регулярного выражения для извлечения идентификатора видео
	pattern := `(?:https?:\/\/)?(?:www\.)?youtube\.com\/live\/([^\/\?\&]+)`
	re := regexp.MustCompile(pattern)

	// Найти совпадение в URL
	matches := re.FindStringSubmatch(url)

	// Если есть совпадение, вернуть идентификатор видео
	if len(matches) >= 2 {
		return matches[1], nil
	}

	// В противном случае вернуть ошибку
	return "", fmt.Errorf("Не удалось извлечь идентификатор видео из URL")
}
