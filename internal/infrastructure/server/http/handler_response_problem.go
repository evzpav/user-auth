package http

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/evzpav/user-auth/pkg/errors"
)

//RFC7807Title is an internal type to handle with RFC7807 title
type RFC7807Title string

const (
	//RFC7807TitleInvalidArgument is the RFC7807 title for invalid argument error
	RFC7807TitleInvalidArgument RFC7807Title = "invalid argument"

	//RFC7807TitleNotFound is the RFC7807 title for not found error error
	RFC7807TitleNotFound RFC7807Title = "resource not found"

	//RFC7807TitleRuleNotSatisfied is the RFC7807 title for rules not satisfied error
	RFC7807TitleRuleNotSatisfied RFC7807Title = "internal rule was not satisfied"

	//RFC7807TitleDuplicated is the RFC7807 title for duplicated record error
	RFC7807TitleDuplicated RFC7807Title = "duplicated record"

	//RFC7807TitleBadRequest is the RFC7807 title for bad request error
	RFC7807TitleBadRequest RFC7807Title = "bad request"

	//RFC7807TitleForbidden is the RFC7807 title for forbidden error
	RFC7807TitleForbidden RFC7807Title = "forbidden"

	//RFC7807TitleNotAuthorized is the RFC7807 title for not authorized error
	RFC7807TitleNotAuthorized RFC7807Title = "not authorized"

	//RFC7807TypeBlank is the RFC7807 default type
	RFC7807TypeBlank string = "about:blank"
)

var (
	//ErrInternalServer is the RFC7807 internal server error. It should be used when an unknown error occur
	ErrInternalServer = NewRFC7807Failure(RFC7807TypeBlank, "the server encountered an unexpected condition that prevented it from fulfilling the request", "", "", http.StatusInternalServerError, "", nil)
)

//RFC7807Failure wraps the problem json that we send to client according to https://tools.ietf.org/html/rfc7807
type RFC7807Failure struct {
	// Type contains a URI that identifies the problem type. This URI will,
	// ideally, contain human-readable documentation for the problem when
	// dereferenced
	Type string `json:"type"`

	// The HTTP status code for this occurrence of the problem
	Status int `json:"status,omitempty"`

	Code string `json:"code,omitempty"`

	// Title is a short, human-readable summary of the problem type. This title
	// SHOULD NOT change from occurrence to occurrence of the problem, except for
	// purposes of localization
	Title string `json:"title"`

	// A human-readable explanation specific to this occurrence of the problem.
	Detail string `json:"detail,omitempty"`

	// A URI that identifies the specific occurrence of the problem. This URI
	// may or may not yield further information if dereferenced.
	Instance string `json:"instance,omitempty"`

	//Args storage the key value of some parameter that has unexpected value or type
	Args map[string]interface{} `json:"arguments,omitempty"`
}

//Error is the type of error
func (f *RFC7807Failure) Error() string {
	return f.Detail
}

//NewRFC7807Failure creates a new instance of RFC7807Failure
func NewRFC7807Failure(tp, tittle, detail, instance string, status int, code string, args map[string]interface{}) *RFC7807Failure {
	err := &RFC7807Failure{}
	err.Type = tp
	err.Title = tittle
	err.Status = status
	err.Instance = instance
	err.Detail = detail
	err.Args = args
	err.Code = code
	return err
}

func (h *handler) responseProblem(c *gin.Context, err error) {
	if nf, ok := errors.NotFoundCast(err); ok {
		h.responseProblemBuilder(c, nf, http.StatusNotFound, RFC7807TitleNotFound)
		return
	}

	if ia, ok := errors.InvalidArgumentCast(err); ok {
		h.responseProblemBuilder(c, ia, http.StatusBadRequest, RFC7807TitleInvalidArgument)
		return
	}

	if dub, ok := errors.DuplicatedRecordCast(err); ok {
		h.responseProblemBuilder(c, dub, http.StatusConflict, RFC7807TitleDuplicated)
		return
	}

	if rns, ok := errors.RuleNotSatisfiedCast(err); ok {
		h.responseProblemBuilder(c, rns, http.StatusBadRequest, RFC7807TitleRuleNotSatisfied)
		return
	}

	if rns, ok := errors.NotAuthorizedCast(err); ok {
		h.responseProblemBuilder(c, rns, http.StatusUnauthorized, RFC7807TitleNotAuthorized)
		return
	}

	// if fdn, ok := errors.ForbiddenCast(err); ok {
	// 	h.responseProblemBuilder(c, fdn, http.StatusForbidden, RFC7807TitleForbidden)
	// 	return
	// }

	h.responseProblemBuilder(c, err, http.StatusInternalServerError, "")
}

func (h *handler) responseProblemBuilder(c *gin.Context, err error, status int, title RFC7807Title) {
	describer, ok := errors.DescriberCast(err)
	var rfc *RFC7807Failure
	if ok {
		rfc = &RFC7807Failure{
			Type:   RFC7807TypeBlank,
			Status: status,
			Code:   string(describer.GetCode()),
			Title:  string(title),
			Detail: describer.GetMessage(),
			Args:   describer.GetArgs(),
		}
	} else {
		rfc = ErrInternalServer
	}

	_ = c.Error(err)
	h.responseProblemWriter(c, rfc)
}

func (h *handler) responseProblemWriter(c *gin.Context, rfc *RFC7807Failure) {
	if problemJSON, err := json.Marshal(rfc); err == nil {
		c.Data(rfc.Status, "application/problem+json; charset=utf-8", problemJSON)
		return
	}
	c.AbortWithStatus(http.StatusInternalServerError)
}
