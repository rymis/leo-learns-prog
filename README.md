Leo learns programming
======================
This project contains a simple environment to teach our son how to write programs in JavaScript.

API
---
Each API call for server has prefix `/api/`. Then each method name is appended to this prefix, so to call method `test` you need to send POST request to URL `/api/test`. Request body contains arguments of the request as a JSON. Function returns result as a JSON message.

If request is successful, resulting JSON contains a field named `result` with message described in following sections. If there was an error function returns JSON with field `error` containing description of the error.

### commit
Allows to commit a new version of the file.

Request:
```
{
    "name": "name of a file",
    "content": "content of a file",
    "comment": "user comment for this change"
}
```

Response - single string representing commited version.

### checkout
Get file content for the last version of for specific version

Request:
```
{
    "name": "name of a file",
    "version": "version of a file or empty string for the last version"
}
```

Response - content of the object.

### versions
List all known versions of the object.

Request:
```
"name of the file"
```

Response:
```
{
    "name": "name of the object",
    "versions": [
        {
            "version": "version string",
            "time": 12356, // JavaScript time
            "comment": "commit message",
            "parent": "version string for previous version"
        },
        ...
    ]
```

### list
List all known objects for the user.

Request:
```
"shell mask of files you want to list"
```

Response:
```
[ "file1", "file2", ... ]
```
