package colorx

func RedPrint(s string) {
	Print(WordRed, s)
}

func GreenPrint(s string) {
	Print(WordGreen, s)
}

func RedSprint(s string) string {
	return Sprint(WordRed, s)
}

func GreenSprint(s string) string {
	return Sprint(WordGreen, s)
}

func RedPrintf(s string, a ...interface{}) {
	Printf(WordRed, s, a...)
}

func GreenPrintf(s string, a ...interface{}) {
	Printf(WordGreen, s, a...)
}

func RedSprintf(s string, a ...interface{}) string {
	return Sprintf(WordRed, s, a...)
}

func GreenSprintf(s string, a ...interface{}) string {
	return Sprintf(WordGreen, s, a...)
}

func RedPrintln(s string) {
	Println(WordRed, s)
}

func GreenPrintln(s string) {
	Println(WordGreen, s)
}

func RedKeywordsSprintf(s, keywords string) string {
	return KeywordsSprintf(WordRed, s, keywords)
}

func GreenKeywordsSprintf(s, keywords string) string {
	return KeywordsSprintf(WordGreen, s, keywords)
}
