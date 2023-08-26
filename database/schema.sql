CREATE TYPE target_class AS ENUM ('os-pkgs', 'lang-pkgs');

create table scans (
    id uuid primary key not null,
    created_at timestamp with time zone not null default now(),
    artifact_name text not null,
    artifact_type text not null,

    repository text,
    branch text,

    page_seq serial
);

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

create table result_packages (
    id uuid primary key not null,
    result_id uuid not null references results(id),
    pkg_id text not null references packages(id)
);

create table vulnerabilities (
    id text primary key not null,

    title text not null,
    description text not null,
    severity text not null,
    cwe_ids text[],
    cvss JSONB,
    reference text[],

    published_at timestamp with time zone,
    last_modified_at timestamp with time zone
);

create table result_vulnerabilities (
    id uuid primary key not null,
    result_id uuid not null references results(id),
    vuln_id text not null references vulnerabilities(id),

    pkg_id text not null references packages(id),
    fixed_version text,
    primary_url text
);
