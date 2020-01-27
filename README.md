# work-item-reporter

Generate weekly report of Azure DevOps.

## Prerequisites

- Azure Boards Query UUIDs (shown in browser address bar)
    - My to-do items
        - `State` = `Accepted`
        - Or `State` = `Ready`
        - Or `State` = `Approved`
        - And `Assigned To` = `(YOU)`
    - My doing items
        - `State` = `Committed`
        - Or `State` = `Doing`
        - Or `State` = `In Progress`
        - Or `State` = `Active`
        - And `Assigned To` = `(YOU)`
    - My done items last week
        - `Closed Date` >= `@StartOfDay('-7d')`
        - And `Assigned To` = `(YOU)`
- Personal Access Token
- Slack Incoming Webhook URL

## Usage

Standard output (for testing)

```
go run main.go -o <ORG> -t <TOKEN> -todo <TODO_QUERY> -doing <DOING_QUERY> -done <DONE_QUERY> -n "<NAME>"
```

Slack

```
go run main.go -o <ORG> -t <TOKEN> -todo <TODO_QUERY> -doing <DOING_QUERY> -done <DONE_QUERY> -n "<NAME>" -s <WEBHOOK>
```

Cron tab

```
0 10 * * 1 work-item-reporter -o <ORG> -t <TOKEN> -todo <TODO_QUERY> -doing <DOING_QUERY> -done <DONE_QUERY> -n "<NAME>" -s <WEBHOOK>
```

## Author

- [mikan](https://github.com/mikan)
