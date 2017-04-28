package resources

import "html/template"

var HomePageTemplate *template.Template

func init() {
	t := template.New("home page")

	t, _ = t.Parse(`
<!doctype html>
<html>
	<head>
		<style>
		 html, body { margin: 50px; }
		 div.postbody h1 { font-size: 12pt; }
		</style>
	</head>
	<body>
		<h1>Hello Templates</h1>
		<h2>Authors</h2>
		<ul>
		{{range .Authors}}
			<li>{{ .Email }}</li>
		{{end}}
		</ul>
		<h2>Posts</h2>
		{{ range .Posts }}
			<article id={{.Id}}>
				<h3>{{ .Slugline }} by {{ .Author }}</h3>
				<h4>{{ .DateCreated }}</h4>
				<div class='postbody'>
					{{ .Text }}
				</div>
			</article>
		{{ end }}
	</body>
</head>
`)

	HomePageTemplate = t
}
