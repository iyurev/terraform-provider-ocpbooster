test-release:
	goreleaser --snapshot --skip-publish --rm-dist

release:
	goreleaser release  --rm-dist
	