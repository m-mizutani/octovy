package usecase

func (x *Usecase) DisableInvokeThread() {
	x.disableInvokeThread = true
}

func SetErrorHandler(uc *Usecase, handler func(error)) {
	uc.testErrorHandler = handler
}

func RunScanThread(uc *Usecase) error {
	return uc.runScanThread()
}

func CloseScanQueue(uc *Usecase) {
	close(uc.scanQueue)
}

type PostGitHubCommentInput postGitHubCommentInput

func PostGitHubComment(input *PostGitHubCommentInput) error {
	d := postGitHubCommentInput(*input)
	return postGitHubComment(&d)
}
