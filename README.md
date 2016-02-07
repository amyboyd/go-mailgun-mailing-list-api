Mailgun mailing list API
========================

A simple API, written in Go, for subscribing new emails to a Mailgun mailing list.

Usage
-----

Set environment variables to configure the application:

```bash
export MAILGUN_DOMAIN=mg.yourdomain.com
export MAILGUN_API_KEY=key-abcdefghijkl
export MAILGUN_MAILING_LIST=newsletter@mg.yourdomain.com
export SUBSCRIBE_HTTP_PORT=8000
export SUBSCRIBE_REDIRECT_URL=http://yourdomain.com/newsletter/success
```

Run the application:

```
./mailgun-mailing-list-api
```

Add a form to your web page that submits the user's email address to the API:

```html
<form action="http://host:8000/subscribe">
	<input type="email" name="email" placeholder="Email your email address..." />
	<input type="button" value="Subscribe to the newsletter" />
</form>
```


Monitoring
----------

If you want to set up monitoring, to check the API is running at all times, the API has a URL for that, `/health-check`
