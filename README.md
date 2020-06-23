# AdMediation
AdMediation - Assignment

# Create and Deploy App
- Go to google cloud platform and open App Engine in Console
- Create a new project and within it a new folder called admediation on path /home/{your_user}/gopath/src/admediation
- Copy files admediation.go, app.yaml and output.txt in the before created folder
- Import all needed imports: <br/>
`$ go get github.com/gorilla/mux`<br/>
`$ go get github.com/thedevsaddam/gojsonq`
- Run command `$ gcloud app deploy`
- Run command `$ gcloud app browse` to see the URL where the App runs

# Usage - REST Examples
- Get All Ad Networks <br/>
Method:`GET`<br/>
URL:`/adnetworks`
- Get Ad Networks by Ad type (Banner, Reward, Interstitial)<br/>
Method:`GET`<br/>
URL:`/adnetwork/{adType}`
- Custom query params<br/>
Method:`GET`<br/>
URL:`/adnetwork?param1=value1&param2=value2...`
- Create a new Ad Network (ID is autoincremented)<br/>
Method:`POST`<br/>
URL:`/adnetwork`<br/>
Body:
```json
{
  "description": "admob",
  "value": 6,
  "platform": "android",
  "osversion": "9",
  "appname": "my talking angela",
  "appversion": "2.4.2",
  "countrycode": "slo",
  "adtype": "banner"
}
```

- Update value in existing Ad Network<br/>
Method:`POST`<br/>
URL:`/adnetwork/{id}`<br/>
Body:
```json
{
  "value": 6
}
```
- Delete existing Ad Network<br/>
Method:`DELETE`<br/>
URL:`/adnetwork/{id}`

# Attributes
id = int64<br/>
description = string<br/>
value = int<br/>
platform = string<br/>
osversion = string<br/>
appname = string<br/>
appversion = string<br/>
countrycode = string<br/>
adtype = string

# Assumptions
- For added logic `AdMob doesn’t work on any android os version 9 so it should be returned to the app, but works on other os versions`, assumed it meant shouldn't and not should
- No database was created to save the values, everything is saved in `output.txt`
- When Query-ing data and the response is empty, the response is replaced with a list of all Ad Networks (if the parameter platform was used in Query, the response is also Query-ed to that platform)
