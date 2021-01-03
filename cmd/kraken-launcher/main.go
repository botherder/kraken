package main

import (
	"fmt"

	"github.com/botherder/go-savetime/hashes"
	"github.com/botherder/kraken/storage"
	"github.com/go-resty/resty/v2"
)

// Check-in to the API server.
func apiVersionCheck() (string, error) {
	type Response struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Message string `json:"msg"`
		URL     string `json:"url"`
	}

	hash, _ := hashes.FileSHA1(storage.StorageExe)

	client := resty.New()
	response, err := client.R().
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetBody(fmt.Sprintf(`{"sha1": "%s"}`, hash)).
		SetResult(&Response{}).
		Post("https://%s/api/versioncheck/")

	if err != nil {
		return "", fmt.Errorf("Unable to check version with REST API: %s", err.Error())
	}

	data := response.Result().(*Response)

	if data.Code == "OK_DOWNLOAD" {
		return data.URL, nil
	}

	// This should mean there are no updates.
	return "", nil
}

func download(url string) error {
	client := resty.New()
	_, err := client.R().
		SetOutput(storage.StorageExe).
		Get(url)

	if err != nil {
		return fmt.Errorf("Unable to download URL %s: %s", url, err.Error())
	}

	return nil
}

func main() {
	url, err := apiVersionCheck()
	if err != nil {
		fmt.Println("[!] ERROR: ", err.Error())
	}

	if url != "" {
		fmt.Println("[+] Instructed to download new agent from: ", url)

		err = download(url)
		if err != nil {
			fmt.Println("[!] ERROR: ", err.Error())
		}
	} else {
		fmt.Println("[-] Nothing new to download.")
	}

	err = launchAgent()
	if err != nil {
		fmt.Println("[*] Agent has been launched! Exiting launcher...")
	}

	// TODO: can't find how to detach the process, so fuck it, we just wait...
	// cmd.Wait()
}
