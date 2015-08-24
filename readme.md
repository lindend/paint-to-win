# Paint to win

## Back end required tools

### Linux installation

**golang** *http://golang.org/doc/install*

**make** *sudo apt-get install make*

### Windows installation

**golang** *http://golang.org/doc/install*

**make** http://gnuwin32.sourceforge.net/packages/make.htm

### Building

Type make in src/paintToWin/lobby and /src/paintToWin/gameserver

## Back end dependencies

### PostgreSQL (Ubuntu)

*sudo apt-get install postgresql-9.3*

*sudo -u postgres psql*

*CREATE DATABASE paint2win;*

*CREATE USER p2wuser WITH PASSWORD 'devpassword';*

*GRANT ALL ON DATABASE paint2win TO p2wuser;*

### Redis (Ubuntu)

Navigate to empty temp folder.

*wget http://download.redis.io/releases/redis-3.0.2.tar.gz*

*dtrx redis-3.0.2.tar.gz*

*cd redis-3.0.2*

*make*

*make test*

*sudo make install*

*cd utils*

*sudo ./install_server.sh*