CREATE TABLE IF NOT EXISTS drones (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    speed FLOAT NOT NULL,
    range FLOAT NOT NULL,
    charging_time INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    return_time TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(255) PRIMARY KEY,
    range FLOAT NOT NULL,
    drone_id VARCHAR(255) REFERENCES drones(id),
    start_time TIMESTAMP NOT NULL
);