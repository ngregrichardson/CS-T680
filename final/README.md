# How to Use

## Automatic
You can automatically build and run all three APIs plus Redis by running `docker compose up --build`.

## Manual
You can manually build and run each API manually by running `docker compose up --build` in each directory.

**Important** Running the APIs manually will expect a network `ngr27` to exist. By default, the Redis container will create this network. If you are not running the Redis container, you should first create the network with `docker network create ngr27`, then add the `REDIS_URL` environment variable to the individual API `docker-compose.yml` files to point to the Redis instance.

## Testing with Postman
You can use the Postman collection either in the `postman` folder or from here (make sure to use the included environment and adjust it to your needs):<br> [![Run in Postman](https://run.pstmn.io/button.svg)](https://god.gw.postman.com/run-collection/10145298-0b49fab4-126c-4f40-8ff3-bafadd7c300c?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D10145298-0b49fab4-126c-4f40-8ff3-bafadd7c300c%26entityType%3Dcollection%26workspaceId%3Db9b1aa99-edd3-44b2-8076-9a45c4321bb2#?env%5BCS-T680%20Final%20Project%5D=W3sia2V5IjoicG9sbHNfYXBpX3VybCIsInZhbHVlIjoiaHR0cDovL2xvY2FsaG9zdDoxMDgwIiwiZW5hYmxlZCI6dHJ1ZSwidHlwZSI6ImRlZmF1bHQifSx7ImtleSI6InZvdGVyc19hcGlfdXJsIiwidmFsdWUiOiJodHRwOi8vbG9jYWxob3N0OjEwODEiLCJlbmFibGVkIjp0cnVlLCJ0eXBlIjoiZGVmYXVsdCJ9LHsia2V5Ijoidm90ZXNfYXBpX3VybCIsInZhbHVlIjoiaHR0cDovL2xvY2FsaG9zdDoxMDgyIiwiZW5hYmxlZCI6dHJ1ZSwidHlwZSI6ImRlZmF1bHQifV0=)

# About

## General Flow
The general flow of a request:
- **Middleware** (`middleware`): Tracks request and response times for statistics.
- **Router** (`api`): Handles routing to the proper handler and validates inputs as well as formats responses.
- **Service** (`services`): Handles logic for getting, setting, and updating data.
- **Cache** (`cache`): Handles storing data.

Some other packages are available:
- `schema`: Defines types and a few interface methods.
- `utils`: Defines utility functions for a variety of different purposes.

## Hypermedia
The hypermedia integration between these 3 APIs can be seen in almost all of the requests. There are different parameters returned, such as `links` or `pollLinks`, depending on the object, that contain links to other objects. For example, a `pollLinks` object will contain links to the poll's options (such as get, update, and delete). This allows for easy traversal of the data across the APIs. Along with each URLis the method that can be used to make the proper request.

# Votes API Usage

## Endpoints

### Base Url
[`http://localhost:1082`](http://localhost:1082)

### [GET] `/votes/health`
#### Description
Returns the health and statistics of the Votes API.

#### Query Parameters
| Name | Type  | Description                                                      | Required? |
|:----:|:-----:|------------------------------------------------------------------|:---------:|
| hang | `int` | Causes the endpoint to hang for the specified number of seconds. |           |

### [GET] `/votes/:id`
#### Description
Returns a vote specified by `id`.

#### Path Parameters
| Name | Type  | Description          | Required? |
|:----:|:-----:|----------------------|:---------:|
|  id  | `int` | The id of the vote. |     ✅     |

### [POST] `/votes`
#### Description
Creates a new vote.

#### Body
```json
{
    "id": 0,
    "voterId": 0,
    "pollId": 0,
    "optionId": 0
}
```

### [PATCH] `/votes/:id`
#### Description
Updates a vote specified by `id`. Since the method is `PATCH`, only the fields that are specified in the body will be updated (in this case, only `optionId` can be).

#### Path Parameters
| Name | Type  | Description          | Required? |
|:----:|:-----:|----------------------|:---------:|
|  id  | `int` | The id of the vote. |     ✅     |

#### Body
```json
{
    "optionId": 0
}
```

### [DELETE] `/votes/:id`
#### Description
Deletes a vote specified by `id`.

#### Path Parameters
| Name | Type  | Description          | Required? |
|:----:|:-----:|----------------------|:---------:|
|  id  | `int` | The id of the vote. |     ✅     |

### [GET] `/voters`
#### Description
Returns the list of voters.

### [GET] `/polls`
#### Description
Returns the list of polls.

### [GET] `/voters/health`
#### Description
Returns the health and statistics of the Voters API.

#### Query Parameters
| Name | Type  | Description                                                      | Required? |
|:----:|:-----:|------------------------------------------------------------------|:---------:|
| hang | `int` | Causes the endpoint to hang for the specified number of seconds. |           |

### [GET] `/polls/health`
#### Description
Returns the health and statistics of the Polls API.

#### Query Parameters
| Name | Type  | Description                                                      | Required? |
|:----:|:-----:|------------------------------------------------------------------|:---------:|
| hang | `int` | Causes the endpoint to hang for the specified number of seconds. |           |

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