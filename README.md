# Go Auctions - Auctions Management

This application provides tools to manage auctions, from its creation to bidding.
Its toolset includes checking for the winning bid, batch inserting bids and automatic auction closure give its defined lifetime.

### Execution steps

It is advised to use `docker` to run this application through the command `docker compose up -d`, executed in this project root folder.

To test the application, make sure the containers are running and execute the command `docker compose exec app go test ./...`. There is one test that provides confirmation that the auction is being created and successfully completed for the given lifetime.
