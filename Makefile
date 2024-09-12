test:
	go test ./pyenv -run "(TestIntegration|TestDependencies)"  -count=1
test-remove:
	go test ./pyenv -run TestRemove  -count=1