package usecase

/*
type checkRun struct {
	repo    *model.GitHubRepo
	checkID int64
	app     githubapp.Interface
}

func newCheckRun(app githubapp.Interface) *checkRun {
	return &checkRun{
		app: app,
	}
}

func (x *checkRun) create(repo *model.GitHubRepo, commitID string) error {
	checkID, err := x.app.CreateCheckRun(repo, commitID)
	if err != nil {
		return err
	}
	x.repo = repo
	x.checkID = checkID
	return nil
}

func (x *checkRun) complete(scanID string, changes *pkgChanges, frontendURL string) error {
	logger.Debug().Str("scanID", scanID).
		Str("url", frontendURL).
		Msg("updating check run")

	// TODO: ignore vulnerabilities having status

	opt := &github.UpdateCheckRunOptions{
		Name:       "Octovy: package vulnerability check",
		Status:     github.String("completed"),
		Conclusion: github.String("success"),
		DetailsURL: github.String(frontendURL + "/scan/" + scanID),
		Output: &github.CheckRunOutput{
			Title:   github.String("Scanned new commit"),
			Summary: github.String("OK"),
			Text:    github.String("testing"),
		},
	}

	if err := x.app.UpdateCheckRun(x.repo, x.checkID, opt); err != nil {
		return err
	}

	return nil
}
*/
