package usecase

import (
	"database/sql"
	"encoding/json"

	ttypes "github.com/aquasecurity/trivy/pkg/types"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/db"
)

type DiffResultKey struct {
	Class      string `json:"class"`
	Target     string `json:"target"`
	TargetType string `json:"type"`
}

type DiffStatus string

const (
	DiffStatusAdd = DiffStatus("add")
	DiffStatusDel = DiffStatus("del")
	DiffStatusMod = DiffStatus("mod")
)

type DiffResult struct {
	DiffResultKey
	Status DiffStatus                     `json:"status"`
	Add    []ttypes.DetectedVulnerability `json:"add"`
	Del    []ttypes.DetectedVulnerability `json:"del"`
}

func getVulnDiffForGitHubRepo(ctx *model.Context, dbClient *sql.DB, report *ttypes.Report, commit *model.GitHubCommit) ([]DiffResult, error) {
	reportResults := map[DiffResultKey]*ttypes.Result{}
	for i, result := range report.Results {
		key := DiffResultKey{
			Class:      string(result.Class),
			Target:     result.Target,
			TargetType: string(result.Type),
		}
		reportResults[key] = &report.Results[i]
	}

	query := db.New(dbClient)
	resp, err := query.GetLatestResultsByCommit(ctx, db.GetLatestResultsByCommitParams{
		RepoID:   commit.RepoID,
		CommitID: commit.CommitID,
	})
	if err != nil {
		return nil, goerr.Wrap(err, "failed to get latest results by commit").With("commit", commit)
	}

	var diffResults []DiffResult
	for _, r := range resp {
		oldVulns, err := query.GetVulnerabilitiesByResultID(ctx, r.ID)
		if err != nil {
			return nil, goerr.Wrap(err, "failed to get vulnerabilities by result ID").With("resultID", r.ID)
		}
		oldVulnMap := map[string]ttypes.DetectedVulnerability{}
		for _, old := range oldVulns {
			var v ttypes.DetectedVulnerability
			if err := json.Unmarshal(old.Data.RawMessage, &v); err != nil {
				return nil, goerr.Wrap(err, "failed to unmarshal vulnerability").With("resultID", r.ID).With("data", old.Data)
			}
			oldVulnMap[v.VulnerabilityID] = v
		}

		key := DiffResultKey{
			Class:      string(r.Class),
			Target:     r.Target,
			TargetType: r.TargetType,
		}

		result, ok := reportResults[key]
		if !ok {
			var vulnList []ttypes.DetectedVulnerability
			for _, v := range oldVulnMap {
				vulnList = append(vulnList, v)
			}
			diffResults = append(diffResults, DiffResult{
				DiffResultKey: key,
				Status:        DiffStatusDel,
				Del:           vulnList,
			})
			continue
		}

		var add, del []ttypes.DetectedVulnerability
		for _, v := range result.Vulnerabilities {
			if _, ok := oldVulnMap[v.VulnerabilityID]; !ok {
				add = append(add, v)
			}
			delete(oldVulnMap, v.VulnerabilityID)
		}
		for _, v := range oldVulnMap {
			del = append(del, v)
		}

		if len(add) > 0 || len(del) > 0 {
			diffResults = append(diffResults, DiffResult{
				DiffResultKey: key,
				Status:        DiffStatusMod,
				Add:           add,
				Del:           del,
			})
		}

		delete(reportResults, key)
	}

	for key, result := range reportResults {
		diffResults = append(diffResults, DiffResult{
			DiffResultKey: key,
			Status:        DiffStatusAdd,
			Add:           result.Vulnerabilities,
		})
	}

	return diffResults, nil
}
