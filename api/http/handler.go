package httpHuvalk

import (
	"encoding/json"
	"github.com/huvalk/tech-db-huvalk/api/models"
	"github.com/huvalk/tech-db-huvalk/api/repository"
	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
	"net/http"
	"strings"
)

type Handler struct {
	repo *repository.PostgresRepository
}

func NewHandler(r *repository.PostgresRepository) *Handler {
	return &Handler{
		repo: r,
	}
}

var err = models.Error{Message: "Can't find user with id\n"}

func (h *Handler) ForumCreate(c echo.Context) error {
	var newForum models.Forum
	//json.NewDecoder(c.Request().Body).Decode(&newForum)
	easyjson.UnmarshalFromReader(c.Request().Body, &newForum)
	json.NewDecoder(c.Request().Body).Decode(&newForum)
	c.Response().Header().Set("Content-Type", "application/json")

	result, resultStatus := h.repo.CreateForum(&newForum)

	switch resultStatus {
	case 201:
		//resultJSON, _ := newForum.MarshalJSON()
		//return c.JSON(resultStatus, newForum)
		c.Response().WriteHeader(http.StatusCreated)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(newForum, c.Response())
		return errResult
	case 404:
		//resultJSON, _ := json.Marshal(err)
		//return c.JSON(resultStatus, err)
		c.Response().WriteHeader(http.StatusNotFound)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(err, c.Response())
		return errResult
	case 409:
		//resultJSON, _ := json.Marshal(*resultForum)
		//return c.JSON(resultStatus, resultForum)
		c.Response().WriteHeader(http.StatusConflict)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(result, c.Response())
		return errResult
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) ForumCreateThread(c echo.Context) error {
	var newThread models.Thread
	//json.NewDecoder(c.Request().Body).Decode(&newThread)
	easyjson.UnmarshalFromReader(c.Request().Body, &newThread)
	newThread.Forum = c.Param("slug")
	c.Response().Header().Set("Content-Type", "application/json")

	result, resultStatus := h.repo.CreateThread(&newThread)

	switch resultStatus {
	case 201:
		//resultJSON, _ := json.Marshal(*resultThread)
		//return c.JSON(resultStatus, resultThread)
		c.Response().WriteHeader(http.StatusCreated)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(result, c.Response())
		return errResult
	case 404:
		//resultJSON, _ := json.Marshal(err)
		//return c.JSON(resultStatus, err)
		c.Response().WriteHeader(http.StatusNotFound)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(err, c.Response())
		return errResult
	case 409:
		//resultJSON, _ := json.Marshal(*resultThread)
		//return c.JSON(resultStatus, resultThread)
		c.Response().WriteHeader(http.StatusConflict)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(result, c.Response())
		return errResult
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) ForumGet(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")

	result, resultStatus := h.repo.GetForum(c.Param("slug"))

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(*resultForum)
		//return c.JSON(resultStatus, resultForum)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(result, c.Response())
		return errResult
	case 404:
		//resultJSON, _ := json.Marshal(err)
		//return c.JSON(resultStatus, err)
		c.Response().WriteHeader(http.StatusNotFound)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(err, c.Response())
		return errResult
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) ForumGetListOfThreads(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	slug := c.Param("slug")
	limit := c.QueryParam("limit")
	since := c.QueryParam("since")
	desc := c.QueryParam("desc")

	result, resultStatus := h.repo.GetListOfThreads(slug, limit, since, desc)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(resultForum)
		//return c.JSON(resultStatus, resultForum)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(result, c.Response())
		return errResult
	case 404:
		//resultJSON, _ := json.Marshal(err)
		//return c.JSON(resultStatus, err)
		c.Response().WriteHeader(http.StatusNotFound)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(err, c.Response())
		return errResult
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) ForumGetListOfUsers(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	slug := c.Param("slug")
	limit := c.QueryParam("limit")
	since := c.QueryParam("since")
	desc := c.QueryParam("desc")

	result, resultStatus := h.repo.GetListOfUsers(slug, limit, since, desc)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(resultForum)
		//return c.JSON(resultStatus, resultForum)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(result, c.Response())
		return errResult
	case 404:
		//resultJSON, _ := json.Marshal(err)
		//return c.JSON(resultStatus, err)
		c.Response().WriteHeader(http.StatusNotFound)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(err, c.Response())
		return errResult
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) PostChangeDetails(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	slug := c.Param("id")
	var newPost models.PostUpdate
	//json.NewDecoder(c.Request().Body).Decode(&newPost)
	easyjson.UnmarshalFromReader(c.Request().Body, &newPost)

	result, resultStatus := h.repo.ChangePost(slug, &newPost)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(resultThread)
		//return c.JSON(resultStatus, resultThread)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(result, c.Response())
		return errResult
	case 404:
		//resultJSON, _ := json.Marshal(err)
		//return c.JSON(resultStatus, err)
		c.Response().WriteHeader(http.StatusNotFound)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(err, c.Response())
		return errResult
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) PostGetDetails(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	slug := c.Param("id")
	related := strings.Split(c.QueryParam("related"), ",")

	result, resultStatus := h.repo.PostDetails(slug, related)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(resultForum)
		//return c.JSON(resultStatus, resultForum)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(result, c.Response())
		return errResult
	case 404:
		//resultJSON, _ := json.Marshal(err)
		//return c.JSON(resultStatus, err)
		c.Response().WriteHeader(http.StatusNotFound)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(err, c.Response())
		return errResult
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) ServiceGetStatus(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	result, resultStatus := h.repo.GetStatus()

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(*resultStat)
		//return c.JSON(resultStatus, resultStat)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(result, c.Response())
		return errResult
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) ServiceClear(c echo.Context) error {
	resultStatus := h.repo.ClearAll()

	switch resultStatus {
	case 200:
		return c.NoContent(resultStatus)
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) PostsCreate(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	slug := c.Param("slug_or_id")
	var newPosts models.Posts
	easyjson.UnmarshalFromReader(c.Request().Body, &newPosts)
	//json.NewDecoder(c.Request().Body).Decode(&newPosts)

	result, resultStatus := h.repo.CreatePosts(slug, newPosts)

	switch resultStatus {
	case 201:
		//resultJSON, _ := json.Marshal(resultThread)
		//return c.JSON(resultStatus, resultThread)
		c.Response().WriteHeader(http.StatusCreated)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(result, c.Response())
		return errResult
	case 404:
		//resultJSON, _ := json.Marshal(err)
		//return c.JSON(resultStatus, err)
		c.Response().WriteHeader(http.StatusNotFound)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(err, c.Response())
		return errResult
	case 409:
		//resultJSON, _ := json.Marshal(err)
		//return c.JSON(resultStatus, err)
		c.Response().WriteHeader(http.StatusConflict)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(err, c.Response())
		return errResult
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) ThreadGetDetales(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	slug := c.Param("slug_or_id")
	var newPosts models.Posts
	easyjson.UnmarshalFromReader(c.Request().Body, &newPosts)
	//json.NewDecoder(c.Request().Body).Decode(&newPosts)

	result, resultStatus := h.repo.GetThread(slug)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(resultThread)
		//return c.JSON(resultStatus, resultThread)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(result, c.Response())
		return errResult
	case 404:
		//resultJSON, _ := json.Marshal(err)
		//return c.JSON(resultStatus, err)
		c.Response().WriteHeader(http.StatusNotFound)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(err, c.Response())
		return errResult
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) ThreadChange(c echo.Context) error {
	var changeThread models.ThreadUpdate
	easyjson.UnmarshalFromReader(c.Request().Body, &changeThread)
	//json.NewDecoder(c.Request().Body).Decode(&changeThread)
	slug := c.Param("slug_or_id")
	c.Response().Header().Set("Content-Type", "application/json")

	result, resultStatus := h.repo.ChangeThread(slug, &changeThread)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(*resultThread)
		//return c.JSON(resultStatus, resultThread)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(result, c.Response())
		return errResult
	case 404:
		//resultJSON, _ := json.Marshal(err)
		//return c.JSON(resultStatus, err)
		c.Response().WriteHeader(http.StatusNotFound)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(err, c.Response())
		return errResult
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) ThreadGetListOfPost(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	slug := c.Param("slug_or_id")
	limit := c.QueryParam("limit")
	since := c.QueryParam("since")
	desc := c.QueryParam("desc")
	sort := c.QueryParam("sort")

	result, resultStatus := h.repo.GetListOfPosts(slug, limit, since, sort, desc)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(resultForum)
		//return c.JSON(resultStatus, resultForum)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(result, c.Response())
		return errResult
	case 404:
		//resultJSON, _ := json.Marshal(err)
		//return c.JSON(resultStatus, err)
		c.Response().WriteHeader(http.StatusNotFound)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(err, c.Response())
		return errResult
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) ThreadVote(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	slug := c.Param("slug_or_id")
	var newVote models.Vote
	easyjson.UnmarshalFromReader(c.Request().Body, &newVote)
	//json.NewDecoder(c.Request().Body).Decode(&newVote)

	result, resultStatus := h.repo.VoteForThread(slug, &newVote)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(resultThread)
		//return c.JSON(resultStatus, resultThread)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(result, c.Response())
		return errResult
	case 404:
		//resultJSON, _ := json.Marshal(err)
		//return c.JSON(resultStatus, err)
		c.Response().WriteHeader(http.StatusNotFound)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(err, c.Response())
		return errResult
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) UserCreate(c echo.Context) error {
	var newUser models.User
	easyjson.UnmarshalFromReader(c.Request().Body, &newUser)
	//json.NewDecoder(c.Request().Body).Decode(&newUser)
	newUser.Nickname = c.Param("nickname")
	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	result, resultStatus := h.repo.CreateUser(&newUser)

	switch resultStatus {
	case 201:
		//resultJSON, _ := json.Marshal(newUser)
		//return c.JSON(resultStatus, newUser)
		c.Response().WriteHeader(http.StatusCreated)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(newUser, c.Response())
		return errResult
	case 409:
		//resultJSON, _ := json.Marshal(resultUsers)
		//return c.JSON(resultStatus, resultUsers)
		c.Response().WriteHeader(http.StatusConflict)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(result, c.Response())
		return errResult
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) UserGet(c echo.Context) error {
	userNickname := c.Param("nickname")
	c.Response().Header().Set("Content-Type", "application/json")

	result, resultStatus := h.repo.GetUser(userNickname)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(*resultUser)
		//return c.JSON(resultStatus, resultUser)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(result, c.Response())
		return errResult
	case 404:
		//resultJSON, _ := json.Marshal(err)
		//return c.JSON(resultStatus, err)
		c.Response().WriteHeader(http.StatusNotFound)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(err, c.Response())
		return errResult
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) UserChange(c echo.Context) error {
	var changeUser models.UserUpdate
	easyjson.UnmarshalFromReader(c.Request().Body, &changeUser)
	//json.NewDecoder(c.Request().Body).Decode(&changeUser)
	userNickname := c.Param("nickname")

	result, resultStatus := h.repo.ChangeUser(userNickname, &changeUser)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(*resultUser)
		//return c.JSON(resultStatus, resultUser)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(result, c.Response())
		return errResult
	case 404:
		//resultJSON, _ := json.Marshal(err)
		//TODO что за ошибка?
		return c.JSON(resultStatus, err)
		//c.Response().WriteHeader(http.StatusNotFound)
		//_, _, errResult := easyjson.MarshalToHTTPResponseWriter(err, c.Response())
		//return errResult
	case 409:
		//resultJSON, _ := json.Marshal(err)
		//return c.JSON(resultStatus, err)
		c.Response().WriteHeader(http.StatusConflict)
		_, _, errResult := easyjson.MarshalToHTTPResponseWriter(err, c.Response())
		return errResult
	default:
		return echo.ErrInternalServerError
	}
}
