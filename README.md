# Chat Service

## Way 1 : run service without `micro server`.

run server:
```bash
$ go run main.go --profile=dev --service_address=localhost:8080
```

run client:
```bash
$ go run client/main.go --profile=dev
```

you can debug ....

## Way 2 : run service in `micro server`.

### Calling the service

You can call the service via the CLI:

Create a chat:
```bash
> micro chat new --user_ids=JohnBarry
{
	"chat_id": "3c9ea66c-d516-45d4-abe8-082089e18b27"
}
```

Send a message to the chat:
```bash
> micro chat send --chat_id=bed4f0f0-da12-46d2-90d2-17ae1714a214 --user_id=John --subject=Hello --text='Hey Barry'
{}
```

View the chat history
```bash
> micro chat history --chat_id=bed4f0f0-da12-46d2-90d2-17ae1714a214
{
	"messages": [
		{
			"id": "a61284a8-f471-4734-9192-640d89762e98",
			"client_id": "6ba0d2a6-96fa-47d8-8f6f-7f75b4cc8b3e",
			"chat_id": "bed4f0f0-da12-46d2-90d2-17ae1714a214",
			"user_id": "John",
			"subject": "Hello",
			"text": "Hey Barry"
		}
	]
}
```