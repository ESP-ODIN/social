package service

import "errors"

var ErrPostForbidden = errors.New("only the author can delete this post")
