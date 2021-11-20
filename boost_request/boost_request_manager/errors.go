package boost_request_manager

type BoostRequestSignupError struct {
	message string
}

func (e *BoostRequestSignupError) Error() string {
	return e.message
}

var (
	ErrNoPrivileges error = &BoostRequestSignupError{
		message: "the user has no privileges",
	}
	ErrNotPreferredAdvertiser error = &BoostRequestSignupError{
		message: "the user is not a preferred advertiser",
	}
)
