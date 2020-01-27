# work-item-reporter

Generate weekly report of Azure DevOps Work Items.

## Example Output

> *mikan's Weekly Report*
> 
> 
> *===== To-Do / 16 =====*
> 
> ● [scrum1] 📘27: Sample Product Backlog Item / 8
> 
> ● [scrum1] 📘33: Sample Product Backlog Item / 8
> 
> *===== DOING / 0 =====*
> 
> ● [agile2] 📘3: Sample User Story / 0
> 
> ● [scrum1] 🏆21: Sample Epic / 0
> 
> ● [scrum1] 🏆37: Sample Epic / 0
> 
> ● [agile2] 📋43: Sample Task / 0
> 
> *===== DONE / 50 =====*
> 
> ● [scrum1] 📘4: Sample Product Backlog Item / 16
> 
> ● [scrum1] 🏆19: Sample Epic / 0
> 
> ● [scrum1] 📘26: Sample Product Backlog Item / 5
> 
> ● [scrum1] 📘28: Sample Product Backlog Item / 8
> 
> ● [scrum1] 📘29: Sample Product Backlog Item / 13
> 
> ● [scrum1] 📘32: Sample Product Backlog Item / 8
> 
> ● [agile3] 🐞39: Sample Bug / 0
> 
> ● [agile2] 📋40: Sample Task / 0
> 
> ● [agile3] 🐞51: Sample Bug / 0

Format:

- each headers: `===== header / total-effort =====`
- each items: `[project] number: title / effort`

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

## License

[BSD 3-clause](LICENSE)

## Author

- [mikan](https://github.com/mikan)
