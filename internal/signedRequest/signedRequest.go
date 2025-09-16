package signedrequest

import (
	"fmt"

	"github.com/google/uuid"
)

type SignedRequest struct {
  Id        uuid.UUID `json:"id,omitempty"`
  Username  string `json:"username"`
  Signature string `json:"signature"`
  Message   string `json:"message"`
  PublicKey string `json:"publicKey,omitempty"`
  Timestamp int64  `json:"timestamp"`  
  Password  string  `json:"password,omitempty"`
}

type ErrorResponse struct {
  Status  int             `json:"status"`
  Error   string	 	      `json:"error,omitempty"`
  Code    string	 			  `json:"code,omitempty"`
  ErrorDescription string `json:"errorDescription,omitempty"`
}

func FormatErrorResponse(err ErrorResponse) string {
  return fmt.Sprintf("status %d, error %s, code %s, description %s", err.Status, err.Error, err.Code, err.ErrorDescription)
}