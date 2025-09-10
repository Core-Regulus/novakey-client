package signedrequest



type SignedRequest struct {
  Username  string `json:"username"`
  Signature string `json:"signature"`
  Message   string `json:"message"`
  PublicKey string `json:"publicKey,omitempty"`
  Timestamp int64  `json:"timestamp"`  
}