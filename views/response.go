package views

import (
	"errors"
	"net/http"

	"github.com/<%= organization %>/<%= repo %>/session"
	"gorm.io/gorm"
)

type ResponseView struct {
	Data  interface{} `json:"data,omitempty"`
	Error error       `json:"error,omitempty"`
	Prev  string      `json:"prev,omitempty"`
	Next  string      `json:"next,omitempty"`
}

func RenderDataResponse(w http.ResponseWriter, r *http.Request, view interface{}) {
	session.Render(r.Context()).JSON(w, http.StatusOK, ResponseView{Data: view})
}

func RenderErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	sessionError, ok := err.(session.Error)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		sessionError = session.ValidationError("record not found")
	} else if !ok {
		sessionError = session.ServerError(err)
	}
	if sessionError.Code == 10001 {
		sessionError.Code = 500
	}
	session.Render(r.Context()).JSON(w, sessionError.Status, ResponseView{Error: sessionError})
}

func RenderBlankResponse(w http.ResponseWriter, r *http.Request) {
	session.Render(r.Context()).JSON(w, http.StatusOK, ResponseView{})
}
