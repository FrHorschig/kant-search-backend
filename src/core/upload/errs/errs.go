package errs

type UploadError struct {
	HasError       bool
	DomainError    error
	TechnicalError error
}

func New(domainErr error, technicalError error) UploadError {
	hasError := false
	if domainErr != nil || technicalError != nil {
		hasError = true
	}
	return UploadError{
		HasError:       hasError,
		DomainError:    domainErr,
		TechnicalError: technicalError,
	}
}

func Nil() UploadError {
	return UploadError{false, nil, nil}
}
