# How to Use

## Automatic
You can automatically build and run all three APIs plus Redis by running `docker compose up --build`.

## Manual
You can manually build and run each API manually by running `docker compose up --build` in each directory.

**Important** Running the APIs manually will expect a network `ngr27` to exist. By default, the Redis container will create this network. If you are not running the Redis container, you should first create the network with `docker network create ngr27`, then add the `REDIS_URL` environment variable to the individual API `docker-compose.yml` files to point to the Redis instance.

# Voters API Usage

## Endpoints

### Base Url
[`http://localhost:1081`](http://localhost:1081)

### [GET] `/voters/health`
#### Description
Returns the health and statistics of the Voters API.

#### Query Parameters
| Name | Type  | Description                                                      | Required? |
|:----:|:-----:|------------------------------------------------------------------|:---------:|
| hang | `int` | Causes the endpoint to hang for the specified number of seconds. |           |

### [GET] `/voters`

#### Description
Returns the list of voters.

### [GET] `/voters/:id`

#### Description
Returns a voter specified by `id`.

#### Path Parameters
| Name | Type  | Description          | Required? |
|:----:|:-----:|----------------------|:---------:|
|  id  | `int` | The id of the voter. |     ✅     |

### [POST] `/voters`

#### Description
Creates a new voter.

#### Body
```json
{
    "id": 0,
    "firstName": "string",
    "lastName": "string"
}
```

### [PUT] `/voters/:id`

#### Description
Updates a voter specified by `id`. Since the method is `PUT`, the whole updatable voter object should be sent in the body. Therefore, if `firstName` or `lastName` is not specified, they will be turned to `""`.

#### Path Parameters
| Name | Type  | Description          | Required? |
|:----:|:-----:|----------------------|:---------:|
|  id  | `int` | The id of the voter. |     ✅     |

#### Body
```json
{
    "firstName": "string",
    "lastName": "string"
}
```

### [DELETE] `/voters/:id`

#### Description
Deletes a voter specified by `id`.

#### Path Parameters
| Name | Type  | Description          | Required? |
|:----:|:-----:|----------------------|:---------:|
|  id  | `int` | The id of the voter. |     ✅     |

### [GET] `/voters/:id/votes`

#### Description
Returns the list of votes a voter specified by `id` has made.

#### Path Parameters
| Name | Type  | Description          | Required? |
|:----:|:-----:|----------------------|:---------:|
|  id  | `int` | The id of the voter. |     ✅     |

### [GET] `/voters/:id/votes/:pollId`

#### Description
Returns a vote a voter specified by `id` made on a poll specified by `pollId`.

#### Path Parameters
|  Name  | Type  | Description          | Required? |
|:------:|:-----:|----------------------|:---------:|
|   id   | `int` | The id of the voter. |     ✅     |
| pollId | `int` | The id of the poll.  |     ✅     |

### [POST] `/voters/:id/votes`

#### Description
Creates a vote for a voter specified by `id` on a poll specified by `pollId`.

#### Path Parameters
|  Name  | Type  | Description          | Required? |
|:------:|:-----:|----------------------|:---------:|
|   id   | `int` | The id of the voter. |     ✅     |

#### Body
```json
{
    "pollId": 0
}
```

### [PUT] `/voters/:id/votes/:pollId`

#### Description
Updates a vote made by a voter specified by `id` on a poll specified by `pollId`.

#### Path Parameters
|  Name  | Type  | Description          | Required? |
|:------:|:-----:|----------------------|:---------:|
|   id   | `int` | The id of the voter. |     ✅     |
| pollId | `int` | The id of the poll.  |     ✅     |

### [DELETE] `/voters/:id/votes/:pollId`

#### Description
Deletes a vote made by a voter specified by `id` on a poll specified by `pollId`.

#### Path Parameters
|  Name  | Type  | Description          | Required? |
|:------:|:-----:|----------------------|:---------:|
|   id   | `int` | The id of the voter. |     ✅     |
| pollId | `int` | The id of the poll.  |     ✅     |






# Polls API Usage

## Endpoints

### Base Url
[`http://localhost:1080`](http://localhost:1080)

### [GET] `/polls/health`
#### Description
Returns the health and statistics of the Polls API.

#### Query Parameters
| Name | Type  | Description                                                      | Required? |
|:----:|:-----:|------------------------------------------------------------------|:---------:|
| hang | `int` | Causes the endpoint to hang for the specified number of seconds. |           |

### [GET] `/polls`

#### Description
Returns the list of polls.

### [GET] `/polls/:id`

#### Description
Returns a poll specified by `id`.

#### Path Parameters
| Name | Type  | Description          | Required? |
|:----:|:-----:|----------------------|:---------:|
|  id  | `int` | The id of the poll. |     ✅     |

### [POST] `/polls`

#### Description
Creates a new poll.

#### Body
```json
{
    "id": 0,
    "title": "string",
    "question": "string",
    "options": [{
        "id": 0,
        "title": "string"
    }]
}
```

### [PUT] `/polls/:id`

#### Description
Updates a poll specified by `id`. Since the method is `PUT`, the whole updatable voter object should be sent in the body. Therefore, if `title` or `question` is not specified, they will be turned to `""`.

#### Path Parameters
| Name | Type  | Description          | Required? |
|:----:|:-----:|----------------------|:---------:|
|  id  | `int` | The id of the poll. |     ✅     |

#### Body
```json
{
    "title": "string",
    "question": "string"
}
```

### [DELETE] `/polls/:id`

#### Description
Deletes a poll specified by `id`.

#### Path Parameters
| Name | Type  | Description          | Required? |
|:----:|:-----:|----------------------|:---------:|
|  id  | `int` | The id of the poll. |     ✅     |

### [GET] `/polls/:id/options`

#### Description
Returns the list of options of a poll specified by `id`.

#### Path Parameters
| Name | Type  | Description          | Required? |
|:----:|:-----:|----------------------|:---------:|
|  id  | `int` | The id of the poll. |     ✅     |

### [GET] `/polls/:id/options/:optionId`

#### Description
Returns the option, specified by `optionId`, of a poll, specified by `id`.

#### Path Parameters
|  Name  | Type  | Description          | Required? |
|:------:|:-----:|----------------------|:---------:|
|   id   | `int` | The id of the poll. |     ✅     |
| optionId | `int` | The id of the option.  |     ✅     |

### [POST] `/polls/:id/options`

#### Description
Creates an option on a poll specified by `id`.

#### Path Parameters
|  Name  | Type  | Description          | Required? |
|:------:|:-----:|----------------------|:---------:|
|   id   | `int` | The id of the poll. |     ✅     |

#### Body
```json
{
    "optionId": 0,
    "title": "string"
}
```

### [PATCH] `/polls/:id/options/:optionId`

#### Description
Updates a poll option, where the poll is specified by `id` and the option specified by `optionId`.

#### Path Parameters
| Name | Type  | Description          | Required? |
|:----:|:-----:|----------------------|:---------:|
|  id  | `int` | The id of the poll. |     ✅     |
|  optionId  | `int` | The id of the option. |     ✅     |

#### Body
```json
{
    "title": "string"
}
```

### [DELETE] `/polls/:id/options/:optionId`

#### Description
Deletes an option, specified by `optionId`, of a poll, specified by `id`.

#### Path Parameters
|  Name  | Type  | Description          | Required? |
|:------:|:-----:|----------------------|:---------:|
|   id   | `int` | The id of the poll. |     ✅     |
| optionId | `int` | The id of the option.  |     ✅     |