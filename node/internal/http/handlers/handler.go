package handlers

import "net/http"

// if request is ok - [INFO] log time, id, method, uri, status code, duration
// on internal error - [ERROR] log time, id, method, uri, status code, duration, error


// custom handlers which allows logging errors
type Handler = func(w http.ResponseWriter, r *http.Request) error

func WithError = 

