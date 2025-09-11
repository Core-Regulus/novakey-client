package workspaces

import (
	signedrequest "novakeyauth/internal/signedRequest"
	"github.com/google/uuid"
)

type SetWorkspaceRequest struct {
	signedrequest.SignedRequest
  Id        string  `json:"id,omitempty"`
  Email     string  `json:"email"`
	Name     string   `json:"name"`
}

type DeleteWorkspaceRequest struct {
	signedrequest.SignedRequest
  Id  uuid.UUID      `json:"id"`
  Password  string  `json:"password,omitempty"` 
}

type SetWorkspaceResponse struct {
	Id uuid.UUID      `json:"id"`
	Name string   		`json:"name"`  
	Password string   `json:"password"`
  Status int        `json:"status"`
	Error string	 	  `json:"error,omitempty"`
	Code string	 			`json:"code,omitempty"`
	ErrorDescription string `json:"errorDescription,omitempty"`
}

type DeleteWorkspaceResponse struct {
	Id uuid.UUID      `json:"id"`	
  Status int        `json:"status"`
	Error string	 	  `json:"error,omitempty"`
	Code string	 			`json:"code,omitempty"`
	ErrorDescription string `json:"errorDescription,omitempty"`
}