## General info
This project implements an offer and reservation platform accommodation. The project unites the following departments: SOA, NoSQL and MRS.
	
## Technologies
Project is created with:
* GoLang
* Angular
* Docker
* NoSQL Databases (MongoDB, Cassandra, Neo4j, Redis)

# Service and its database
* auth-service: MongoDB
* profile-service: MongoDB
* accommodation-service: Cassandra
* reservation-service: Cassandra
* rating-service: Cassandra
* notification-service: Neo4j
* image caching: Redis

# Type of communication
* auth-service: synchronous
* profile-service: synchronous
* accommodation-service: synchronous and asynchronous
* reservation-service: synchronous
* rating-service: synchronous
* notification-service: asynchronous
	
## Setup
To run this project, install it locally using npm:

```
$ npm install
$ ng serve --open
```

```
$ docker-compose up --build
```

```
Create TABLES:
CREATE KEYSPACE accommodations WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
CREATE KEYSPACE reservations WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

CREATE TABLE accommodations.accommodations (
    id uuid PRIMARY KEY,
    amenities text,
    base_price double,
    guests int,
    location text,
    name text,
    price double,
    price_strategy text,
    user_id text,
    images list<text>
);

CREATE INDEX ON accommodations.accommodations (user_id);


CREATE TABLE availability (
    id UUID PRIMARY KEY,
    accommodation_id UUID,
    start_date DATE,
    end_date DATE
);


CREATE TABLE prices (
    id UUID PRIMARY KEY,
    accommodation_id UUID,
    start_date DATE,
    end_date DATE,
    amount DOUBLE,
    strategy TEXT  -- 'per_guest' or 'per_unit'
);

CREATE TABLE IF NOT EXISTS ratings (
    id UUID PRIMARY KEY,
    user_id TEXT,
    target_id UUID,
    rating INT,
    comment TEXT,
    rated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ratings2 (
    id UUID PRIMARY KEY,
    user_id TEXT,
    target_id UUID,
    rating INT,
    comment TEXT,
    rated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS average_ratings (
    target_id UUID PRIMARY KEY,
    average_rating FLOAT,
    total_ratings INT
);

CREATE TABLE IF NOT EXISTS reservations (
    id UUID PRIMARY KEY,
    accommodation_id UUID,
    guest_id text,
    start_date date,
    end_date date
);

Neo4j: - cypher-shell -u neo4j -p StrongPassword123
       - MATCH (n) RETURN n;

MongoDB: - mongo
         - use authdb
         - db.users.find().pretty()

Cassandra: - cqlsh
           - use accommodations

Redis: - redis-cli
       - keys *
       - get <key>
```