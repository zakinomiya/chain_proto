package repository

type RepoError int

const (
	ErrUnknown RepoError = iota
	ErrNotFound
)

func (re RepoError) Error() string {
	return re.String()
}

func (re RepoError) String() string {
	switch re {
	case ErrNotFound:
		return "Resource Not Found"
	default:
		return "Unknown Error"
	}
}
