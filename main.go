package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
)

type printItem struct {
	ID      int
	Project string
	Type    string
	Title   string
	Effort  float64
}

func (i printItem) String() string {
	return fmt.Sprintf("[%s] %s%d: %s / %.0f", i.Project, i.emojiType(), i.ID, i.Title, i.Effort)
}

func (i printItem) emojiType() string {
	switch i.Type {
	case "Epic":
		return "👑"
	case "Product Backlog Item":
		fallthrough
	case "User Story":
		return "📘"
	case "Feature":
		return "🏆"
	case "Impediment":
		fallthrough
	case "Issue":
		return "🚨"
	case "Task":
		return "📋"
	case "Bug":
		return "🐞"
	case "Test Case":
		return "🧪"
	default:
		return i.Type
	}
}

func main() {
	organization := flag.String("o", "", "organization")
	todo := flag.String("todo", "", "todo query id")
	doing := flag.String("doing", "", "doing query id")
	done := flag.String("done", "", "done query id")
	token := flag.String("t", "", "personal access token")
	slack := flag.String("s", "", "slack web hook url")
	name := flag.String("n", "", "name for report")
	flag.Parse()
	if len(*organization) == 0 || len(*todo) == 0 || len(*doing) == 0 ||
		len(*done) == 0 || len(*token) == 0 || len(*name) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	connection := azuredevops.NewPatConnection("https://dev.azure.com/"+*organization, *token)
	ctx := context.Background()
	client, err := workitemtracking.NewClient(ctx, connection)
	if err != nil {
		log.Fatal(err)
	}

	msg := "*" + *name + "'s Weekly Report*\n\n"
	todoMsg, err := subTotal(client, ctx, *todo, "To-Do")
	if err != nil {
		log.Fatal(err)
	}
	msg += todoMsg
	doingMsg, err := subTotal(client, ctx, *doing, "DOING")
	if err != nil {
		log.Fatal(err)
	}
	msg += doingMsg
	doneMsg, err := subTotal(client, ctx, *done, "DONE")
	if err != nil {
		log.Fatal(err)
	}
	msg += doneMsg

	if len(*slack) > 0 {
		if err := postSlack(*slack, msg); err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println(msg)
	}
}

func subTotal(client workitemtracking.Client, ctx context.Context, queryID string, title string) (string, error) {
	items, err := query(client, ctx, queryID)
	if err != nil {
		return "", err
	}
	total := 0.0
	for _, i := range items {
		total += i.Effort
	}
	msg := fmt.Sprintf("*===== %s / %.0f =====*\n", title, total)
	if len(items) == 0 {
		msg += "Nothing" + "\n"
	}
	for _, i := range items {
		msg += fmt.Sprintf("● %s\n", i)
	}
	return msg, nil
}

func query(client workitemtracking.Client, ctx context.Context, queryID string) ([]printItem, error) {
	id, err := uuid.Parse(queryID)
	if err != nil {
		return nil, err
	}
	resp, err := client.QueryById(ctx, workitemtracking.QueryByIdArgs{Id: &id})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	var list []printItem
	for _, id := range *resp.WorkItems {
		wi, err := client.GetWorkItem(ctx, workitemtracking.GetWorkItemArgs{Id: id.Id})
		if err != nil {
			return nil, err
		}
		effort, ok := (*wi.Fields)["Microsoft.VSTS.Scheduling.Effort"]
		if !ok {
			effort = 0.0
		}
		list = append(list, printItem{
			ID:      *wi.Id,
			Project: (*wi.Fields)["System.TeamProject"].(string),
			Type:    (*wi.Fields)["System.WorkItemType"].(string),
			Title:   (*wi.Fields)["System.Title"].(string),
			Effort:  effort.(float64),
		})
	}
	return list, nil
}

func postSlack(url, msg string) error {
	body, err := json.Marshal(struct {
		Text string `json:"text"`
	}{msg})
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Print(err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("HTTP %s: %v", resp.Status, data)
	}
	return nil
}
