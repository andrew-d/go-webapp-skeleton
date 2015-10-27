# go-webapp-skeleton

This repo contains a skeleton of a Go webapp.  It has some useful features,
mostly cribbed from other, similar projects around the Internet:

- [Automatic migration handling][migrate]
- [Embedded static files][efiles] and [templates][etmpls], including
  [layout support][elayouts].
- [Graceful shutdown support][graceful]
- Robust logging

All tied together with a useful set of tooling.  See below for more information.


## Usage

1. Clone this repository into your project's location:  
   `git clone https://github.com/andrew-d/go-webapp-skeleton.git $GOPATH/github.com/your-user/your-project`
1. Run the `rename.sh` script in order to rename this project and point it at
   the new import path.
1. Type `make clean all` in order to clean the project and rebuild it.  
   *Note*: you need [`go-bindata`][bindata] installed and in your `$PATH` in
   order for the build to complete.


## Tooling

This project uses [gvt][gvt] in order to manage dependencies.  It comes with
all the dependencies already vendored as part of the repository.  In order to
update a given dependency, you can run `gvt update path/to/dep`.


## Organization

**Note**: somewhat of a work in progress

- The base of the project can be found in the `main.go` file, which imports
	most of the other packages and kicks everything off.  You probably won't need
	to modify this as much.
- The `model` directory contains database models.  These models should not
	interact directly with the database.
- The `datastore` directory is responsible for mapping the models to the
	database in use.  It defines interfaces which provide the interface to
	interact with the underlying, concrete, datastore.
- The `datastore/database` directory contains the actual code that is
	responsible for interacting with the underlying database.
- The `datastore/migrate` directory contains the SQL code for migrations
	performed by the app upon startup.
- The `handler` directory contains useful functions that are generic between
	API and frontend routes.
- The `handler/api` directory contains the API route handler functions.
- The `handler/frontend` directory contains the frontend route handler
	functions.
- The `router` directory contains the main router, which registers each of the
	handler functions on their respective routes.




[migrate]: https://github.com/andrew-d/go-webapp-skeleton/blob/master/datastore/migrate/migrate.go
[efiles]: https://github.com/andrew-d/go-webapp-skeleton/tree/master/static
[etmpls]: https://github.com/andrew-d/go-webapp-skeleton/tree/master/handler/frontend/templates
[elayouts]: https://github.com/andrew-d/go-webapp-skeleton/tree/master/handler/frontend/layouts
[graceful]: https://github.com/andrew-d/go-webapp-skeleton/blob/master/main.go#L124
[bindata]: https://github.com/jteeuwen/go-bindata
[gvt]: https://github.com/FiloSottile/gvt
