package main

import (
	"text/template"
)

type rootTemplateArgs struct {
	ListingEndpoint string
	Message         string
}

var rootTemplate = template.Must(template.New("root").Parse(`
<html>
<head>
       <title>Upload file</title>
</head>
<body>
{{ .Message }}
<form enctype="multipart/form-data" action="/upload" method="POST">
	upload file:
    <input type="file" name="source" />
    <input type="submit" value="upload" />
</form>
<br />
browse dir: <a href="files/">{{ .ListingEndpoint }}</a>
</body>
</html>
`))
