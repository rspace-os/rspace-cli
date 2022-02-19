Notes for making  new release

- Ensure tests pass
- update CHANGELOG.md
- Regenerate markdown docs: cd docs && go run generateDocs.go
- run ./build.sh VERSION to generate binaries.
- commit changes and push
- tag and push
- Make new Github release from the tag, uploading binaries from dist/ .
