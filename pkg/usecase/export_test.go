package usecase

func SetErrorHandler(uc *Usecase, handler func(error)) {
	uc.testErrorHandler = handler
}

func CloseScanQueue(uc *Usecase) {
	close(uc.scanQueue)
}

type PostGitHubCommentInput postGitHubCommentInput

func PostGitHubComment(input *PostGitHubCommentInput) error {
	d := postGitHubCommentInput(*input)
	return postGitHubComment(&d)
}
