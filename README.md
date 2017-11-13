# Trego

Simple Trello console client (but with UI). 

### Screenshots

![trego](https://user-images.githubusercontent.com/17130832/27430489-1401f12e-5749-11e7-8dc3-c7749d561c33.png)

![trego2](https://user-images.githubusercontent.com/17130832/27430490-1402c0f4-5749-11e7-9007-140e3af8fed8.png)

### Building

You can find built binaries for GNU/Linux, Windows and MacOS in the 
[releases](https://github.com/aQaTL/Trego/releases) tab.

If you want to built the app by yourself, perform steps listed below: 

* Download and install [golang](https://golang.org/dl/)
* Set the `GOPATH` environment variable (set it to the folder, where you want to store `go` 
projects)
* Download Trego via `go get -u github.com/aqatl/trego` (needs [git](https://git-scm.com/) 
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
