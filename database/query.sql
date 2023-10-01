-- name: SaveScan :exec
INSERT INTO scans (
    id,
    artifact_name,
    artifact_type
) VALUES (
    $1, $2, $3
);

-- name: SaveMetaGithubRepository :exec
INSERT INTO meta_github_repository (
    id,
    scan_id,
    owner,
    repo_name,
    branch,
    is_default_branch,
    commit_id,
    base_commit_id,
    pull_request_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
);

-- name: SaveResult :exec
INSERT INTO results (
    id,
    scan_id,
    target,
    target_type,
    class
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: SavePackage :exec
INSERT INTO packages (
    id,
    target_type,
    name,
    version
) VALUES (
    $1, $2, $3, $4
);

-- name: GetPackages :many
SELECT * FROM packages WHERE id = ANY($1::text[]);

-- name: SaveResultPackage :exec
INSERT INTO result_packages (
    id,
    result_id,
    pkg_id
) VALUES (
    $1, $2, $3
);

-- name: SaveVulnerability :exec
INSERT INTO vulnerabilities (
    id,
    title,
    severity,
    published_at,
    last_modified_at,
    data
) VALUES (
    $1, $2, $3, $4, $5, $6
) ON CONFLICT (id)
DO UPDATE SET
    title = $2,
    severity = $3,
    published_at = $4,
    last_modified_at = $5,
    data = $6
WHERE vulnerabilities.last_modified_at < $5;

-- name: GetVulnerability :one
SELECT * FROM vulnerabilities WHERE id = $1;

-- name: GetVulnerabilities :many
SELECT * FROM vulnerabilities WHERE id = ANY($1::text[]);

-- name: SaveResultVulnerability :exec
INSERT INTO result_vulnerabilities (
    id,
    result_id,
    vuln_id,
    pkg_id,
    installed_version,
    fixed_version
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: GetLatestResultsByCommit :many
SELECT results.* FROM results
INNER JOIN (
    SELECT scans.id AS id FROM meta_github_repository
    INNER JOIN scans ON scans.id = results.scan_id
    WHERE meta_github_repository.commit_id = $1
    AND meta_github_repository.owner = $2
    AND meta_github_repository.repo_name = $3
    ORDER BY scans.created_at DESC
    LIMIT 1
) AS latest_scan ON latest_scan.id = results.scan_id;
