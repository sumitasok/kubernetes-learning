# kubernetes-learning
repo for having sample application that helps in learning kubernetes.

## Requirement and Task analysis

Text files only; check the allowed file extension;

Authentication? For now a token based or http basic auth will be provided. User specific data segregation using Authorization cannot be provided in this time frame.

Ideally a database holding file  information like name size md5 value will be useful in improving  performance of request response.

This is a perfect candidate for Stateful sets; and shared volumes in kubernetes. Because my experience in that is limited, I will be running a single pod (no replicas) for the time being.

While implementing a service exposed through http, following AWS s3 API is recommended as industry is considering it to be a standard. But cant do it for the two day assignment.

Use of the remote-update protocol, to transfer only the difference between files content like rsync.
ref: https://www.tecmint.com/sync-new-changed-modified-files-rsync-linux/

#### Server
- Fetch all files in the system.
- Fetch file by name <- return file name, size, md5 sum
- verify API
    - using the filename, size and md5 send to server; check if this file already exists by name or by md5 match.
- Add, validate if the file exists, with same size and md5 check sum
	- Return error if file exist with details about the uniqueness value(same or new file)
	- Store the md5 of the file for easy comparison of files
- Update
	- Verify a change in md5 before doing change
	- Apply the diff patch
	- Or, Complete replacement of file

#### Client
- Add
	- Check if file exist with same md5, then print no add required
	- If different md5, then through error on add, but suggest using upload command.
		
- Update
	- Make update request with md5 and verify if there is any change in the client copy compared to existing server copy.

How to make sure only changed data is send across network like applying a diff patch.

[] setup pre-commit hooks for analysing code quality.


### Install Pre-commit for code analysis

```
brew install pre-commit

# from inside the repo root
pre-commit install

go get -u golang.org/x/lint/golint
go get -u github.com/fzipp/gocyclo/cmd/gocyclo


brew tap alecthomas/homebrew-tap
brew install gometalinter

```

### code structure standard documentations followed from

https://github.com/golang-standards/project-layout

Notes:

- go-vet in pre-commit hook is throwing error. prioritising the task over go vet in pre-commit. go vet will be run manually.
- using `main` branch as this is a test task. a complete git workflow will be followed in case of usual development.
