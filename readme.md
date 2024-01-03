Spotify service

This is a MVP

## Pre-requisites
1. Install go
2. Should have sql server running
3. Configurure will all values - config.go

## Execution

1. Run
          ``` go mod vendor ```

2. Run 
          ``` go run . ```


## Tasks Pending

1. tests

## Spotify apis

          
               curl --location --request POST 'https://accounts.spotify.com/api/token' \
               --header 'Content-Type: application/x-www-form-urlencoded' \
               --data-urlencode 'grant_type=client_credentials' \
               --data-urlencode 'client_id=<CLIENT_ID>' \
               --data-urlencode 'client_secret=<CLIENT_SECRET>'

