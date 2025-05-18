package prompt

func BuildReviewPrompt(filename, code string) string {
	return "You are a senior Go engineer. Review the following Go source code and suggest improvements, bug fixes, and style corrections.\n\n" +
		"Filename: " + filename + "\n\nCode:\n" + code
}
