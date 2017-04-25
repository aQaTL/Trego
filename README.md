# Trego

Simple Trello console client (but with UI). 

### Building

I haven't published any binaries yet, so if you want to try Trego on your computer, you'll to 
perform following steps:

* Download, install and configure environment variables [golang](https://golang.org/dl/)
* Download Trego via `go get github.com/aqatl/trego` (needs [git](https://git-scm.com/) 
installed)
* Create [token.json](#token.json) file

Trego should be in the `$GOPATH/bin` directory. 

### token.json

Getting all of the necessary tokens is still a TODO feature. So, for now you have to get them 
manually and put them into the `token.json` file. First, you will need to get application key 
from [here](https://trello.com/app-key). Then you will have to generate your token. 
Go to this url: `https://trello.com/1/authorize?response_type=token&expiration=never&name=Trego&scope=read,write,account&key=YOURAPPKEY` remember to replace the `YOURAPPKEY` with 
the key you got previously. 

Now, create the `token.json` file and place it in the same directory as the Trego binary. 
Then, fill that file using the template below.

```
{
	"AppKey": "paste your key here",
	"Token": "paste your generated token here"
}
```
