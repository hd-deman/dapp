	unrecognized       parserState = "unrecognized"
	diffBegin          parserState = "diffBegin"
	diffBody           parserState = "diffBody"
	newFileDiff        parserState = "newFileDiff"
	deleteFileDiff     parserState = "deleteFileDiff"
	modifyFileDiff     parserState = "modifyFileDiff"
	modifyFileModeDiff parserState = "modifyFileModeDiff"
	ignoreDiff         parserState = "ignoreDiff"
		if strings.HasPrefix(line, "deleted file mode ") {
		if strings.HasPrefix(line, "new file mode ") {
		if strings.HasPrefix(line, "old mode ") {
			return p.handleModifyFileModeDiff(line)
		}
		return fmt.Errorf("unexpected diff line in state `%s`: %#v", p.state, line)
	case modifyFileModeDiff:
		if strings.HasPrefix(line, "new mode ") {
			p.state = unrecognized
			return p.writeOutLine(line)
		}
		return fmt.Errorf("unexpected diff line in state `%s`: %#v", p.state, line)

func (p *diffParser) handleModifyFileModeDiff(line string) error {
	p.state = modifyFileModeDiff
	return p.writeOutLine(line)
}
