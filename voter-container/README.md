# How to Use

## Automatic
You can automatically build and run the Voter API in two ways:
- `make docker-start`
- `docker compose up --build`

## Manual
You can manually build and run the Voter API in a few different ways:
- `make docker-build && make docker-run`
- `docker compose build && docker compose up`
- `docker build --tag ngregrichardson/voter-api -f ./Dockerfile . && docker compose up`

## Configuring Redis
Redis is set up to automatically start in its own Docker container before the Voter API starts up (in `docker-compose.yml`). Since the containers are on the same network, the Voter API defaults to `redis:6379` as the Redis connection url. To change this, set `REDIS_URL` to your preferred url for the Voter API container.

## Using Published Version
The published Docker image is [`ngregrichardson/voter-api`](https://hub.docker.com/r/ngregrichardson/voter-api) and can be used without the Dockerfile or source code present.

# Voter API Usage

## Makefile Commands
- `run`: Run the Voter API from code
- `run-bin`: Run the Voter API executable
- `get-voters`: Get all voters
- `get-voter`: Get a voter using id=<id>
- `create-voter`: Create a voter using id=<id>, firstName=<firstName>, lastName=<lastName>
- `update-voter`: Update a voter using id=<id>, firstName=<firstName>, lastName=<lastName>
- `delete-voter`: Delete a voter using id=<id>
- `get-voter-history`: Get a voter's vote history using id=<id>
- `get-vote`: Get a voter's vote using id=<id>, pollId=<pollId>
- `create-vote`: Create a voter's vote using id=<id>, pollId=<pollId>
- `update-vote`: Update a voter's vote using id=<id>, pollId=<pollId>
- `delete-vote`: Delete a voter's vote using id=<id>, pollId=<pollId>
- `health`: Health check using hang=<hang> (optional)

## API Responses
All responses from the API (successful or not) are formatted in the same wrapper. This makes the API easier to work with in some front-end request frameworks where it can be hard to tell if a response is successful.

```json
{
    "statusText": "string",
    "statusCode": 0,
    "message": "string",
    "isSuccess": true/false,
    "data": any
}
```

In use, this is what an actual response from `/voters/health` looks like.

```json
{
    "statusText": "OK",
    "statusCode": 200,
    "message": "Success",
    "isSuccess": true,
    "data": {
        "stats": {
            "averageRequestTime": "0s",
            "errors": 0,
            "requests": 1,
            "totalRequestTime": "0s",
            "uptime": "10.8922853s"
        },
        "version": "1.0.0"
    }
}
```

## Endpoints

### Base Url
[`http://localhost:1080/v1`](http://localhost:1080/v1)

### [GET] `/voters/health`
#### Description
Returns the health and statistics of the Voter API.

#### Query Parameters
| Name | Type  | Description                                                      | Required? |
|:----:|:-----:|------------------------------------------------------------------|:---------:|
| hang | `int` | Causes the endpoint to hang for the specified number of seconds. |           |

### [GET] `/voters`

#### Description
Returns the list of voters.

### [GET] `/voter/:id`

#### Description
Returns a voter specified by `id`.

#### Path Parameters
| Name | Type  | Description          | Required? |
|:----:|:-----:|----------------------|:---------:|
|  id  | `int` | The id of the voter. |     ✅     |

### [POST] `/voter/:id`

#### Description
Creates a new voter. The `id` of the voter is specified in the path parameters, any passed through the body will be ignored.

#### Path Parameters
| Name | Type  | Description          | Required? |
|:----:|:-----:|----------------------|:---------:|
|  id  | `int` | The id of the voter. |     ✅     |

#### Body
```
{
    "firstName": "string",
    "lastName": "string"
}
```

### [PUT] `/voter/:id`

#### Description
Updates a voter specified by `id`. Since the method is `PUT`, the whole updatable voter object should be sent in the body. Therefore, if `firstName` or `lastName` is not specified, they will be turned to `""`.

#### Path Parameters
| Name | Type  | Description          | Required? |
|:----:|:-----:|----------------------|:---------:|
|  id  | `int` | The id of the voter. |     ✅     |

#### Body
```
{
    "firstName": "string",
    "lastName": "string"
}
```

### [DELETE] `/voter/:id`

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

### [GET] `/voter/:id/votes/:pollId`

#### Description
Returns a vote a voter specified by `id` made on a poll specified by `pollId`.

#### Path Parameters
|  Name  | Type  | Description          | Required? |
|:------:|:-----:|----------------------|:---------:|
|   id   | `int` | The id of the voter. |     ✅     |
| pollId | `int` | The id of the poll.  |     ✅     |

### [POST] `/voter/:id/votes/:pollId`

#### Description
Creates a vote for a voter specified by `id` on a poll specified by `pollId`.

#### Path Parameters
|  Name  | Type  | Description          | Required? |
|:------:|:-----:|----------------------|:---------:|
|   id   | `int` | The id of the voter. |     ✅     |
| pollId | `int` | The id of the poll.  |     ✅     |

### [PUT] `/voter/:id`

#### Description
Updates a vote made by a voter specified by `id` on a poll specified by `pollId`.

#### Path Parameters
|  Name  | Type  | Description          | Required? |
|:------:|:-----:|----------------------|:---------:|
|   id   | `int` | The id of the voter. |     ✅     |
| pollId | `int` | The id of the poll.  |     ✅     |

### [DELETE] `/voter/:id`

#### Description
Deletes a vote made by a voter specified by `id` on a poll specified by `pollId`.

#### Path Parameters
|  Name  | Type  | Description          | Required? |
|:------:|:-----:|----------------------|:---------:|
|   id   | `int` | The id of the voter. |     ✅     |
| pollId | `int` | The id of the poll.  |     ✅     |