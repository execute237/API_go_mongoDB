package apperror

import (
	"errors"
	"net/http"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func Middleware(h appHandler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var appErr *AppError
		err := h(writer, request)
		if err != nil {
			writer.Header().Set("Content-Type", "application/json")
			if errors.As(err, &appErr) {
				if errors.Is(err, ErrNotFound) {
					writer.WriteHeader(http.StatusNotFound)
					writer.Write(ErrNotFound.Marshal())
					return
				} //else if errors.Is(err, NoAuthErr) {
				//	writer.WriteHeader(http.StatusUnauthorized)
				//	writer.Write(ErrNotFound.Marshal())
				//	return
				//} TODO

				err = err.(*AppError)
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write(appErr.Marshal())
				return
			}

			writer.WriteHeader(http.StatusTeapot)
			writer.Write(systemError(err).Marshal())
		}
	}
}
