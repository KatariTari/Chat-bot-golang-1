// main.go
package main

import (
	"context"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"

	"bot/app/controllers"
	middleware "bot/app/midleware"
	"bot/app/routes"
)

var (
	db *pgxpool.Pool
)

func init() {
	// Инициализация подключения к базе данных
	// ...

	// Применение миграций (если необходимо)
	// ...
}

func main() {
	// Инициализация бота
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	// Инициализация роутера
	router := mux.NewRouter()

	// Подключение к базе данных
	dbInfo := "your_database_info"
	db, err = pgxpool.Connect(context.Background(), dbInfo)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	// Инициализация контроллеров
	userController := controllers.NewUserController(db)
	adminController := controllers.NewAdminController(db)

	// Инициализация middleware
	adminMiddleware := middleware.AdminCheckMiddleware(db)

	// Установка роутов
	routes.SetUserRoutes(router.PathPrefix("/user").Subrouter(), userController)
	routes.SetAdminRoutes(router.PathPrefix("/admin").Subrouter(), adminController, adminMiddleware)
	labController := controllers.NewLabController(db)
	routes.SetLabRoutes(router.PathPrefix("/lab").Subrouter(), labController)

	// Обработка сообщений от бота
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		// Обработка сообщений
		if update.Message == nil {
			continue
		}

		// Пример использования контроллера для обработки сообщения
		userController.HandleMessage(update.Message)
	}
}
