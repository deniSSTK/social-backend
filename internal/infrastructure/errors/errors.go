package errors

import "errors"

var ImgBBUploadingError = errors.New("[IMG_BB] failed uploading image")
var ImgBBUDeletingError = errors.New("[IMG_BB] failed deleting image")

var EnvironmentVariableNotSet = errors.New("[ENV] environment variable not set: ")

var UserIdDoesNotExists = errors.New("[AUTH] user id does not exists")

var ContextUserIdEmpty = errors.New("[CONTEXT] user id is empty")
var ContextParamNotFound = errors.New("[CONTEXT] context param not found: ")
