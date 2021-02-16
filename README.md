# Shakespeare Dex

Shakespeare Dex (inspired by Pokedex) is a simple API that, given a Pokemon name, will return its description in a Shakesperean style.

To do so, it consumes 2 APIs:

* [Shakespeare API](https://funtranslations.com/api/shakespeare#endpoint)

* [PokeAPI](https://pokeapi.co/)

## Running
### Method 1 - Running directly from command line

* Install the latest version of [Go](https://golang.org/) (this project was built using Go 1.15.8)

* In the command line, from the root of the project run: `go run main.go` (the project uses Go modules, which means all dependencies should be downloaded when running this command)

The program runs on port `5000`, so please make sure this is available when running.

### Method 2 - Using Docker

* Install Docker

Run the following commands from the root of the project:
* `docker build -t shakespeare-dex .`
* `docker run -p <host-port>:5000 shakespeare-dex`

Where `<host-port>` will be the port used to query the API.

The app exposes one endpoint. Here's a sample request using httpie:

`http localhost:<port>/pokemon/<pokemon_name>`

## Testing

The program has a suite of unit tests for each client and service. To run all unit tests, please ensure you have go installed (see above) and run the following from the project root:

`go test ./...`

## Future Improvements

### Rate Limiting and Circuit Breakers

Due to the simple nature of this project, rate limiting abd circuit breaking was not taken into account.

For the rate limiting case, this means that anyone can make a very high number of requests, which would then be propagated to the 3rd party APIs. For the particular case of the Shakespeare API, we are only allowed 5 requests per hour, so the chances of Pokepeare getting rate limited are high. 

This program also doesn't deal with being rate limited by third party APIs, and simply fails requests to the consumers when this happens. A possible implementation would be to halt requests until rate limiting ban is lifted, or apply exponential backoff mechanisms (which isn't practical since consumers would be waiting for a response).

Circuit breaking would be useful to not overload any 3rd party API with requests if the API is down. Given the scale of this exercise, this was not implemented.

### Shakespeare and Pokeapi error messages

The current implementation does not read potential error messages sent by the 3rd party APIs - as they are just basic plain text responses. In a real world scenario, these could be very benificial for a better service operation (for example, with appropriate codes that would lead to better error handling), and would be used.

### Integration Tests

The program relies solely on unit tests to guarantee the correct functionality. In a production environment, Integration tests should also be taken into account (for example by having a script that does a request to the API and checks if the response is as expected).

### Environment Variables

For simplicity, server and API configurations are hardcoded. In a production environment these would be populated by means of environment variables (for example to configure the port of the server).

### A more generic HTTP Client

As it stands, the Shakespeare and Pokemon API clients have some duplicated code to handle http requests. Since the core functionality of these requests is unlikely to change, it's OK to leave it as is for this program. However, if more clients/APIs were to be introduced, it would make sense to create a shared client that did most of the heavylifting.


