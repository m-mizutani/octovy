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
    commit_id,
    base_commit_id,
    pull_request_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
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
    description,
    severity,
    cwe_ids,
    cvss,
    reference,
    published_at,
    last_modified_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
);

-- name: UpdateVulnerability :exec
UPDATE vulnerabilities SET
    title = $2,
    description = $3,
    severity = $4,
    cwe_ids = $5,
    cvss = $6,
    reference = $7,
    published_at = $8,
    last_modified_at = $9
WHERE id = $1 and last_modified_at < $9;

-- name: GetVulnerabilities :many
SELECT * FROM vulnerabilities WHERE id = ANY($1::text[]);

-- name: SaveResultVulnerability :exec
INSERT INTO result_vulnerabilities (
    id,
    result_id,
    vuln_id,
    pkg_id,
    fixed_version,
    primary_url
) VALUES (
    $1, $2, $3, $4, $5, $6
);
