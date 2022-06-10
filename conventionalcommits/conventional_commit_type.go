package conventionalcommits

type ConventionalCommitType int

const (
	Chore ConventionalCommitType = iota
	Fix
	Feature
	BreakingChange
)

func fromString(s string) ConventionalCommitType {
	for _, choreType := range []string{"build", "chore", "ci", "docs", "style", "refector", "perf", "test"} {
		if s == choreType {
			return Chore
		}

		if s == "fix" {
			return Fix
		}

		if s == "feat" {
			return Feature
		}
	}

	panic("string is no conventional commit type")
}
