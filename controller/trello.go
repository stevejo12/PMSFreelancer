package controller

import (
	"errors"
	"fmt"
	"log"

	"github.com/adlio/trello"
	"github.com/gin-gonic/gin"
)

// how to open a trello board
// https://www.trello.com/b/{id}

// key token pribadi
var key = "d8370c65b1067ab8f964c7102544080f"

var token = "1761f2e9417b6855adf4f50dc00d4721086d6fad6c079a966b673f6c8e927432"

// token pinjeman
// var token = "51ee616bd00d17e7615e2aca0dc0d849211863855567446935478143a96c4115"
var client = trello.NewClient(key, token)

func createTrelloBoard(title string, trelloToken string) (string, error) {
	client := trello.NewClient(key, trelloToken)

	boardName := title

	board := trello.NewBoard(boardName)

	// trello.Defaults bisa dirubah ama arguments
	// such as description
	err := client.CreateBoard(&board, trello.Defaults())

	if err != nil {
		log.Fatal(err)
	}

	allBoards := GetUserTrelloBoard()

	var boardIDCreated string

	for i := 0; i < len(allBoards); i++ {
		if allBoards[i].Name == boardName {
			boardIDCreated = allBoards[i].ID
		}
	}

	if boardIDCreated == "" {
		return boardIDCreated, errors.New("Server is unable to find trello board")
	}

	return boardIDCreated, nil
}

// CreateNewBoard godoc
// This function is used to create a new board in trello using user's token
// this target the endpoints of trello public api for creating a board.
func CreateNewBoard(c *gin.Context) {
	client := trello.NewClient(key, token)

	boardName := "testing with auth"

	board := trello.NewBoard(boardName)

	board.Desc = "testing a description"

	// trello.Defaults bisa dirubah ama arguments
	// such as description
	err := client.CreateBoard(&board, trello.Defaults())

	if err != nil {
		log.Fatal(err)
	}

	allBoards := GetUserTrelloBoard()

	var boardIDCreated string

	for i := 0; i < len(allBoards); i++ {
		if allBoards[i].Name == boardName {
			boardIDCreated = allBoards[i].ID
		}
	}

	fmt.Println(boardIDCreated)
}

// GetUserTrelloBoard => Fetching all the users' boards using token and key
func GetUserTrelloBoard() []*trello.Board {
	client := trello.NewClient(key, token)
	//https://api.trello.com/1/members/me/boards?key={yourKey}&token={yourToken}
	args := map[string]string{
		"key":   key,
		"token": token,
	}

	allBoards, err := client.GetMyBoards(args)

	if err != nil {
		log.Fatal(err)
	}

	return allBoards
}
