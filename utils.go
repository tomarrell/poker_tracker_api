package main

import "net/http"

// FailDecode responds to a request letting the caller
// know the request failed during the decoding of their
// JSON body
func FailDecode(w http.ResponseWriter, err error) {
	RespondServerError(w, []string{
		"Failed to create realm",
		err.Error(),
	})
}

// InvalidArgs sends a response to the client describing the
// required arguments that were missing in the request
func InvalidArgs(w http.ResponseWriter, args []string) {
	RespondBadRequest(
		w,
		map[string]interface{}{
			"message": "Invalid args",
			"args":    args,
		},
	)
}
