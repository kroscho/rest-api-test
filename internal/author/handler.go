package author

import (
	"context"
	"encoding/json"
	"net/http"
	"rest-api-test/internal/apperror"
	"rest-api-test/internal/handlers"
	"rest-api-test/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

var _ handlers.Handler = &handler{}

const (
	authorsURL = "/authors"
	authorURL  = "/authors/:uuid"
)

// нужно для обеспечения процессов создания, получения списка и тд
type handler struct {
	logger     *logging.Logger
	repository Repository
}

func NewHandler(repository Repository, logger *logging.Logger) handlers.Handler {
	return &handler{
		repository: repository,
		logger:     logger,
	}
}

// регистрация handler в router
func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, authorsURL, apperror.Middleware(h.CreateAuthor))
	router.HandlerFunc(http.MethodGet, authorsURL, apperror.Middleware(h.GetListAuthors))
	router.HandlerFunc(http.MethodGet, authorURL, apperror.Middleware(h.GetAuthorByID))
	router.HandlerFunc(http.MethodPut, authorsURL, apperror.Middleware(h.UpdateAuthor))
}

func (h *handler) CreateAuthor(w http.ResponseWriter, r *http.Request) error {
	var author Author

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&author)
	if err != nil {
		return err
	}

	err = h.repository.Create(context.TODO(), &author)
	if err != nil {
		return err
	}

	oneBytes, err := json.Marshal(author)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(oneBytes)

	return nil
}

func (h *handler) GetListAuthors(w http.ResponseWriter, r *http.Request) error {
	all, err := h.repository.FindAll(context.TODO())
	if err != nil {
		w.WriteHeader(400)
		return err
	}

	allBytes, err := json.Marshal(all)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(allBytes)

	return nil
}

func (h *handler) GetAuthorByID(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())

	id := params.ByName("uuid")

	one, err := h.repository.FindOne(context.TODO(), id)
	if err != nil {
		w.WriteHeader(400)
		return err
	}

	oneBytes, err := json.Marshal(one)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(oneBytes)

	return nil
}

func (h *handler) UpdateAuthor(w http.ResponseWriter, r *http.Request) error {
	var author Author

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&author)
	if err != nil {
		return err
	}

	err = h.repository.Update(context.TODO(), author)
	if err != nil {
		return err
	}

	oneBytes, err := json.Marshal(author)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(oneBytes)

	return nil
}
