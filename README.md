## emsql

Query the Equinix Metal API with SQL.

### Installation

```
make build
```

will produce a binary (will need Go installed) 

### Usage

`emsql` expects the `PACKET_AUTH_TOKEN` env var to be set to a valid API token.
The first argument should be a SQL query.

```sql
SELECT hostname, json_extract(facility, '$.code') AS facility, json_extract(os, '$.name') AS os, json_extract(plan, '$.slug') AS plan FROM devices('9bdabaa0-fa63-478a-8e34-64785eba8c14')
```

More documentation coming...
