package gordon

// A NilPageError is created when a Handler returns a nil Page, without
// also returning a valid error
type NilPageError struct{}

// Error fulfills the error interface
func (NilPageError) Error() string {
	return "nil page returned from Handler"
}
