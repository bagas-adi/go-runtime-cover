	dumper, err := coverage.NewDumper()
	if err != nil {
		dumper.PrintLog(err)
	}
	dumper.WatchForFileAndDumpCoverage("./generate_coverage")