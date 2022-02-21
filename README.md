# go-smile :)

A Golang implementation of the [Jackson-Smile](https://github.com/FasterXML/smile-format-specification) binary JSON data format.

## What's supported

* Decoding is fully supported minus:
  * Big Decimals

* Encoding is not yet supported

## Developing

Build, format and run all tests with:

```
make
```

## Tests

Test cases can be found under `test/testdata`. Each test case has a `.smile` and equivalent `.json` file. 

Any new pairs of files added to this directory will be automatically included in the test suite. 

## Usage

```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/zencoder/go-smile/smile"
)

func main() {
	var smileFile = "test/testdata/unicode.smile"

	b, err := ioutil.ReadFile(smileFile)
	if err != nil {
		log.Fatal("Error reading Smile file:", smileFile)
	}

	j, err := smile.DecodeToJSON(b)
	if err != nil {
		log.Fatal("Error decoding Smile file:", smileFile)
	}

	fmt.Println(j)
}
```

```shell script
$ go run main.go 
{"child":"Niño","child-jp":"子供","chilllllllllllllllllllllllld":"Niñññññññññññññññññññññññño"}
```

The following convenience method is also provided for decoding to a Go representation of the JSON:
```go
obj, err := smile.DecodeToObject(b)
```
