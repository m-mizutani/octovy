package usecase

func CloseScanQueue(uc *Usecase) {
	close(uc.scanQueue)
}

type PostGitHubCommentInput postGitHubCommentInput

func PostGitHubComment(input *PostGitHubCommentInput) error {
	d := postGitHubCommentInput(*input)
	return postGitHubComment(&d)
}
