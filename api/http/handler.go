package httpHuvalk

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/huvalk/tech-db-huvalk/api/models"
	"github.com/huvalk/tech-db-huvalk/api/repository"
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

func (h *Handler) ForumCreate(w http.ResponseWriter, r *http.Request) {
	var newForum models.Forum
	json.NewDecoder(r.Body).Decode(&newForum)
	w.Header().Set("Content-Type", "application/json")

	resultForum, resultStatus := h.repo.CreateForum(&newForum)

	switch resultStatus {
	case 201:
		resultJSON, _ := json.Marshal(newForum)
		w.WriteHeader(http.StatusCreated)
		w.Write(resultJSON)
	case 404:
		resultJSON, _ := json.Marshal(models.Error{Message: "Can't find user with id\n"})
		w.WriteHeader(http.StatusNotFound)
		w.Write(resultJSON)
	case 409:
		resultJSON, _ := json.Marshal(*resultForum)
		w.WriteHeader(http.StatusConflict)
		w.Write(resultJSON)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) ForumCreateThread(w http.ResponseWriter, r *http.Request) {
	var newThread models.Thread
	json.NewDecoder(r.Body).Decode(&newThread)
	newThread.Forum = mux.Vars(r)["slug"]
	w.Header().Set("Content-Type", "application/json")

	resultThread, resultStatus := h.repo.CreateThread(&newThread)

	switch resultStatus {
	case 201:
		resultJSON, _ := json.Marshal(*resultThread)
		w.WriteHeader(http.StatusCreated)
		w.Write(resultJSON)
	case 404:
		resultJSON, _ := json.Marshal(models.Error{Message: "Can't find user with id\n"})
		w.WriteHeader(http.StatusNotFound)
		w.Write(resultJSON)
	case 409:
		resultJSON, _ := json.Marshal(*resultThread)
		w.WriteHeader(http.StatusConflict)
		w.Write(resultJSON)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) ForumGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resultForum, resultStatus := h.repo.GetForum(mux.Vars(r)["slug"])

	switch resultStatus {
	case 200:
		resultJSON, _ := json.Marshal(*resultForum)
		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	case 404:
		resultJSON, _ := json.Marshal(models.Error{Message: "Can't find user with id\n"})
		w.WriteHeader(http.StatusNotFound)
		w.Write(resultJSON)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) ForumGetListOfThreads(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := mux.Vars(r)["slug"]
	limit := r.FormValue("limit")
	since := r.FormValue("since")
	desc := r.FormValue("desc")

	resultForum, resultStatus := h.repo.GetListOfThreads(slug, limit, since, desc)

	switch resultStatus {
	case 200:
		resultJSON, _ := json.Marshal(resultForum)
		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	case 404:
		resultJSON, _ := json.Marshal(models.Error{Message: "Can't find threads\n"})
		w.WriteHeader(http.StatusNotFound)
		w.Write(resultJSON)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) ForumGetListOfUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := mux.Vars(r)["slug"]
	limit := r.FormValue("limit")
	since := r.FormValue("since")
	desc := r.FormValue("desc")

	resultForum, resultStatus := h.repo.GetListOfUsers(slug, limit, since, desc)

	switch resultStatus {
	case 200:
		resultJSON, _ := json.Marshal(resultForum)
		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	case 404:
		resultJSON, _ := json.Marshal(models.Error{Message: "Can't find forum\n"})
		w.WriteHeader(http.StatusNotFound)
		w.Write(resultJSON)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) PostChangeDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := mux.Vars(r)["id"]
	var newPost models.PostUpdate
	json.NewDecoder(r.Body).Decode(&newPost)

	resultThread, resultStatus := h.repo.ChangePost(slug, &newPost)

	switch resultStatus {
	case 200:
		resultJSON, _ := json.Marshal(resultThread)
		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	case 404:
		resultJSON, _ := json.Marshal(models.Error{Message: "Can't find post\n"})
		w.WriteHeader(http.StatusNotFound)
		w.Write(resultJSON)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) PostGetDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := mux.Vars(r)["id"]
	related := strings.Split(r.FormValue("related"), ",")

	resultForum, resultStatus := h.repo.PostDetails(slug, related)

	switch resultStatus {
	case 200:
		resultJSON, _ := json.Marshal(resultForum)
		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	case 404:
		resultJSON, _ := json.Marshal(models.Error{Message: "Can't find forum\n"})
		w.WriteHeader(http.StatusNotFound)
		w.Write(resultJSON)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) ServiceGetStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resultStat, resultStatus := h.repo.GetStatus()

	switch resultStatus {
	case 200:
		resultJSON, _ := json.Marshal(*resultStat)
		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) ServiceClear(w http.ResponseWriter, r *http.Request) {
	resultStatus := h.repo.ClearAll()

	switch resultStatus {
	case 200:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) PostsCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := mux.Vars(r)["slug_or_id"]
	var newPosts models.Posts
	json.NewDecoder(r.Body).Decode(&newPosts)

	resultThread, resultStatus := h.repo.CreatePosts(slug, newPosts)

	switch resultStatus {
	case 201:
		resultJSON, _ := json.Marshal(resultThread)
		w.WriteHeader(http.StatusCreated)
		w.Write(resultJSON)
	case 404:
		resultJSON, _ := json.Marshal(models.Error{Message: "Can't find threads\n"})
		w.WriteHeader(http.StatusNotFound)
		w.Write(resultJSON)
	case 409:
		resultJSON, _ := json.Marshal(models.Error{Message: "Can't find parent post\n"})
		w.WriteHeader(http.StatusConflict)
		w.Write(resultJSON)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) ThreadGetDetales(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := mux.Vars(r)["slug_or_id"]
	var newPosts models.Posts
	json.NewDecoder(r.Body).Decode(&newPosts)

	resultThread, resultStatus := h.repo.GetThread(slug)

	switch resultStatus {
	case 200:
		resultJSON, _ := json.Marshal(resultThread)
		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	case 404:
		resultJSON, _ := json.Marshal(models.Error{Message: "Can't find thread\n"})
		w.WriteHeader(http.StatusNotFound)
		w.Write(resultJSON)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) ThreadChange(w http.ResponseWriter, r *http.Request) {
	var changeThread models.ThreadUpdate
	json.NewDecoder(r.Body).Decode(&changeThread)
	slug := mux.Vars(r)["slug_or_id"]
	w.Header().Set("Content-Type", "application/json")

	resultThread, resultStatus := h.repo.ChangeThread(slug, &changeThread)

	switch resultStatus {
	case 200:
		resultJSON, _ := json.Marshal(*resultThread)
		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	case 404:
		resultJSON, _ := json.Marshal(models.Error{Message: "Can't find thread with id\n"})
		w.WriteHeader(http.StatusNotFound)
		w.Write(resultJSON)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) ThreadGetListOfPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := mux.Vars(r)["slug_or_id"]
	limit := r.FormValue("limit")
	since := r.FormValue("since")
	desc := r.FormValue("desc")
	sort := r.FormValue("sort")

	resultForum, resultStatus := h.repo.GetListOfPosts(slug, limit, since, sort, desc)

	switch resultStatus {
	case 200:
		resultJSON, _ := json.Marshal(resultForum)
		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	case 404:
		resultJSON, _ := json.Marshal(models.Error{Message: "Can't find threads\n"})
		w.WriteHeader(http.StatusNotFound)
		w.Write(resultJSON)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) ThreadVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := mux.Vars(r)["slug_or_id"]
	var newVote models.Vote
	json.NewDecoder(r.Body).Decode(&newVote)

	resultThread, resultStatus := h.repo.VoteForThread(slug, &newVote)

	switch resultStatus {
	case 200:
		resultJSON, _ := json.Marshal(resultThread)
		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	case 404:
		resultJSON, _ := json.Marshal(models.Error{Message: "Can't find threads\n"})
		w.WriteHeader(http.StatusNotFound)
		w.Write(resultJSON)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) UserCreate(w http.ResponseWriter, r *http.Request) {
	userNickname := mux.Vars(r)["nickname"]

	var newUser models.User
	json.NewDecoder(r.Body).Decode(&newUser)
	newUser.Nickname = userNickname
	w.Header().Set("Content-Type", "application/json")

	resultUsers, resultStatus := h.repo.CreateUser(&newUser)

	switch resultStatus {
	case 201:
		resultJSON, _ := json.Marshal(newUser)

		w.WriteHeader(http.StatusCreated)
		w.Write(resultJSON)

	case 409:
		resultJSON, _ := json.Marshal(resultUsers)
		w.WriteHeader(http.StatusConflict)
		w.Write(resultJSON)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) UserGet(w http.ResponseWriter, r *http.Request) {
	userNickname := mux.Vars(r)["nickname"]
	w.Header().Set("Content-Type", "application/json")

	resultUser, resultStatus := h.repo.GetUser(userNickname)

	switch resultStatus {
	case 200:
		resultJSON, _ := json.Marshal(*resultUser)
		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	case 404:
		resultJSON, _ := json.Marshal(models.Error{Message: "Can't find user with id\n"})
		w.WriteHeader(http.StatusNotFound)
		w.Write(resultJSON)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) UserChange(w http.ResponseWriter, r *http.Request) {
	var changeUser models.UserUpdate
	json.NewDecoder(r.Body).Decode(&changeUser)
	userNickname := mux.Vars(r)["nickname"]
	w.Header().Set("Content-Type", "application/json")

	resultUser, resultStatus := h.repo.ChangeUser(userNickname, &changeUser)

	switch resultStatus {
	case 200:
		resultJSON, _ := json.Marshal(*resultUser)
		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	case 404:
		resultJSON, _ := json.Marshal(models.Error{Message: "Can't find user with id\n"})
		w.WriteHeader(http.StatusNotFound)
		w.Write(resultJSON)
	case 409:
		resultJSON, _ := json.Marshal(models.Error{Message: "Email already in use\n"})
		w.WriteHeader(http.StatusConflict)
		w.Write(resultJSON)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
