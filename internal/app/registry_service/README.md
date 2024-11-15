Registry Service
================

API reference
-------------

- `/equipment/`
  + `/` \[POST\] -- create new piece of software. Required JSON parameters (`id` is assigned automatically, `status` is set to 0):
    * `kind (0...3)` -- kind of piece of equipment;
    * `parameters {JSON}` -- other parameters;
  + `/{id}` \[PATCH\] -- edit existing piece of software with given `id` (if `id` does not exist, the response will contain error). Optional JSON parameters (but at least one is required):
    * `status` -- new status value;
    * `parameters` -- new parameters;
  + `/` \[GET\]-- list pieces of equipment; optional filtering `GET`-parameters:
    * `kind (0...3)` -- equipment with given kind; multiply comma separated values to include multiply kinds; prevents using `no_kind`;
    * `no_kind (0...3)` -- equipment with any kind except given one; multiply comma separated values to exclude multiply kinds; prevents using `kind`;
    * `status (0...3)` -- equipment having given operational status; multiply comma separated values to include multiply statuses; prevents using `no_status`;
    * `no_status (0...2)` -- equipment having any status except given one; multiply comma separated values to exclude multiply statuses; prevents using `status`;
    * `created_since (timestamp)` -- pieces of equipment created not earlier than;
    * `created_until (timestamp)` -- pieces of equipment created not later than;
    * `created_since (timestamp)` -- pieces of equipment updated not earlier than;
    * `created_until (timestamp)` -- pieces of equipment updated not later than;
  + `/{id}` \[GET\] -- the piece of software with given id;
  + `/{id}` \[DELETE\] -- delete the piece of software with given id  (if `id` does not exist, the response will contain error). 
