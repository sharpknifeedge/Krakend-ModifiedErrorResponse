# Health check util
## Iris usage
* Note that it\`s party  must be in the root path, Therefore, This function will be accessible without a prefix `http://SERVER/health`.

``` go

	app := iris.New()
	party := app.Party("/")
	// add health check route
	http_helper.IrisHelthCheck(party, version.Version)

```
## Gorilla Mux usage
* Note that it\`s router  must be in the root path, Therefore, This function will be accessible without a prefix `http://SERVER/health`.

``` go

	baseRouter := mux.NewRouter()
	// add health route
	http_helper.MuxHelthCheck(baseRouter, version.Version)

```