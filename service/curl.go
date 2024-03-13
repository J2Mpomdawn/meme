package service

import (
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/rdegges/go-ipify"
)

// allow ip
func SSHPermission() {
	//get ip
	ip, err := ipify.GetIp()

	if err != nil {
		LogPrint("GetIp", "red", err)
	}

	//set api parameters
	values := url.Values{}
	values.Add("account", os.Getenv("SSH_User"))
	values.Add("server_name", os.Getenv("SSH_Host"))
	values.Add("api_secret_key", os.Getenv("API_Key"))
	values.Add("param[addr]", ip)

	//api request
	req, err := http.NewRequest("POST",
		"https://api.xrea.com/v1/tool/ssh_ip_allow",
		strings.NewReader(values.Encode()))

	if err != nil {
		LogPrint("red", "NewRequest", err)
	}

	//set api request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	//call api
	res, err := client.Do(req)

	if err != nil {
		LogPrint("red", "client-Do", err)
	}

	defer res.Body.Close()
}
