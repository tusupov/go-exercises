# Bank Account

## Params
* `PORT` - address for server listen

## Run width docker
``` bash
$ docker-compose up
```
## API
PUT [`/acoount`](#api-create) - Create account
POST [`/acoount`](#api-deposit) - Deposit (withdrawal) from account
GET [`/acoount`](#api-balance) - Account balance
DELETE [`/acoount`](#api-delete) - Close account

### <a name="api-create"></a>Create account
```
PUT /account
{
  "initialAmount": 100
}
```

### <a name="api-deposit"></a> Deposit (withdrawal) from account
```
POST /account
{
  "amount": 50 // may be negative
}
```

### <a name="api-balance"></a> Account balance
```
GET /account
```

## Example
```
==> GET /account
<== HTTP 400 Bad request
{ "error": "Account is not created" }

==> POST /account
{ "initialAmount": 120 }
<== HTTP 200 OK

==> POST /account
{ "amount": -30 }
<== HTTP 200 OK

==> POST /account
{ "amount": -100 }
<== HTTP 400 Bad request
{ "error": "Not enough money" }

==> GET /account
<== HTTP 200 OK
{ "amount": 90 }

==> DELETE /account
<== HTTP 200 OK

==> GET /account
<== HTTP 400 Bad request
{ "error": "Account is closed" }
```
