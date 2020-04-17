package httpHuvalk

import (
	"encoding/json"
	"github.com/huvalk/tech-db-huvalk/api/models"
	"github.com/huvalk/tech-db-huvalk/api/repository"
	"github.com/labstack/echo/v4"
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
	json.NewDecoder(c.Request().Body).Decode(&newForum)
	c.Response().Header().Set("Content-Type", "application/json")

	resultForum, resultStatus := h.repo.CreateForum(&newForum)

	switch resultStatus {
	case 201:
		//resultJSON, _ := json.Marshal(newForum)
		return c.JSON(resultStatus, newForum)
	case 404:
		//resultJSON, _ := json.Marshal(err)
		return c.JSON(resultStatus, err)
	case 409:
		//resultJSON, _ := json.Marshal(*resultForum)
		return c.JSON(resultStatus, resultForum)
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) ForumCreateThread(c echo.Context) error {
	var newThread models.Thread
	json.NewDecoder(c.Request().Body).Decode(&newThread)
	newThread.Forum = c.Param("slug")
	c.Response().Header().Set("Content-Type", "application/json")

	resultThread, resultStatus := h.repo.CreateThread(&newThread)

	switch resultStatus {
	case 201:
		//resultJSON, _ := json.Marshal(*resultThread)
		return c.JSON(resultStatus, resultThread)
	case 404:
		//resultJSON, _ := json.Marshal(err)
		return c.JSON(resultStatus, err)
	case 409:
		//resultJSON, _ := json.Marshal(*resultThread)
		return c.JSON(resultStatus, resultThread)
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) ForumGet(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")

	resultForum, resultStatus := h.repo.GetForum(c.Param("slug"))

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(*resultForum)
		return c.JSON(resultStatus, resultForum)
	case 404:
		//resultJSON, _ := json.Marshal(err)
		return c.JSON(resultStatus, err)
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

	resultForum, resultStatus := h.repo.GetListOfThreads(slug, limit, since, desc)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(resultForum)
		return c.JSON(resultStatus, resultForum)
	case 404:
		//resultJSON, _ := json.Marshal(err)
		return c.JSON(resultStatus, err)
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

	resultForum, resultStatus := h.repo.GetListOfUsers(slug, limit, since, desc)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(resultForum)
		return c.JSON(resultStatus, resultForum)
	case 404:
		//resultJSON, _ := json.Marshal(err)
		return c.JSON(resultStatus, err)
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) PostChangeDetails(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	slug := c.Param("id")
	var newPost models.PostUpdate
	json.NewDecoder(c.Request().Body).Decode(&newPost)

	resultThread, resultStatus := h.repo.ChangePost(slug, &newPost)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(resultThread)
		return c.JSON(resultStatus, resultThread)
	case 404:
		//resultJSON, _ := json.Marshal(err)
		return c.JSON(resultStatus, err)
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) PostGetDetails(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	slug := c.Param("id")
	related := strings.Split(c.QueryParam("related"), ",")

	resultForum, resultStatus := h.repo.PostDetails(slug, related)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(resultForum)
		return c.JSON(resultStatus, resultForum)
	case 404:
		//resultJSON, _ := json.Marshal(err)
		return c.JSON(resultStatus, err)
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) ServiceGetStatus(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	resultStat, resultStatus := h.repo.GetStatus()

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(*resultStat)
		return c.JSON(resultStatus, resultStat)
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
	json.NewDecoder(c.Request().Body).Decode(&newPosts)

	resultThread, resultStatus := h.repo.CreatePosts(slug, newPosts)

	switch resultStatus {
	case 201:
		//resultJSON, _ := json.Marshal(resultThread)
		return c.JSON(resultStatus, resultThread)
	case 404:
		//resultJSON, _ := json.Marshal(err)
		return c.JSON(resultStatus, err)
	case 409:
		//resultJSON, _ := json.Marshal(err)
		return c.JSON(resultStatus, err)
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) ThreadGetDetales(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	slug := c.Param("slug_or_id")
	var newPosts models.Posts
	json.NewDecoder(c.Request().Body).Decode(&newPosts)

	resultThread, resultStatus := h.repo.GetThread(slug)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(resultThread)
		return c.JSON(resultStatus, resultThread)
	case 404:
		//resultJSON, _ := json.Marshal(err)
		return c.JSON(resultStatus, err)
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) ThreadChange(c echo.Context) error {
	var changeThread models.ThreadUpdate
	json.NewDecoder(c.Request().Body).Decode(&changeThread)
	slug := c.Param("slug_or_id")
	c.Response().Header().Set("Content-Type", "application/json")

	resultThread, resultStatus := h.repo.ChangeThread(slug, &changeThread)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(*resultThread)
		return c.JSON(resultStatus, resultThread)
	case 404:
		//resultJSON, _ := json.Marshal(err)
		return c.JSON(resultStatus, err)
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

	resultForum, resultStatus := h.repo.GetListOfPosts(slug, limit, since, sort, desc)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(resultForum)
		return c.JSON(resultStatus, resultForum)
	case 404:
		//resultJSON, _ := json.Marshal(err)
		return c.JSON(resultStatus, err)
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) ThreadVote(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	slug := c.Param("slug_or_id")
	var newVote models.Vote
	json.NewDecoder(c.Request().Body).Decode(&newVote)

	resultThread, resultStatus := h.repo.VoteForThread(slug, &newVote)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(resultThread)
		return c.JSON(resultStatus, resultThread)
	case 404:
		//resultJSON, _ := json.Marshal(err)
		return c.JSON(resultStatus, err)
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) UserCreate(c echo.Context) error {
	userNickname := c.Param("nickname")

	var newUser models.User
	json.NewDecoder(c.Request().Body).Decode(&newUser)
	newUser.Nickname = userNickname
	c.Response().Header().Set("Content-Type", "application/json")

	resultUsers, resultStatus := h.repo.CreateUser(&newUser)

	switch resultStatus {
	case 201:
		//resultJSON, _ := json.Marshal(newUser)
		return c.JSON(resultStatus, newUser)
	case 409:
		//resultJSON, _ := json.Marshal(resultUsers)
		return c.JSON(resultStatus, resultUsers)
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) UserGet(c echo.Context) error {
	userNickname := c.Param("nickname")
	c.Response().Header().Set("Content-Type", "application/json")

	resultUser, resultStatus := h.repo.GetUser(userNickname)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(*resultUser)
		return c.JSON(resultStatus, resultUser)
	case 404:
		//resultJSON, _ := json.Marshal(err)
		return c.JSON(resultStatus, err)
	default:
		return echo.ErrInternalServerError
	}
}

func (h *Handler) UserChange(c echo.Context) error {
	var changeUser models.UserUpdate
	json.NewDecoder(c.Request().Body).Decode(&changeUser)
	userNickname := c.Param("nickname")
	c.Response().Header().Set("Content-Type", "application/json")

	resultUser, resultStatus := h.repo.ChangeUser(userNickname, &changeUser)

	switch resultStatus {
	case 200:
		//resultJSON, _ := json.Marshal(*resultUser)
		return c.JSON(resultStatus, resultUser)
	case 404:
		//resultJSON, _ := json.Marshal(err)
		return c.JSON(resultStatus, err)
	case 409:
		//resultJSON, _ := json.Marshal(err)
		return c.JSON(resultStatus, err)
	default:
		return echo.ErrInternalServerError
	}
}
