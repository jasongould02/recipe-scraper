#Recipe Scraper
Currently only scrapes wordpress based recipe websites.

Small Golang server that expects a JSON object:
{
	"URL": "[url placed here]"
}

The server will then scrap the data and place information such as
ingredients, instructions, nutritional data and more into a JSON object, which would be sent back to the server that sent the URL.
