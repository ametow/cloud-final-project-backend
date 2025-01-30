.PHONY: signup clean deploy

signup:
	GOOS=linux GOARCH=amd64 go build -o signup/bootstrap signup/main.go
	cd signup && zip deployment.zip bootstrap

clean:
	rm -f */bootstrap */deployment.zip

deploy-signup:
	aws lambda update-function-code --function-name YourLambdaFunctionName --zip-file fileb://deployment.zip
