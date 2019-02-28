package examples

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

// Example 10-3. Using the command line to retrieve block 277,316
func Block277316(t *testing.T) []byte {
	// hash of block of height 277316
	blockHash := "0000000000000001b6b9a13b095e96db41c4a928b97ef2d944a9b31b2cc7bdc4"

	resp, err := http.Get("https://blockchain.info/rawblock/" + blockHash)
	if nil != err {
		t.Fatal(err)
		return nil
	}
	defer resp.Body.Close()

	block, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		t.Fatal(err)
		return nil
	}

	var out bytes.Buffer
	if err := json.Indent(&out, block, "", "  "); nil != err {
		t.Fatal(err)
	}

	//ioutil.WriteFile("block277316.json", out.Bytes(), 0644)

	return out.Bytes()
}

//func ExampleBlock277316() {
func TestBlock277316(t *testing.T) {
	got := Block277316(t)

	fd, err := os.Open("block277316.json")
	if nil != err {
		t.Fatal(err)
		return
	}
	defer fd.Close()

	expect, err := ioutil.ReadAll(fd)
	if nil != err {
		t.Fatal(err)
		return
	}

	if !bytes.Equal(got, expect) {
		t.Fatalf("invalid block: got %s, expect %s", got, expect)
	}
}
