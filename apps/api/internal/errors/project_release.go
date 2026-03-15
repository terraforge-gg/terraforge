package errors

import "errors"

var (
	ErrProjectReleaseNotFound                         = errors.New("project release not found")
	ErrProjectReleaseDependencyNotFound               = errors.New("project release dependency not found")
	ErrCircularProjectReleaseDependency               = errors.New("circular project release dependency")
	ErrDuplicateProjectReleaseDependency              = errors.New("duplicate project release dependencies")
	ErrProjectReleaseNumberAlreadyExists              = errors.New("project release number already exists")
	ErrProjectReleaseDependencyMinVersionDoesNotExist = errors.New("project release min version number does not exist")
	ErrProjectReleaseInvalidFileSize                  = errors.New("invalid project release file size")
	ErrProjectReleaseFailedToParseFileUrl             = errors.New("file to parse uploaded file url")
	ErrProjectReleaseUploadedFileNotFound             = errors.New("uploaded file not found")
)
