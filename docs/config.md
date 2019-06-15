---
layout: default
title: gosaas configuration file
---

[back to main content](index.md)

# Configuration file

You may use a configuration file to enable and set options and behaviors of gosaas.

**File name**: gosaas.json | **Location**: at the root of your package

### Example

```json
{
	"emailFrom": "you@yourdomain.com",
	"emailFromName": "Your Name",
	"emailProvider": "amazonses",
	"stripeKey": "your-stripe-key",
	"signupTemplate": "",
	"sendEmailValidation": false,
	"signupSuccessRedirect": "",
	"signupErrorRedirect": "",
	"signinTemplate": "",
	"signinSuccessRedirect": "",
	"signinErrorRedirect": "/users/login"
}
```

### Options



| name									| type			| description																				|
| ---------------------:|:---------:| ------------ 																			|
| emailFrom							| `string`	| Email used as the `from` of transactional emails 	|
| emailFromName					| `string`	| From name field																		|
| emailProvider					| `string`	| Email provider (only amazonses for now)						|
| stripeKey							| `string`	| Your Stripe key																		|
| signupTemplate				| `string`	| Template to use for sign up html page							|
| sendEmailValidation		| `bool`		| Should it sends a validation email								|
| signupSuccessRedirect	| `string`	| When using HTML, URL to redirect after signup			|
| signupErrorRedirect		| `string`	| When using HTML, URL when there's an error				|

*Sign in options are same as signup so they are not present in the table.*