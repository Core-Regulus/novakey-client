package users

import (
	`novakeyauth/internal/signedRequest`
	"github.com/google/uuid"
)


type SetUserRequest struct {
	signedrequest.SignedRequest
  Id        string  `json:"id,omitempty"`
  Email     string  `json:"email"`
}

type DeleteUserRequest struct {
	signedrequest.SignedRequest
  Id  uuid.UUID      `json:"id"`
  Password  string  `json:"password,omitempty"` 
}

type SetUserResponse struct {
	Id uuid.UUID      `json:"id"`
	Username string   `json:"username"`
  Password string   `json:"password"`
  Status int        `json:"status"`
	Error string	 	  `json:"error,omitempty"`
	Code string	 			`json:"code,omitempty"`
	ErrorDescription string `json:"errorDescription,omitempty"`
}

type DeleteUserResponse struct {
	Id uuid.UUID      `json:"id"`	
  Status int        `json:"status"`
	Error string	 	  `json:"error,omitempty"`
	Code string	 			`json:"code,omitempty"`
	ErrorDescription string `json:"errorDescription,omitempty"`
}