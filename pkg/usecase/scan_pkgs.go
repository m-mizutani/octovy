package usecase

/*
func crawlPackages(req *model.ScanRepositoryRequest, clients *scanClients) ([]*ent.PackageRecord, error) {
	tmp, err := ioutil.TempFile("", "*.zip")
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	defer func() {
		if err := os.Remove(tmp.Name()); err != nil {
			logger.Error().Interface("filename", tmp.Name()).Msg("Failed to remove zip file")
		}
	}()

	if err := clients.GitHubApp.GetCodeZip(&req.GitHubRepo, req.CommitID, tmp); err != nil {
		return nil, err
	}

	zipFile, err := zip.OpenReader(tmp.Name())
	if err != nil {
		return nil, goerr.Wrap(err).With("file", tmp.Name())
	}
	defer func() {
		if err := zipFile.Close(); err != nil {
			logger.Error().Interface("zip", zipFile).Err(err).Msg("Failed to close zip file")
		}
	}()

	var newPkgs []*ent.PackageRecord

	for _, f := range zipFile.File {
		psr, ok := parserMap[filepath.Base(f.Name)]
		if !ok {
			continue
		}

		fd, err := f.Open()
		if err != nil {
			return nil, goerr.Wrap(err)
		}
		defer fd.Close()

		pkgs, err := psr.Parse(fd)
		if err != nil {
			return nil, goerr.Wrap(err)
		}

		parsed := make([]*ent.PackageRecord, len(pkgs))
		for i := range pkgs {
			pkg := &ent.PackageRecord{
				Source:  stepDownDirectory(f.Name),
				Type:    psr.PkgType,
				Name:    pkgs[i].Name,
				Version: pkgs[i].Version,
			}
			parsed[i] = pkg
		}

		newPkgs = append(newPkgs, parsed...)
	}

	return newPkgs, nil
}
*/
