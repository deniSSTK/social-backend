package errors

import "errors"

var ImgBBUploadingError = errors.New("[IMG_BB] failed uploading image")

var EnvironmentVariableNotSet = errors.New("[ENV] environment variable not set: ")
