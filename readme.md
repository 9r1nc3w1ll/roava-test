# Setup Guide
## Pulsar
NB: Pulsar volume is not persistent

## Database
### Create migration
`$ migrate create -seq -ext sql -dir migrations initialize_schema`
NB: Postgres volume is persistent