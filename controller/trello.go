package controller

import (
	// "PMSFreelancer/config"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/adlio/trello"
	"github.com/gin-gonic/gin"
	"github.com/stevejo12/PMSFreelancer/config"
)

// how to open a trello board
// https://www.trello.com/b/{id}

// user key
var key = "d8370c65b1067ab8f964c7102544080f"

func manageTrelloBoard(title string, trelloToken string, freelancerID int) (string, error) {
	client := trello.NewClient(key, trelloToken)

	boardName := title

	board := trello.NewBoard(boardName)

	// trello.Defaults bisa dirubah ama arguments
	// such as description
	err := client.CreateBoard(&board, trello.Defaults())

	if err != nil {
		fmt.Println(err)
		return "", errors.New("Server is unable to create trello board")
	}

	allBoards, err := getUserTrelloBoard(key, trelloToken)

	if err != nil {
		fmt.Println(err)
		return "", errors.New("Server is unable to create trello board")
	}

	var boardIDCreated string

	for i := 0; i < len(allBoards); i++ {
		if allBoards[i].Name == boardName {
			boardIDCreated = allBoards[i].ID
		}
	}

	if boardIDCreated == "" {
		return boardIDCreated, errors.New("Server is unable to find trello board")
	}

	// invite freelancers
	var email string
	err = config.DB.QueryRow("SELECT email FROM login WHERE id=?", freelancerID).Scan(&email)

	if err != nil {
		return boardIDCreated, errors.New("Server is unable to retrieve user email")
	}

	userBoard, err := client.GetBoard(boardIDCreated, trello.Defaults())
	if err != nil {
		return boardIDCreated, errors.New("Server is unable to get the user board with id")
	}

	member := trello.Member{Email: email}
	_, err = userBoard.AddMember(&member, trello.Defaults())

	if err != nil {
		return boardIDCreated, errors.New("Server is unable to add freelancer to the trello board")
	}

	return boardIDCreated, nil
}

// CreateNewBoard godoc
// This function is used to create a new board in trello using user's token
// this target the endpoints of trello public api for creating a board.
func CreateNewBoard(c *gin.Context) {
	var token = "1761f2e9417b6855adf4f50dc00d4721086d6fad6c079a966b673f6c8e927432"

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

	allBoards, err := getUserTrelloBoard(key, token)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server can't get user trello board"})
		return
	}

	var boardIDCreated string

	for i := 0; i < len(allBoards); i++ {
		if allBoards[i].Name == boardName {
			boardIDCreated = allBoards[i].ID
		}
	}

	fmt.Println(boardIDCreated)
}

func getUserTrelloBoard(key string, token string) ([]*trello.Board, error) {
	client := trello.NewClient(key, token)
	//https://api.trello.com/1/members/me/boards?key={yourKey}&token={yourToken}
	args := map[string]string{
		"key":   key,
		"token": token,
	}

	allBoards, err := client.GetMyBoards(args)

	if err != nil {
		return allBoards, errors.New("Server is unable to retrieve user boards")
	}

	return allBoards, nil
}
