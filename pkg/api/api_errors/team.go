package api_errors

import "errors"

var (
	ErrorOnlyOneTeamAllowed = errors.New("only one team allowed")
	ErrorTeamNotFound       = errors.New("team not found")
	ErrorUserAlreadyInTeam  = errors.New("user already exists in a team")
)
