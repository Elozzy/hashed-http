package main
import "testing"

func TestFetchUrl(t *testing.T) {
	//test to check if url is avaliable
	url := parseRawUrl("http://google.com") 
		if url != "http://google.com"{
			t.Error("url must be a string")
		}
	}

	
