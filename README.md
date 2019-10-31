# Instatask API
### **Staging host: [instatasks.herokuapp.com](instatasks.herokuapp.com)**
**Ping:** [instatasks.herokuapp.com/ping](instatasks.herokuapp.com/ping)

# API endpoints examples
## Admin context

Admin endpoints mount to ***/admin***  protected route.
For route access used BasicAuthorization strategy

| HEADER | VALUE |
|--|--|
|Content-Type|"application/json"|
|Authorization|"Basic " + Base64 encoded "username:password"|

username and password on server stored in environment variables:
```go
SUPERADMIN_USERNAME=superadmin
SUPERADMIN_PASSWORD=superadminsupersecretpassword
```
> **Note:**  On the **Production** server all **Environments Variables** stored in ***.env*** file (see ***.env.exemple*** by example usage).  On the **staging** hosting  **Heroku** ([https://dashboard.heroku.com/apps/instatasks](https://dashboard.heroku.com/apps/instatasks))  they are in the **Settings** section as the **Reveal Config Vars**
> 
Authorization HEADER example:
```go
Basic c3VwZXJhZG1pbjpzdXBlcmFkbWluc3VwZXJzZWNyZXRwYXNzd29yZA==
```
 ### Create User Agent 
**POST */admin/useragent***

*Request Body*
```json
{
	"name": "useragent1", 
	"activitylimit": 0,
	"like": true,
	"follow": true,
	"pricelike": 1,
	"pricefollow": 5,
	"pricerateus": 20
}
```

> In this exemple is **set of default values**. Rquest **Required only**
> User Agent **"name"** key.

*Response Body*

```json
{
	"activitylimit":1,
	"follow":true,
	"like":true,
	"name":"useragent1",
	"pricefollow": 5,
	"pricelike": 1,
	"pricerateus": 20,
	"rsa_public_key":"-----BEGIN RSA PUBLIC KEY-----\nMII...QAB\n-----END RSA PUBLIC KEY-----\n"
}
```
### Get RSA Public Key
**GET */admin/useragent/pkey***

*Request Body*
```json
{
	"name": "useragent1"
}
```
*Response Body*

```json
{
	"rsa_public_key":"-----BEGIN RSA PUBLIC KEY-----\nMII...QAB\n-----END RSA PUBLIC KEY-----\n"
}
```
## Main API context

### Accaunt route
**POST */accaunt***

|HEADER|VALUE|
|--|--|
|Content-Type|application/json|

*Request Body*
```json
{
	"instagramid": 666,
	"deviceid":    "device1"
}
```
All keys required.

*Response Body*

```json
{
	"coins":0,
	"instagramid":666,
	"rateus":true
}
```

### User Agent settings

**POST */setting***

| HEADER | VALUE |
|--|--|
|Content-Type|application/json|
|User-Agent|useragentname|

*Request Body*
```json
not required
```
All keys required.

*Response Body*

```json
{
	"activitylimit":0,
	"follow":true,
	"like":true,
	"pricefollow":5,
	"pricelike":1,
	"pricerateus":20
}
```
	
### Crete Task
**POST */newwork***

| HEADER | VALUE |
|--|--|
|Content-Type|application/json|
|User-Agent|useragentname|

*Request Body*
```json
{
	"instagramid": 777,
	"count": 10,
	"type": "like",
	"mediaid":"mediaid1",
	"photourl":"url/blabla", 
	"instagramusername":"url/blabla"
}
```
*Response Body*

```json
{
	"coins":10
}
```
### Tasks history
**POST */history***

| HEADER | VALUE |
|--|--|
|Content-Type|application/json|

*Request Body*
```json
{
	"instagramid": 777
}
```
*Response Body*

```json
[
	{	
		"taskid": "15",
		"created_at": "2019-10-31T14:22:44.927428295Z",
		"deleted_at": null,
		"taskid": "15",
		"type": "like",
		"count": 20,
		"left_counter": 20,
		"photourl": "url/blabla",
		"instagramusername": "url/blabla",
		"mediaid": "mediaid2",
		"instagramid": 777
	}, {......}....
]
```
###  Get Tasks
**POST */gettasks***

| HEADER | VALUE |
|--|--|
|Content-Type|application/json|


*Request Body*
```json
{
	"instagramid": 777,
	"type":"all"
}
```
*Response Body*

```json
[
	{
		"taskid":"2",
		"type":"like",
		"photourl":"url/blabla",
		"instagramusername":"url/blabla",
		"mediaid":"mediaid2"
	},{....}...
]
```
###  Done Task
**POST */done***

| HEADER | VALUE |
|--|--|
|Content-Type|application/json|
|User-Agent|useragentname|

*Request Body*
```json
{
	"instagramid": 777,
	"taskid": "19",
	"status": "cancel" // or any not "cancel" as "ok"
}
```
**If status cancel:**
http status code: 205

```json
{"error":"Task Canceled"}
```
**If status ok:**
http status code: 200

```json
{"coins":15}
```

###  Rateus
**POST */rateus***

| HEADER | VALUE |
|--|--|
|Content-Type|application/json|
|User-Agent|useragentname|

*Request Body*
```json
{
	"instagramid": 888
}
```
**If already done:**
http status code: 406

```json
{"error":"Task already done"}
```
**If status ok:**
http status code: 200

```json
{"coins":35}
```