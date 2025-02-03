.PHONY: signup signin clean deploy presignedurl updateimageurl

updateimageurl:
	GOOS=linux GOARCH=amd64 go build -o updateimageurl/bootstrap updateimageurl/main.go
	cd updateimageurl && zip updateurl.zip bootstrap && mv updateurl.zip ..

presignedurl:
	GOOS=linux GOARCH=amd64 go build -o presignedurl/bootstrap presignedurl/main.go
	cd presignedurl && zip presign.zip bootstrap && mv presign.zip ..

signin:
	GOOS=linux GOARCH=amd64 go build -o signin/bootstrap signin/main.go
	cd signin && zip signin.zip bootstrap && mv signin.zip ..

signup:
	GOOS=linux GOARCH=amd64 go build -o signup/bootstrap signup/main.go
	cd signup && zip signup.zip bootstrap && mv signup.zip ..

clean:
	rm -f */bootstrap */*.zip