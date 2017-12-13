package main

const page string = `
{{define "uiclass"}}{{if eq .Updated .Install}}good{{else}}bad{{end}}{{end}}
{{define "delete"}}<form method="post" action="/act"><input type="hidden" name="action" value="remove"><input type="hidden" name="addon" value="{{.}}"><input type="submit" value="[X]"></form>{{end}}
{{define "mod"}}
				<td class="{{template "uiclass" .}}">{{.Install}}<br>{{.Updated}}</td>
				<td>{{.Installed.FileName}}</td>
				<td><a href="{{.CurseURL}}">Curse</a></td>
				<td><a href="{{.DownloadURL .Newest}}">Download</a></td>
				<td>{{template "delete" .Addon}}</td>
{{end}}

{{define "page"}}<!DOCTYPE HTML>
<html>
	<head>
		<title>wooddons</title>

		<link href="style.css" rel="stylesheet" type="text/css">
	</head>
	<body>
		<a href="/refresh" target="tf">update</a>
		<form method="post" action="/act">
			<input type="hidden" name="action" value="add">
			<input type="text" name="addon">
			<input type="submit" value="Add">
		</form>
		<table border="1">
			<tr>
				<th>Name</th>
				<th>Installed<br>Updated</th>
				<th>Current Version</th>
				<th>Curse</th>
				<th>Curse CDN</th>
				<th>Delete</th>
			</tr>
		{{range $name, $mod :=  .}}
			<tr>
				<td>{{$name}}</td>
				{{template "mod" $mod}}
			</tr>
		{{end}}
		</table>
	</body>
</html>
{{end}}`

const style string = `
html {
	background-color: #111111;
	font-size: 16px;
	text-align: center;
	white-space: nowrap;
	color: white;
}

a {
	text-decoration: none;
	color: inherit;
}

table {
	margin-left: auto;
	margin-right: auto;
}

.good {
	color: green;
}
.bad {
	color: red;
}
`
