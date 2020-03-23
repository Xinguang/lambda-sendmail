<pre align="center" type="ascii-art">
 _                 _         _                                 _                 _ _ 
| |               | |       | |                               | |               (_) |
| | __ _ _ __ ___ | |__   __| | __ _ ______ ___  ___ _ __   __| |_ __ ___   __ _ _| |
| |/ _` | '_ ` _ \| '_ \ / _` |/ _` |______/ __|/ _ \ '_ \ / _` | '_ ` _ \ / _` | | |
| | (_| | | | | | | |_) | (_| | (_| |      \__ \  __/ | | | (_| | | | | | | (_| | | |
|_|\__,_|_| |_| |_|_.__/ \__,_|\__,_|      |___/\___|_| |_|\__,_|_| |_| |_|\__,_|_|_|
</pre>
# Lambda sendmail
Send contact form emails using AWS Lambda

Contact forms. They are all around us. Almost every website has one. 

**Lambda sendmail** a solution where you would use Amazon Web Services (AWS) to make your life simple and handle your contact form in a very efficient and cheap way.

## Getting setup

* [Build Project](#build)
* [Get Google reCaptcha Site Key And Secret Key](docs/reCaptcha.md)
* [Create Lambda function](docs/Lambda.md)
* [Create API gateway endpoint](docs/APIgateway.md)
* [Integrate Google reCAPTCHA in your website](#website)
* [Environment variables](#environment)

## Build Project <a name="build" />

Preparing a binary to deploy to AWS Lambda requires that it is compiled for Linux and placed into a .zip file.


## Via docker

``` shell
docker run -it --rm -v $(pwd):/work -w /work lambci/lambda:build-go1.x sam build
cd .aws-sam/build/SendmailFunction
zip main.zip main
```

## For developers on Linux and macOS

``` shell
# Remember to build your handler executable for Linux!
GOOS=linux GOARCH=amd64 go build -o main main.go
zip main.zip main
```

**(docker)Golang Alpine images** doesn't build statically linked binary:
ELF 64-bit LSB executable, x86-64, version 1 (SYSV), dynamically linked, interpreter /lib/ld-musl-x86_64.so.1, stripped

The interpreter doesn't exist in Lambda environment so that's will throws "no such file or directory error".

To solve that use
``` shell
GOOS=linux go build -v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo -o bin/service service/main.go
```


## For developers on Windows

Windows developers may have trouble producing a zip file that marks the binary as executable on Linux. To create a .zip that will work on AWS Lambda, the `build-lambda-zip` tool may be helpful.

Get the tool
``` shell
go.exe get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip
```

Use the tool from your `GOPATH`. If you have a default installation of Go, the tool will be in `%USERPROFILE%\Go\bin`. 

in cmd.exe:
``` bat
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o main main.go
%USERPROFILE%\Go\bin\build-lambda-zip.exe -o main.zip main
```

in Powershell:
``` posh
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -o main main.go
~\Go\Bin\build-lambda-zip.exe -o main.zip main
```

## Integrate Google reCAPTCHA in your website <a name="website" />

To integrate it into your website you need to put it in the client side as well as in Server side. In client HTML page you need to integrate this line before the tag.

```html
<script src="https://www.google.com/recaptcha/api.js?render=put your site key here"></script>
```

Google reCAPTCHA v3 is invisible. You won’t see a captcha form of any sort on your web page. You need to capture the google captcha response in your JavaScript code. Here is a small snippet.


```html
<script>
  grecaptcha.ready(function() {
      grecaptcha.execute('put your site key here', {action: 'homepage'}).then(function(token) {
        // pass the token to the backend script for verification
      });
  });
</script>
```

Let’s look at the [sample code](example.html) to understand it better.
**change the `API_Gateway url` to yours**

```url
"https://API_Gateway_id.execute-api.ap-northeast-1.amazonaws.com/v1/sendmail";
```

## Environment variables <a name="environment" />

| Key  | Value |
| ------------- | ------------- |
| ReCAPTCHA_SECRET  | reCaptcha Secret Key |
| SES_REGION  | AWS Region (E.g.,'us-west-1') |
| MAIL_SUBJECT  | email subject  |
| IS_HTML  | true or false  |
| MAIL_FROM  | from(E.g., recipient <info@site.com> or info@site.com)   |
| TO_ADDRESSES  | recipients (E.g., recipient <info@site.com>,recipient <info@other.com>) |
| CC_ADDRESSES  | recipients (E.g., recipient <info@site.com>,recipient <info@other.com>) |
| BCC_ADDRESSES  | recipients (E.g., recipient <info@site.com>,recipient <info@other.com>) |
| SMTP_USER  | smtp user name |
| SMTP_PASSWORD  | smtp password  |
| SMTP_HOST  | smtp host(E.g., example.com)  |
| SMTP_PORT  | smtp port(E.g., 587)  |

### Dev environment

```sh
sam init --runtime

sam local invoke SendmailFunction --event test-event.json --env-vars .env
# or
sam local start-api

curl -i -X POST \
   -H "content-type:application/json" \
   -d \
'{
  "name":"にほんご",
  "email":"test@test.com",
  "message":"にほんご",
  "token":"tokenxxxxx"
}' \
 'http://127.0.0.1:3000/sendmail'
```
