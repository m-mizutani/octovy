package usecase_test

// TODO: To be fixed
/*
func setupHandleGitHubEvent(t *testing.T) (interfaces.Usecases, *mockSet) {
	const secretsARN = "arn:aws:secretsmanager:us-east-0:123456789012:secret:tutorials/MyFirstSecret-jiObOV"

	cfg := &model.Config{
		SecretsARN:                 secretsARN,
		ScanRequestQueue:           "https://scanreq.queue.url",
		RulePullReqCommentTriggers: "opened",
	}

	uc := usecase.New(cfg)
	svc := usecase.ExposeService(uc)

	// DB mock
	dbClient := newTestTable(t)
	svc.Infra.NewDB = func(region, tableName string) (interfaces.DBClient, error) {
		return dbClient, nil
	}

	// SQS NOT mocked
	svc.Infra.NewSQS = aws.NewSQS

	// SecretsManager
	newSM, mockSM := aws.NewMockSecretsManagerSet()
	mockSM.OutData[secretsARN] = map[string]string{
		"github_app_private_key": base64.StdEncoding.EncodeToString([]byte("zatsu")),
		"github_app_id":          "123",
	}
	svc.Infra.NewSecretManager = newSM

	// GitHubApp
	newApp, mockApp := githubapp.NewMock()
	svc.Infra.NewGitHubApp = newApp
	mockApp.CreateCheckRunMock = func(repo *model.GitHubRepo, commit string) (int64, error) {
		return 503, nil
	}

	return uc, &mockSet{
		db:        dbClient,
		githubapp: mockApp,
	}
}

func TestHandleGitHubEvent(t *testing.T) {
	ts := time.Now().UTC()
	t.Run("Push event", func(t *testing.T) {
		uc, mock := setupHandleGitHubEvent(t)

		pushEvent := github.PushEvent{
			Ref: github.String("refs/heads/master"),
			Repo: &github.PushEventRepository{
				Name: github.String("blue"),
				Owner: &github.User{
					Name: github.String("five"),
				},
				HTMLURL:       github.String("https://github-enterprise.example.com/blue/five"),
				DefaultBranch: github.String("default"),
			},

			Commits: []github.PushEventCommit{
				{
					ID:        github.String("abcdef123"),
					Timestamp: &github.Timestamp{Time: ts},
				},
				{
					ID:        github.String("beefcafe"),
					Timestamp: &github.Timestamp{Time: ts.Add(time.Minute)},
				},
				{
					ID:        github.String("bbbbbbbb"),
					Timestamp: &github.Timestamp{Time: ts.Add(time.Second)},
				},
			},
			Installation: &github.Installation{
				ID: github.Int64(1234),
			},
		}

		require.NoError(t, uc.HandleGitHubPushEvent(&pushEvent))


		require.Equal(t, 1, len(mock.sqs.Input))
		var req model.ScanRepositoryRequest
		require.NoError(t, json.Unmarshal([]byte(*mock.sqs.Input[0].MessageBody), &req))
		assert.Equal(t, "five", req.Owner)
		assert.Equal(t, "blue", req.RepoName)
		assert.Equal(t, "master", req.Branch)
		assert.Equal(t, int64(1234), req.InstallID)
		assert.Equal(t, "beefcafe", req.CommitID)
		assert.Equal(t, ts.Add(time.Minute).Unix(), req.UpdatedAt)
		assert.False(t, req.IsPullRequest)
		assert.False(t, req.IsTargetBranch)
		require.NotNil(t, req.Feedback)
		assert.Equal(t, int64(503), *req.Feedback.CheckID)

		repo, err := mock.db.FindRepoByFullName("five", "blue")
		require.NoError(t, err)
		require.NotNil(t, repo)
		require.Equal(t, "default", repo.DefaultBranch)
	})

	t.Run("Push event with default branch", func(t *testing.T) {
		uc, mock := setupHandleGitHubEvent(t)

		pushEvent := github.PushEvent{
			Ref: github.String("refs/heads/master"),
			Repo: &github.PushEventRepository{
				Name: github.String("blue"),
				Owner: &github.User{
					Name: github.String("five"),
				},
				HTMLURL:       github.String("https://github-enterprise.example.com/blue/five"),
				DefaultBranch: github.String("master"),
			},

			Commits: []github.PushEventCommit{
				{
					ID:        github.String("abcdef123"),
					Timestamp: &github.Timestamp{Time: ts},
				},
			},
			Installation: &github.Installation{
				ID: github.Int64(1234),
			},
		}

		require.NoError(t, uc.HandleGitHubPushEvent(&pushEvent))
		var req model.ScanRepositoryRequest
		require.NoError(t, json.Unmarshal([]byte(*mock.sqs.Input[0].MessageBody), &req))
		assert.True(t, req.IsTargetBranch)
	})

	t.Run("Pull request event", func(t *testing.T) {
		uc, mock := setupHandleGitHubEvent(t)
		pullReqEvent := makePullRequestEvent(&ts)
		require.NoError(t, uc.HandleGitHubPullReqEvent(pullReqEvent))

		var req model.ScanRepositoryRequest
		require.NoError(t, json.Unmarshal([]byte(*mock.sqs.Input[0].MessageBody), &req))

		assert.Equal(t, "five", req.Owner)
		assert.Equal(t, "blue", req.RepoName)
		assert.Equal(t, "ao:1", req.Branch)
		assert.Equal(t, int64(1234), req.InstallID)
		assert.Equal(t, "xyz", req.CommitID)
		assert.Equal(t, ts.Unix(), req.UpdatedAt)
		assert.True(t, req.IsPullRequest)
		assert.False(t, req.IsTargetBranch)

		assert.NotNil(t, req.Feedback)
		require.NotNil(t, req.Feedback.CheckID)
		assert.Equal(t, int64(503), *req.Feedback.CheckID)
		require.NotNil(t, req.Feedback.PullReqID)
		assert.Equal(t, 875, *req.Feedback.PullReqID)

		repo, err := mock.db.FindRepoByFullName("five", "blue")
		require.NoError(t, err)
		require.NotNil(t, repo)
		assert.Equal(t, "five", repo.Owner)
		assert.Equal(t, "blue", repo.RepoName)
	})

	t.Run("Pull request event of synchronize", func(t *testing.T) {
		uc, mock := setupHandleGitHubEvent(t)
		pullReqEvent := makePullRequestEvent(&ts)
		pullReqEvent.Action = github.String("synchronize")
		require.NoError(t, uc.HandleGitHubPullReqEvent(pullReqEvent))

		require.Equal(t, 1, len(mock.sqs.Input))

		t.Run("Not feedback to PR if action is sysynchronize", func(t *testing.T) {
			var req model.ScanRepositoryRequest
			require.NoError(t, json.Unmarshal([]byte(*mock.sqs.Input[0].MessageBody), &req))
			assert.Nil(t, req.Feedback.PullReqID)
		})
	})

	t.Run("Ignore pull request event not opened or sync", func(t *testing.T) {
		uc, mock := setupHandleGitHubEvent(t)
		pullReqEvent := makePullRequestEvent(&ts)
		pullReqEvent.Action = github.String("ready_for_review")
		require.NoError(t, uc.HandleGitHubPullReqEvent(pullReqEvent))

		require.Equal(t, 0, len(mock.sqs.Input))
	})
}

func makePullRequestEvent(ts *time.Time) *github.PullRequestEvent {
	return &github.PullRequestEvent{
		Action: github.String("opened"),
		Repo: &github.Repository{
			Name: github.String("blue"),
			Owner: &github.User{
				Login: github.String("five"),
			},
			HTMLURL:       github.String("https://github-enterprise.example.com/blue/five"),
			DefaultBranch: github.String("default"),
		},
		PullRequest: &github.PullRequest{
			Head: &github.PullRequestBranch{
				SHA:   github.String("xyz"),
				Label: github.String("ao:1"),
			},
			Base: &github.PullRequestBranch{
				Ref: github.String("master"),
			},
			CreatedAt: ts,
			Number:    github.Int(875),
		},
		Installation: &github.Installation{
			ID: github.Int64(1234),
		},
	}
}
*/
