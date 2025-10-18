## This is an RSS feed aggregator that we call _Gator_!

It's a multi-user **CLI tool** that allows users to:

* Collect RSS feeds from across the internet
* Store the feeds and posts in a PostgreSQL database
* Follow and unfollow RSS feeds of other users
* Look at summaries of the posts in the terminal, with links to the full posts

Websites publish updates to their contents with RSS feeds - and it is possible to gather the data and keep track of your favourite sites, blogs, etc.

What you will need for this application to work:

* Install _Go_

    https://webinstall.dev/golang/

* Install _PostgreSQL_ and set up the database

    https://webinstall.dev/postgres/
    
        sudo -u postgres psql
        CREATE DATABASE gator;
        \c gator
        ALTER USER postgres PASSWORD 'postgres';

    You can type exit to leave the psql shell.

* Download the files from this repository to your computer

    `git clone https://github.com/ValeriiaGrebneva/BlogAggregator`

* Get your connection string
    
    Your connection string is a URL with the format:

    `protocol://username:password@host:port/database`

    In my case, it was `postgres://postgres:postgres@localhost:5432/gator`

* Install _Goose_ and do migrations

    https://github.com/pressly/goose#install
    
    Migrations can be done using (do not forget to change to _your_ connection string): 
    
    `goose -dir sql/schema postgres postgres://postgres:postgres@localhost:5432/gator up`

* Build _Gator_ CLI tool

    Ensure that you are in the Gator repository on your computer, than run:

    `go build`

* Set up _Config_ file

    Manually create a config file in your home directory, `~/.gatorconfig.json`, with the following content (do not forget to change to _your_ connection string):

        {
            "db_url": "postgres://postgres:postgres@localhost:5432/gator"
        }

* Check the names of commands before using

    * _**register** name_ - to register user with the username;
    * _**login** name_ - to switch the current user to the user with the name (if it exist);
    * _**reset**_ - resets the databases;
    * _**users**_ - shows the list of registered users and the current user;
    * _**addfeed** name url_ - registers the rss-feed at url with the name;
    * _**feeds**_ - shows the list of registered feeds;
    * _**follow** url_ - to follow the registered feed by url with the current user;
    * _**unfollow** url_ - to unfollow the registered feed by url with the current user;
    * _**following**_ - shows the list of the feeds the current user is following;
    * _**agg** time-delay_ - fetches the rss feeds, one at a time, from the registered url (starting from oldest) and saves the posts into the database. The fetch happens every time-delay. Time format: 1m (for 1 mintute), 1h (for 1 hour). Please do not run it shorter than 1 minute. This command runs an infinite loop. For this reason, run the second instance of the goGator for the posts aggregation or aggregate and then use Ctrl+C to stop it and continue using other commands.
    * _**browse** number (optional)_ - browses the number of posts (or 2 if not chosen) for the current user, showing newest first, according to the feeds the user is following.

This program used Go, PostgreSQL, sqlc, and goose. Thank you for your attention :)