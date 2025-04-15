package errors

type UploadError struct {
	HasError       bool
	DomainError    error
	TechnicalError error
}

func NewError(domainErr error, technicalError error) UploadError {
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

func NilError() UploadError {
	return UploadError{false, nil, nil}
}
