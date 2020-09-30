package note

type Config struct {
	Repository       string
	BranchWhitelist  []string
	BranchBlacklist  []string
	DatabaseLocation string
	Encrypt          bool
	Port             int
}
