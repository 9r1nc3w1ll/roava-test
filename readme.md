# Setup Guide
## Pulsar
- Pulsar volume is not persistent (this is delibrate)
<br>

## Database
- Postgres volume is persistent
- Check and update database url before running migration.

Create migration

```$ migrate create -seq -ext sql -dir migrations initialize_schema```

<br>

## Proto File
Pro files are git ignored so remeber to run `make proto_generate`

<br>

## What is Left?
These are the list of things I haven't done
- Unit test deathstar gRPC methods.
- Unit test destroyer gRPC methods.
- Add Hashicorp vault.
