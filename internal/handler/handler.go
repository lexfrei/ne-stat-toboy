// Package handler provides HTTP handlers for the "Ne Stat Toboy" film website.
package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lexfrei/ne-stat-toboy/internal/model"
	"github.com/lexfrei/ne-stat-toboy/web/template"
)

// TelegramMessage represents the structure for sending messages to Telegram API
type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// Handler contains all HTTP handlers and their dependencies.
type Handler struct {
	FilmInfo       model.FilmInfo
	TelegramToken  string
	TelegramChatID string
}

// HandlerOption is a functional option for configuring the handler
type HandlerOption func(*Handler)

// WithTelegramConfig configures Telegram notification settings
func WithTelegramConfig(token, chatID string) HandlerOption {
	return func(h *Handler) {
		h.TelegramToken = token
		h.TelegramChatID = chatID
	}
}

// New creates a new Handler with initialized dependencies.
func New(opts ...HandlerOption) *Handler {
	h := &Handler{
		// Default empty values for Telegram config
		TelegramToken:  "",
		TelegramChatID: "",
		FilmInfo: model.FilmInfo{
			Title:           "Не стать тобой",
			Tagline:         "Жизнь в семье, где любовь выражается насилием, ложью и запретами, подталкивает школьницу к тяжелому выбору - смерть или предательство ради жизни",
			Synopsis:        "Ева мечтает стать музыкантом, но её жизнь жёстко контролируют мать и бабушка, запрещая любые увлечения и навязывая ей своё будущее. Отец ушёл, но говорить об этом нельзя — мать создаёт видимость «идеальной семьи». Дома же Еву наказывают за любые ошибки. Когда после визита подруги Вики у Евы пропадают наушники - подарок отца, бабушка обвиняет в этом Вику, но Ева случайно находит их в вещах бабушки и начинает сомневаться в семье. Вика убеждает Еву впервые нарушить правила и прогулять школу. Девушки веселятся, и Ева раскрывает правду о семье, не зная, что её снимают на видео. Запись попадает в сеть, начинается травля и сильные побои дома. После побоев Ева решает покончить с собой, но разговор с подругой подталкивает ее к другому решению - позвонить в опеку и попросить о помощи. Этим звонком Ева предает мать, но одновременно с этим обретает свободу и возможность сделать первый шаг в исполнении своей мечты.",
			Director:        "Мария Свиридкина",
			Producer:        "Илья Минин",
			ContactEmail:    "me@masha.film",
			ContactPhone:    "+7 916 467 13 00", // Same as Director's phone for backward compatibility
			DurationMinutes: 20,
			Genre:           "Драма",
			Locations: []model.Location{
				{
					Name:        "Квартира",
					Description: "В квартире Еве очень одиноко, на что указывает холодное вечернее освещение. Тесное пространство намекает на то, что героине здесь \"тяжело дышать\" во всех смыслах.",
					ImageURL:    "/static/img/location-apartment.webp",
				},
				{
					Name:        "Заброшка",
					Description: "На заброшке Ева \"разрушает фасад\" идеальной семьи. Здесь берет начало разрушение ее старой жизни, которое освобождает место мечтам и надеждам.",
					ImageURL:    "/static/img/location-abandoned.webp",
				},
				{
					Name:        "Артплэй",
					Description: "В Артплэе Ева чувствует себя свободно и счастливо. Рядом с поддерживающим отцом она наконец получает возможность осуществить свои мечты.",
					ImageURL:    "/static/img/location-artplay.webp",
				},
			},
			TeamMembers: []model.TeamMember{
				{Role: "Режиссер", Name: "Мария Свиридкина", Email: "me@masha.film", Phone: "+79164671300", ImageURL: "/static/img/cast/maria.webp"},
				{Role: "Продюсер", Name: "Илья Минин", Email: "minin-ilya@yandex.ru", Phone: "+79110904359", ImageURL: "/static/img/cast/ilya.webp"},
				{Role: "Второй режиссер", Name: "Елизавета Федорова", Email: "", Phone: "", ImageURL: ""},
				{Role: "Директор площадки", Name: "Арина Анисова", Email: "arishasov@gmail.com", Phone: "+79853675219", ImageURL: "/static/img/cast/arina.webp"},
				{Role: "Кастинг-директор", Name: "Максим Головач", Email: "", Phone: "", ImageURL: ""},
				{Role: "Костюмер", Name: "Лола Самадова", Email: "", Phone: "", ImageURL: ""},
				{Role: "Художник-постановщик", Name: "Дарья Зеленкова", Email: "", Phone: "", ImageURL: "/static/img/cast/daria.webp"},
				{Role: "Оператор-постановщик", Name: "Аслан Бададов", Email: "", Phone: "", ImageURL: "/static/img/cast/aslan.webp"},
				{Role: "Фокус-пуллер", Name: "Александр Грушовец", Email: "", Phone: "", ImageURL: ""},
			},
			Cast: []model.CastMember{
				{Role: "Ева", ActorName: "Будет объявлено позже", ImageURL: "/static/img/cast-eva.jpg"},
				{Role: "Бабушка Вера", ActorName: "Будет объявлено позже", ImageURL: "/static/img/cast-grandmother.jpg"},
				{Role: "Мать Оксана", ActorName: "Будет объявлено позже", ImageURL: "/static/img/cast-mother.jpg"},
			},
		},
	}

	// Apply all options
	for _, opt := range opts {
		opt(h)
	}

	return h
}

// HomeHandlerEcho renders the home page.
func (h *Handler) HomeHandlerEcho(c echo.Context) error {
	component := template.Home(h.FilmInfo)
	if err := component.Render(c.Request().Context(), c.Response().Writer); err != nil {
		return handleTemplateError(err, c, "Failed to render home page")
	}
	return nil
}

// AboutHandlerEcho renders the about page.
func (h *Handler) AboutHandlerEcho(c echo.Context) error {
	component := template.About(h.FilmInfo)
	if err := component.Render(c.Request().Context(), c.Response().Writer); err != nil {
		return handleTemplateError(err, c, "Failed to render about page")
	}
	return nil
}

// TeamHandlerEcho renders the team page.
func (h *Handler) TeamHandlerEcho(c echo.Context) error {
	component := template.Team(h.FilmInfo)
	if err := component.Render(c.Request().Context(), c.Response().Writer); err != nil {
		return handleTemplateError(err, c, "Failed to render team page")
	}
	return nil
}

// LocationsHandlerEcho renders the locations page.
func (h *Handler) LocationsHandlerEcho(c echo.Context) error {
	component := template.Locations(h.FilmInfo)
	if err := component.Render(c.Request().Context(), c.Response().Writer); err != nil {
		return handleTemplateError(err, c, "Failed to render locations page")
	}
	return nil
}

// ContactHandlerEcho renders the contact page.
func (h *Handler) ContactHandlerEcho(c echo.Context) error {
	// We need to properly pass the CSRF token to the template
	csrfToken := c.Get("csrf").(string)

	// Create a custom data structure with film info and CSRF token
	data := struct {
		Film      model.FilmInfo
		CSRFToken string
	}{
		Film:      h.FilmInfo,
		CSRFToken: csrfToken,
	}

	// Pass the data to the template
	component := template.ContactWithCSRF(data)
	if err := component.Render(c.Request().Context(), c.Response().Writer); err != nil {
		return handleTemplateError(err, c, "Failed to render contact page")
	}
	return nil
}

// ContactSubmitHandlerEcho processes contact form submissions.
func (h *Handler) ContactSubmitHandlerEcho(c echo.Context) error {
	// Validate CSRF token (handled by middleware)

	// Get form data
	name := strings.TrimSpace(c.FormValue("name"))
	email := strings.TrimSpace(c.FormValue("email"))
	message := strings.TrimSpace(c.FormValue("message"))
	timeStamp := time.Now().Format(time.RFC3339)

	// Validate input
	errors := make(map[string]string)

	if name == "" {
		errors["name"] = "Имя обязательно"
	} else if len(name) > 100 {
		errors["name"] = "Имя слишком длинное (максимум 100 символов)"
	}

	if email == "" {
		errors["email"] = "Email обязателен"
	} else if !validateEmail(email) {
		errors["email"] = "Неверный формат email"
	} else if len(email) > 100 {
		errors["email"] = "Email слишком длинный (максимум 100 символов)"
	}

	if message == "" {
		errors["message"] = "Сообщение обязательно"
	} else if len(message) > 5000 {
		errors["message"] = "Сообщение слишком длинное (максимум 5000 символов)"
	}

	if len(errors) > 0 {
		// For HTMX requests, return form with errors
		c.Response().Header().Set("HX-Trigger", "{\"showFormErrors\": true}")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"errors": errors,
		})
	}

	// Sanitize input for logging/telegram
	name = sanitizeString(name)
	email = sanitizeString(email)
	message = sanitizeString(message)

	// Log the submission
	slog.Info("Contact form submission",
		"name", name,
		"email", email,
		"message_length", len(message),
		"time", timeStamp)

	// Format message for Telegram
	telegramMsg := fmt.Sprintf(
		"<b>Новое сообщение с сайта!</b>\n\n"+
			"<b>Имя:</b> %s\n"+
			"<b>Email:</b> %s\n"+
			"<b>Время:</b> %s\n\n"+
			"<b>Сообщение:</b>\n%s",
		html.EscapeString(name), html.EscapeString(email), timeStamp, html.EscapeString(message))

	// Send to Telegram
	err := h.sendTelegramMessage(telegramMsg)
	if err != nil {
		slog.Error("Failed to send message to Telegram", "error", err)
		// Continue anyway - don't show error to user
	}

	// Return success template
	component := template.ContactSuccess()
	if err := component.Render(c.Request().Context(), c.Response().Writer); err != nil {
		return handleTemplateError(err, c, "Failed to render contact success page")
	}
	return nil
}

// sendTelegramMessage sends a message to the specified Telegram chat
func (h *Handler) sendTelegramMessage(text string) error {
	// Check if Telegram is configured
	if h.TelegramToken == "" || h.TelegramChatID == "" {
		slog.Info("Telegram notification skipped - token or chat ID not configured")
		return nil
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", h.TelegramToken)

	msg := TelegramMessage{
		ChatID:    h.TelegramChatID,
		Text:      text,
		ParseMode: "HTML", // Allow HTML formatting
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		slog.Error("Failed to marshal Telegram message", "error", err)
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		slog.Error("Failed to send Telegram message", "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Telegram API error", "status", resp.Status)
		return fmt.Errorf("telegram API error: %s", resp.Status)
	}

	return nil
}

// validateEmail and sanitizeString functions are in validation.go

// handleTemplateError logs the error and returns a proper HTTP error response.
func handleTemplateError(err error, _ echo.Context, message string) error {
	slog.Error(message, "error", err)
	return echo.NewHTTPError(500, "Internal Server Error")
}
