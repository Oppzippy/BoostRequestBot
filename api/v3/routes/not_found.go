package routes

import (
	"net/http"

	"github.com/oppzippy/BoostRequestBot/api/responder"
)

func NotFoundHandler(rw http.ResponseWriter, r *http.Request) {
	responder.RespondError(rw, http.StatusNotFound, "The requested API endpoint does not exist.")
}
