# hub-server

## Simple TCP-based chat server. Client can be found: https://github.com/alankritjoshi/hub-client

## Setup
1. `go run server.go 1234`

## How it works
1. Use SyncMap for maintaining concurrent access to different clients
2. Use Go Routines for handling connection
3. Use handlers for input:
    1. "whoami" -> returns clientID of the requesting client
    2. "whoelse" -> returns space-separated clientIDs connected to server apart from requesting client
    3. "send [message] [clientID 1] [clientID 2] ... [clientID 3]" -> sends [message] to specified space-separated clientIDs

## Checks in place
1. Hard limit of 255 simultaneous connections
2. Rejection of commands other than described above

## Testing

### Manual Testing
1. Tested basic functionality
2. Spun up 100 chat connections using a script and applied a sequence of 1. whoelse and 2. send [random_word] [list of clientIDs from 1.]

### Unit Testing
1. Unable to make connection mocking work in interest of time so added a unit test outline once connection mocking is available

### To Improve
1. Use structs to simplify data sharing for handlers
2. Use better error handling structs/package for consistent error to the client
3. Figure out how to mock connections for unit tests
4. Find an optimal way of making the "whoelse" handler work. I forgot that there's no concept of sets in GoLang but the Map can be utilized better instead of printing each clientID on the fly.
