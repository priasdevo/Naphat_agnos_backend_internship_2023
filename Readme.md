How to deploy locally
1. Clone the repository
2. Change the user, password, dbname of postgres to correct local protgres database
3. Create the table in database with following command
    CREATE TABLE logs (
        id SERIAL PRIMARY KEY,
        time TIMESTAMP NOT NULL,
        status INTEGER NOT NULL,
        latency TEXT NOT NULL,
        client_ip TEXT NOT NULL,
        method TEXT NOT NULL,
        path TEXT NOT NULL
    );
4. At the root folder run "go run main.go"
How to run unit test
1. At the root folder run "go test"