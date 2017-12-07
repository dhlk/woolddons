package main


const page string = `
{{define "uiclass"}}{{if eq .Updated .Install}}good{{else}}bad{{end}}{{end}}
{{define "mod"}}
				<td class="{{template "uiclass" .}}">{{.Install}}<br>{{.Updated}}</td>
				<td>{{.Version}}</td>
				<td><a href="{{.CurseURL}}">Curse</a></td>
				<td><a href="{{.ProjectURL}}">Project</a></td>
				<td><a href="{{.DownloadURL}}">Download</a></td>
{{end}}

{{define "page"}}<!DOCTYPE HTML>
<html>
	<head>
		<title>wooddons</title>
		<link rel="stylesheet" type="text/css" href="style.css">
	</head>
	<body>
		<iframe name="tf" frameborder="0" height="0" width="0"></iframe>
		<a href="/refresh" target="tf">update</a>
		<form method="post" action="/act" target="tf">
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
				<th>Curseforge</th>
				<th>Curse CDN</th>
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
.good {
	color: green;
}
.bad {
	color: red;
}
`

