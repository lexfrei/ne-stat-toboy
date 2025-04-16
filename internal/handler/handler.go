// Package handler provides HTTP handlers for the "Ne Stat Toboy" film website.
package handler

import (
	"log/slog"
	"github.com/labstack/echo/v4"
	"github.com/lexfrei/ne-stat-toboy/internal/model"
	"github.com/lexfrei/ne-stat-toboy/web/template"
)

// Handler contains all HTTP handlers and their dependencies.
type Handler struct {
	FilmInfo model.FilmInfo
}

// New creates a new Handler with initialized dependencies.
func New() *Handler {
	return &Handler{
		FilmInfo: model.FilmInfo{
			Title:           "Не стать тобой",
			Tagline:         "Жизнь в семье, где любовь выражается насилием, ложью и запретами, подталкивает школьницу к тяжелому выбору - смерть или предательство ради жизни",
			Synopsis:        "Ева мечтает стать музыкантом, но её жизнь жёстко контролируют мать и бабушка, запрещая любые увлечения и навязывая ей своё будущее. Отец ушёл, но говорить об этом нельзя — мать создаёт видимость «идеальной семьи». Дома же Еву наказывают за любые ошибки. Когда после визита подруги Вики у Евы пропадают наушники - подарок отца, бабушка обвиняет в этом Вику, но Ева случайно находит их в вещах бабушки и начинает сомневаться в семье. Вика убеждает Еву впервые нарушить правила и прогулять школу. Девушки веселятся, и Ева раскрывает правду о семье, не зная, что её снимают на видео. Запись попадает в сеть, начинается травля и сильные побои дома. После побоев Ева решает покончить с собой, но разговор с подругой подталкивает ее к другому решению - позвонить в опеку и попросить о помощи. Этим звонком Ева предает мать, но одновременно с этим обретает свободу и возможность сделать первый шаг в исполнении своей мечты.",
			Director:        "Мария Свиридкина",
			Producer:        "Илья Минин",
			ContactEmail:    "me@masha.film",
			ContactPhone:    "+7 916 467 13 00",
			DurationMinutes: 20,
			Genre:           "Драма",
			Locations: []model.Location{
				{
					Name:        "Квартира",
					Description: "В квартире Еве очень одиноко, на что указывает холодное вечернее освещение. Тесное пространство намекает на то, что героине здесь \"тяжело дышать\" во всех смыслах.",
					ImageURL:    "/static/img/location-apartment.jpg",
				},
				{
					Name:        "Заброшка",
					Description: "На заброшке Ева \"разрушает фасад\" идеальной семьи. Здесь берет начало разрушение ее старой жизни, которое освобождает место мечтам и надеждам.",
					ImageURL:    "/static/img/location-abandoned.jpg",
				},
				{
					Name:        "Артплэй",
					Description: "В Артплэе Ева чувствует себя свободно и счастливо. Рядом с поддерживающим отцом она наконец получает возможность осуществить свои мечты.",
					ImageURL:    "/static/img/location-artplay.jpg",
				},
			},
			TeamMembers: []model.TeamMember{
				{Role: "Режиссер", Name: "Мария Свиридкина", Email: "me@masha.film"},
				{Role: "Продюсер", Name: "Илья Минин", Email: ""},
				{Role: "Второй режиссер", Name: "Елизавета Федорова", Email: ""},
				{Role: "Директор площадки", Name: "Арина Анисова", Email: ""},
				{Role: "Кастинг-директор", Name: "Максим Головач", Email: ""},
				{Role: "Костюмер", Name: "Лола Самадова", Email: ""},
				{Role: "Художник-постановщик", Name: "Дарья Зеленкова", Email: ""},
				{Role: "Оператор-постановщик", Name: "Аслан Бададов", Email: ""},
				{Role: "Фокус-пуллер", Name: "Александр Грушовец", Email: ""},
			},
			Cast: []model.CastMember{
				{Role: "Ева", ActorName: "Кира Аллен", ImageURL: "/static/img/cast-eva.jpg"},
				{Role: "Бабушка Вера", ActorName: "Светлана Крючкова", ImageURL: "/static/img/cast-grandmother.jpg"},
				{Role: "Мать Оксана", ActorName: "Наталья Тетенова", ImageURL: "/static/img/cast-mother.jpg"},
			},
		},
	}
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
	component := template.Contact(h.FilmInfo)
	if err := component.Render(c.Request().Context(), c.Response().Writer); err != nil {
		return handleTemplateError(err, c, "Failed to render contact page")
	}
	return nil
}

// ContactSubmitHandlerEcho processes contact form submissions.
func (h *Handler) ContactSubmitHandlerEcho(c echo.Context) error {
	// In a real app, process the form submission here
	component := template.ContactSuccess()
	if err := component.Render(c.Request().Context(), c.Response().Writer); err != nil {
		return handleTemplateError(err, c, "Failed to render contact success page")
	}
	return nil
}

// handleTemplateError logs the error and returns a proper HTTP error response.
func handleTemplateError(err error, c echo.Context, message string) error {
	slog.Error(message, "error", err)
	return echo.NewHTTPError(500, "Internal Server Error")
}