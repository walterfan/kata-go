package prompt

func BuildRefactorPrompt(code string) string {
	return `You are a senior Go engineer.
		Refactor the following Go code to be more idiomatic and readable.
		Output only the updated code.\n\nCode:\n` + code
}
