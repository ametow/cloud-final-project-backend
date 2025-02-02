.PHONY: signup signin clean deploy presignedurl updateimageurl

updateimageurl:
	GOOS=linux GOARCH=amd64 go build -o updateimageurl/bootstrap updateimageurl/main.go
	cd updateimageurl && zip updateurl.zip bootstrap

presignedurl:
	GOOS=linux GOARCH=amd64 go build -o presignedurl/bootstrap presignedurl/main.go
	cd presignedurl && zip presign.zip bootstrap

signin:
	GOOS=linux GOARCH=amd64 go build -o signin/bootstrap signin/main.go
	cd signin && zip signin.zip bootstrap

signup:
	GOOS=linux GOARCH=amd64 go build -o signup/bootstrap signup/main.go
	cd signup && zip signup.zip bootstrap

clean:
	rm -f */bootstrap */*.zip