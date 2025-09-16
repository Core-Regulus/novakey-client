package workspaces

import (
	signedrequest "novakeyauth/internal/signedRequest"
	"github.com/google/uuid"
)

type SetWorkspaceRequest struct {	
  Id       uuid.UUID 		 `json:"id,omitempty"`  
	Name     string  			 `json:"name"`
	User 		 signedrequest.SignedRequest `json:"user"`
}

type DeleteWorkspaceRequest struct {
	Id  		uuid.UUID      `json:"id"`
	User 		signedrequest.SignedRequest `json:"user"`
}

type SetWorkspaceResponse struct {
	Id 				uuid.UUID   `json:"id"`
	Name 			string   		`json:"name"`  
	Password 	string   		`json:"password"`
  signedrequest.ErrorResponse
}

type DeleteWorkspaceResponse struct {
	Id uuid.UUID      `json:"id"`	
  signedrequest.ErrorResponse	
}