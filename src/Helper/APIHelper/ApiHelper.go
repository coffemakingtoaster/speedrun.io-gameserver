package ApiHelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	ErrorHelper "gameserver.speedrun.io/Helper/Errorhelper"
	ObjectStructures "gameserver.speedrun.io/Helper/Objecthelper"
)

var apiUrl = "https://api.speedrun.io"

type LobbyReport struct {
	ID             string `json:"lobbyCode"`
	MapCode        string `json:"mapSlug"`
	LobbyName      string `json:"Name"`
	IP             string `json:"ip"`
	Region         string `json:"region"`
	MaxPlayerCount int    `json:"maxPlayerCount"`
	PlayerCount    int    `json:"playerCount"`
}

func setApiURL(url string) {
	apiUrl = url
}

//Generate report struct for player count change
func ReportClientChange(playerCount int, lobby ObjectStructures.LobbyData) LobbyReport {
	return ReportLobbyChange(LobbyReport{PlayerCount: playerCount, ID: lobby.ID, MapCode: lobby.MapCode, IP: getIP(), MaxPlayerCount: 69})
}

//Generate report struct for mapchange
func ReportMapChange(lobby ObjectStructures.LobbyData) LobbyReport {
	return ReportLobbyChange(LobbyReport{MapCode: lobby.MapCode, ID: lobby.ID, IP: getIP(), MaxPlayerCount: 69})
}

//Report changed data to Masterserver
func ReportLobbyChange(data LobbyReport) LobbyReport {
	ip := getIP()
	if ip == "" {
		ErrorHelper.OutputToConsole("Error", "No valid local IP found")
	}

	requestData, err := json.Marshal(data)
	fmt.Println(string(requestData))
	if err != nil {
		ErrorHelper.OutputToConsole("Error", err.Error())
	}

	req, err := http.NewRequest("PATCH", apiUrl+"/v1/lobbies/"+data.ID, bytes.NewBuffer(requestData))
	if err != nil {
		ErrorHelper.OutputToConsole("Error", err.Error())
	}

	//send request to API
	doRequest(req)
	return data
}

//Report the creation of a lobby alongside it´s metadata
func ReportLobby(lobby ObjectStructures.LobbyData) LobbyReport {

	ip := getIP()
	if ip == "" {
		ErrorHelper.OutputToConsole("Error", "No valid local IP found")
	}
	data := LobbyReport{
		ID:             lobby.ID,
		MapCode:        lobby.MapCode,
		LobbyName:      lobby.LobbyName,
		IP:             ip,
		Region:         "EUW",
		MaxPlayerCount: 69,
		PlayerCount:    0,
	}

	requestData, err := json.Marshal(data)
	fmt.Println(string(requestData))
	if err != nil {
		ErrorHelper.OutputToConsole("Error", err.Error())
	}

	req, err := http.NewRequest("POST", apiUrl+"/v1/lobbies", bytes.NewBuffer(requestData))
	if err != nil {
		ErrorHelper.OutputToConsole("Error", err.Error())
	}

	doRequest(req)

	return data
}

//Send Delete request to server to delete lobby from masterserver
func CloseLobby(lobby ObjectStructures.LobbyData) string {
	secret, err := ioutil.ReadFile("./cert/apiSecret.txt")
	if err != nil {
		fmt.Println(err)
	}
	reqUrl := apiUrl + "/v1/lobbies/" + lobby.ID
	req, err := http.NewRequest("DELETE", reqUrl, nil)
	if err != nil {
		ErrorHelper.OutputToConsole("Error", err.Error())
	}
	req.Header.Add("Authorization", "Basic "+strings.TrimSuffix(string(secret), "\n"))
	client := &http.Client{}
	if apiUrl == "" {
		return reqUrl
	}
	resp, err := client.Do(req)
	if err != nil {
		ErrorHelper.OutputToConsole("Error", err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
	return reqUrl
}

//sends request to API
func doRequest(req *http.Request) {
	if apiUrl == "" {
		return
	}
	secret, err := ioutil.ReadFile("./cert/apiSecret.txt")
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+strings.TrimSuffix(string(secret), "\n"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ErrorHelper.OutputToConsole("Error", err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

//get the IP of the current system to send it to the masterserver
func getIP() string {
	//this is a terrible solution. Yet I have not found an efficient way to
	//get the domain of a server. However for the time being this works just fine as
	//there currently is only one gameserver
	//Without a domain the TLS certificate does not work
	return "gameserver.speedrun.io"
	/*
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			return ""
		}
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return localAddr.IP.String()
	*/
}
