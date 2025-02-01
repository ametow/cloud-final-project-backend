.PHONY: signup signin clean deploy presignedurl updateimageurl

updateimageurl:
	GOOS=linux GOARCH=amd64 go build -o updateimageurl/bootstrap updateimageurl/main.go
	cd updateimageurl && zip deployment.zip bootstrap

presignedurl:
	GOOS=linux GOARCH=amd64 go build -o presignedurl/bootstrap presignedurl/main.go
	cd presignedurl && zip deployment.zip bootstrap

signin:
	GOOS=linux GOARCH=amd64 go build -o signin/bootstrap signin/main.go
	cd signin && zip deployment.zip bootstrap

signup:
	GOOS=linux GOARCH=amd64 go build -o signup/bootstrap signup/main.go
	cd signup && zip deployment.zip bootstrap

clean:
	rm -f */bootstrap */deployment.zip

deploy-signup:
	aws lambda update-function-code --function-name YourLambdaFunctionName --zip-file fileb://deployment.zip
