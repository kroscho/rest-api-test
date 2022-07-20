package apperror

import (
	"errors"
	"net/http"
)

type apphandler func(w http.ResponseWriter, r *http.Request) error

func Middleware(h apphandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var appErr *AppError
		err := h(w, r)
		if err != nil {
			if errors.As(err, &appErr) {
				w.Header().Set("Content-Type", "application/json")
				if errors.Is(err, ErrNotFound) {
					w.WriteHeader(http.StatusNotFound)
					w.Write(ErrNotFound.Marshal())
					return
				} // else if errors.Is(err, NoAuthErr) {
				//	w.WriteHeader(http.StatusUnauthorized)
				//	w.Write(ErrNotFound.Marshal())
				//	return
				//}

				err = err.(*AppError)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(appErr.Marshal())
			}

			w.WriteHeader(http.StatusTeapot) // 418 ошибка - точно знаем что накосячил код
			w.Write(systemError(err).Marshal())
		}
	}
}
