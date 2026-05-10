# Ticketing

## Prepare service

```
go mod tidy
```

### If using SQL database, run this commands for migration

UP migration:
```
make migration
```

Rollback migration:
```
make migration down
```

## Build and run service
```
make run
```