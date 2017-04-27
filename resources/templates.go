package resources

import "html/template"

var HomePageTemplate *template.Template

func init() {
	t := template.New("home page")

	t, _ = t.Parse(`
<h1>Hello Templates</h1>
<h2>Authors</h2>
<ul>
{{range .Authors}}
	<li>{{ .Email }}</li>
{{end}}
</ul>
<h2>Posts</h2>
{{ range .Posts }}
	<article>
		<p>{{ .Slugline }} by {{ .Author }}</p>
		<p>{{ .DateCreated }}</p>
		<div>
			{{ .Text }}
		</div>
	</article>
{{ end }}
`)

	HomePageTemplate = t
}
