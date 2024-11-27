# Mythic Deals Archive

A web app that frequently collects a list of free and discounted games

Gathers games from Amazon Prime, Epic Games, Steam, and GOG

Games on the discounts page can be sorted by price, discount, or title via the filters drop-down 


Why?
----

As an avid enjoyer of video games, I'm always on the lookout for a good deal.  However, frequenting several sites often to find said deals can be a pain.  I found it especially bothersome when I would forget to check in and miss one, especially if it was a free offering!
To solve this issue I created an application that will gather these games all in one place.  The app gives important details for each game, such as the discount, price, and rating.  Redeeming or buying a game is easy, simply click on the game card and you will be redirected to the specific market page for said game.


Requirements
------------

- **Go** 1.23.2+ (Only needed if planning modifications)

If no modifications are desired then just run gsource-lib


Dependencies
------------

- [goose](https://github.com/pressly/goose)
- [postgresql](https://www.postgresql.org)

	go install github.com/pressly/goose/v3/cmd/goose@latest
	apt install postgresql


Install
-------
Clone the repository and install it as an editable package:

	git clone https://github.com/ScooballyD/gsource-lib
	cd gsource-lib/sql/schema
 	goose postgres postgres://postgres:@localhost:5432/archive up
  	goose postgres postgres://postgres:@localhost:5432/archive up
   	cd ../..
 	go run .


Future Plans
------------
- Get results for humble bundle
- Add user accounts to track owned game
- Implement the ability to redeem or purchase games without needing to leave the application
- Add automatic redemption functionality for free games
- Conglomerate owned games and achievements from various sites allowing collectors and hunters to show off (possibly add a leaderboard)


## Contributing
Any contributions are welcome
Submit pull requests to the main branch
