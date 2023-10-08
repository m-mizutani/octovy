CREATE TYPE target_class AS ENUM ('os-pkgs', 'lang-pkgs');

create table scans (
    id uuid primary key not null,
    created_at timestamp with time zone not null default now(),
    artifact_name text not null,
    artifact_type text not null,

    page_seq serial
);

create table meta_github_repository (
    id uuid primary key not null,
    scan_id uuid not null references scans(id),

    owner text not null,
    repo_name text not null,
    commit_id text not null,
    branch text,
    is_default_branch boolean,
    base_commit_id text,
    pull_request_id int,

    page_seq serial
);

CREATE INDEX meta_github_repository_commit ON meta_github_repository (commit_id);

create table results (
    id uuid primary key not null,
    scan_id uuid not null references scans(id),

    target text not null,
    target_type text not null,
    class target_class not null
);

create table packages (
    -- id is hash of target_type, name, version
    id text primary key not null,

    target_type text not null,
    name text not null,
    version text not null
);

create table detected_packages (
    id uuid primary key not null,
    result_id uuid not null references results(id),
    pkg_id text not null references packages(id)
);

create table vulnerabilities (
    id text primary key not null,

    title text not null,
    severity text not null,
    published_at timestamp with time zone,
    last_modified_at timestamp with time zone,

    data JSONB,
    page_seq serial
);

create table detected_vulnerabilities (
    id uuid primary key not null,
    result_id uuid not null references results(id),
    vuln_id text not null references vulnerabilities(id),

    pkg_id text not null references packages(id),
    fixed_version text,
    installed_version text,

    data JSONB
);
