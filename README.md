Game Library
=============

A web-hosted library of free and discounted games

Gathers games from Amazon Prime, Epic Games, Steam, and GOG


Requirements
------------

- **Go** 1.23.2+ (Only needed if planning modifications)

If no modifications are desired then just run gsource-lib


Dependencies
------------

-[goose](https://github.com/pressly/goose)
-[postgresql](https://www.postgresql.org)

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
