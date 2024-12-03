# Registry Service

## API reference

All actions are subpath of `/equipment`.

### Create \[POST\]

`/` -- create new piece of equipment.

Required JSON parameters (`id` is assigned automatically, `status` is set to 0) are:

* `kind` (0...3) -- kind of piece of equipment;
* `parameters` (JSON) -- other parameters;

Response:

* 201 (Created) on success. Assigned `id` is returned as JSON string.
* Error on fail.

Sending array of JSON objects -- create multiply new pieces of equipment.

* 201 (Created) on success. Assigned `id`s are returned as JSON array of strings.
* Error on fail.

### Update \[PATCH\]

`/{id}` -- update existing piece of equipment with given `id`.

Optional JSON parameters (but at least one is required):

* `status` (0...2) -- new status value;
* `parameters` (JSON) -- new parameters;

Response:

* 204 (No content) on success.
* 404 (Not found) if `id` was not found.
* Other errors on fail.

`/{id0,id1...}` -- update existing pieces of equipment with given `id`s.

* 200 (OK) on success. Returned JSON object contains fields:
    - `"updated"` -- array of updated `id`s (missing if no requested `id` was found);
    - `"unfound"` -- array of not found `id`s (missing if every requested `id` was found).
* Errors on fail.

`/` -- update existing pieces of equipment filtered by URL-parameters (see section **Filtering** below).

* 200 (OK) on success. Updated Ids are returned as JSON array of strings.
* Errors on fail.

### Delete \[DELETE\]

`/{id}` -- delete existing piece of equipment with given `id`.

Response:

* 204 (No content) on success.
* 404 (Not found) if `id` was not found.
* Other errors on fail.

`/{id0,id1...}` -- delete existing pieces of equipment with given `id`s.

* 200 (OK) on success. Returned JSON object contains fields:
    - `"deleted"` -- array of deleted `id`s (missing if no requested `id` was found);
    - `"unfound"` -- array of not found `id`s (missing if every requested `id` was found).
* Errors on fail.

`/` -- delete existing pieces of equipment filtered by URL-parameters (see section **Filtering** below).

### Find \[GET\]

`/{id}` -- find piece of equipment with given `id`.

Response:

* 200 (OK) on success. Returned JSON object contains equipment data.
* 404 (Not found) if `id` was not found.
* Other errors on fail.

`/{id0,id1}` -- find pieces of equipment with given `id`s.

Response:

* 200 (OK) on success. Returned JSON array contains fields:
    `"found"` -- array of found *pieces of equipment* (missing if no requested `id` was found);
    `"unfound"` -- array of not found `id`s (missing if every requested `id` was found).
* Errors on fail.

`/` -- find pieces of equipment filtered by URL-parameters (see section **Filtering** below).

### Filtering
  
Optional filtering parameters (specified in URL after `?`):

* `kind (0...3)` -- equipment of given kind; multiply comma separated values to include multiply kinds; prevents using `no_kind`;
* `no_kind (0...3)` -- equipment of any kind except given; multiply comma separated values to exclude multiply kinds; prevents using `kind`;
* `status (0...3)` -- equipment having given operational status; multiply comma separated values to include multiply statuses; prevents using `no_status`;
* `no_status (0...2)` -- equipment having any status except given; multiply comma separated values to exclude multiply statuses; prevents using `status`;
* `created_since (timestamp)` -- pieces of equipment created not earlier than;
* `created_until (timestamp)` -- pieces of equipment created not later than;
* `created_since (timestamp)` -- pieces of equipment updated not earlier than;
* `created_until (timestamp)` -- pieces of equipment updated not later than.
