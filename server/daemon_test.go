package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/triasteam/go-streamnet/store"

	"github.com/triasteam/go-streamnet/types"
)

/*func TestStart(t *testing.T) {
	Start()

}*/

func TestStartAndStop(t *testing.T) {
	// open database
	db := store.Storage{}
	db.Init("/tmp/gorocksdb_http_test")

	// start http server
	go Start(&db)
	time.Sleep(1 * time.Second)
	defer Stop()

	/*resp, err := http.Get("http://127.0.0.1:14700/")
	if err != nil {
		t.Fatalf("Failed to get: %v", err)
	}
	defer resp.Body.Close()*/

	data := types.StoreData{
		Attester: "192.168.130.1",
		Attestee: "192.168.130.2",
		Score:    "1",
	}

	j, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("json failed: %v\n", err)
	}

	resp1, err := http.Post("http://127.0.0.1:14700/save", "application/json", strings.NewReader(string(j)))
	if err != nil {
		t.Fatalf("Save message failed: %v\n", err)
	}
	defer resp1.Body.Close()

	body, err := ioutil.ReadAll(resp1.Body)
	t.Logf("body: %s\n", string(body))

}

/*
func TestSaveHandle(t *testing.T) {
	time.Sleep(3 * time.Second)
	go Start()
	time.Sleep(1 * time.Second)
	defer Stop()

	t.Log("Hello...\n")
	resp, err := http.PostForm("http://127.0.0.1:14700/save",
		url.Values{"Attester": {"192.168.130.1"}, "Attestee": {"192.168.130.2"}, "Score": {"1"}})
	if err != nil {
		t.Fatalf("Save message failed: %v\n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	t.Logf("body: %s\n", string(body))
}*/

/*func TestGetHandle(t *testing.T) {

}*/
