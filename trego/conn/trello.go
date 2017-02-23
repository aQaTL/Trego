package conn

import (
	. "github.com/VojtechVitek/go-trello"
	"encoding/json"
	"io/ioutil"
	"log"
)

func Connect() *Member {
	jsonData, err := ioutil.ReadFile("token.json")
	if err != nil {
		log.Panicln(err)
	}

	var credentials struct{ AppKey, Token, Nick string }
	err = json.Unmarshal(jsonData, &credentials)
	if err != nil {
		log.Panicln(err)
	}

	client, err := NewAuthClient(credentials.AppKey, &credentials.Token)
	if err != nil {
		log.Panicln(err)
	}
	usr, err := client.Member(credentials.Nick)
	if err != nil {
		log.Panicln(err)
	}
	return usr
}

func GetLists(board *Board) []List {
	lists, err := board.Lists()
	if err != nil {
		log.Panicln(err)
	}
	return lists
}

func GetBoard(member *Member, boardName string) *Board {
	boards, err := member.Boards()
	if err != nil {
		log.Panicln(err)
	}
	var searchedBoard *Board
	for _, board := range (boards) {
		if board.Name == boardName {
			searchedBoard = &board
		}
	}
	return searchedBoard
}
