## Node API description

### General information

Stateless service to run and monitor xray server.

Support add and remove users and get its traffic stats.

API is protected, available only for NodeManager.

### `POST /node/start`

**Description:** Start xray node with added users

**Params:**
```json
{
    "users": [
        {
            "name": "User1",
            "key": "UserKey1"
        },
        {
            "name": "User2",
            "key": "UserKey2"
        }
        ...
    ]
}
```

**Response:** ```200 OK```

### `POST /node/stop`

**Description:**
Stop xray node

**Params:**
`None`

**Response:**
`200 OK`

### `GET /node/status`

**Description:**
Get node status

**Params:**
`None`

**Response:**
```json
{
    "status": "Starting/Running/Stopping/Stopped",
    "cpuload": 75,    // cpu load, %
    "memoryload": 34, // memory load, %

    "userstats": [
        {
            "name": "User1",
            "incoming": 123 // incoming traffic, Kb
            "outcoming": 7  // outcoming traffic, Kb
        },
        {
            "name": "User2",
            "incoming": 3 // incoming traffic, Kb
            "outcoming": 0  // outcoming traffic, Kb
        },
        ...
    ]
}
```


### `POST /node/stop`

**Description:**
Stop xray node

**Params:**
`None`

**Response:**
`200 OK`

### `POST /users/add`

**Description:**
Add users to node

**Params:**
```json
{
    "users": [
        {
            "name": "User3",
            "key": "UserKey3"
        },
        {
            "name": "User4",
            "key": "UserKey4"
        }
        ...
    ]
}
```

**Response:**
  - `200 OK`
  - `400 Node not running`

### `POST /users/del`

**Description:**
Delete users from node

**Params:**
```json
{
    "users": [
        {
            "name": "User3",
            "key": "UserKey3"
        },
        {
            "name": "User4",
            "key": "UserKey4"
        }
        ...
    ]
}
```

**Response:**
  - ```200 OK```
  - ```400 Node not running```


## Node Manager API description

### General information

- Stateless service to monitor xray nodes.

- Add or remove users

- Add or remove nodes

- Provide available nodes for users


### `POST /nodes/add`

**Description:**
Add new node

**Params:**
```json
{
    "url": "https://path-to.node"
    "accesskey": "secret-node-key"
}
```

**Response:**
```200 OK```

### `POST /nodes/del`

**Description:**
Delete node

**Params:**
```json
{
    "url": "https://path-to.node"
}
```

**Response:**
```200 OK```


### `POST /users/new`

**Description:**
Create new user, keep user added to all available nodes

**Params:**
```json
{
    "namesuggest": "name" // basic human-readable name
}
```

**Response:**
```json
{
    "subscription": "https://xray-node.manager/subscription/{user-id}" // page with xray nodes for user
}
```

### `POST /users/del`

**Description:**

Delete user, remove it from all nodes

**Params:**
```json
{
    "id": "user-id"
}
```

**Response:**
```200 OK```

### `GET /subscription/{user-id}`
```json
{
    [
        {
            // connection config for node A
        },

        {
            // connection config for node B
        }
    ]
}
```

### `GET /`

**Description:**

HTML (SSR) page with:
- list of users (and traffic stats)
- list of nodes

**Params:**
None

**Response:**
```200 OK + index.html```


## Node Manager database schema

### Users

- `UserID (int)`
- `Name (string)`

### Nodes

- `NodeID (int)`
- `Address (string)`
- `Key (string)`
- `Active (book)`

### NodeUsers

- `UserID (int)`
- `NodeID (int)`
- `Connected (bool)`
- `TotalIncoming (int, Kb)`
- `TotalOutcoming (int, Kb)`
- `LastCheckIncoming (int, Kb)`
- `LastCheckOutcoming (int, Kb)`

## Required Practices

### Logging

Log requests and responces in nodes and node manager

### Database

Redis with RDB and manual saving on add/del users or nodes.

Persistent users and nodes, may lost some traffic statistics.

### Encryption

Connection between nodes and nodemanager is encrypted

### Parallelism

Nodes and Manager: incoming requests are processed in different goroutines

Manager: background goroutine to check nodes state


## Iterations

### 1. Node basic functionality

Add or remove users, get users stats, start/stop server

### 2. Node manager: watchdog functionality

Connect to nodes and database, restart failed nodes

### 3. Node manager: users adding/deleting

Process add/del user requests.

### 4. Node manager: users subscriptions page

Show users subscription page

### 5. Node manager: UI

SSR monitoring page

### 5. Requests advanced functionality

Logging, encryption, auth

### 6. Node manager: users statistics

Add users statistics 
