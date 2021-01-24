// Code generated by qtc from "http_stuff.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line templates/http_stuff.qtpl:1
package templates

//line templates/http_stuff.qtpl:1
import "github.com/bouncepaw/mycorrhiza/util"

//line templates/http_stuff.qtpl:2
import "github.com/bouncepaw/mycorrhiza/user"

//line templates/http_stuff.qtpl:4
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line templates/http_stuff.qtpl:4
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line templates/http_stuff.qtpl:4
func StreamBaseHTML(qw422016 *qt422016.Writer, title, body string, u *user.User, headElements ...string) {
//line templates/http_stuff.qtpl:4
	qw422016.N().S(`
<!doctype html>
<html>
	<head>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="stylesheet" type="text/css" href="/static/common.css">
		<title>`)
//line templates/http_stuff.qtpl:10
	qw422016.E().S(title)
//line templates/http_stuff.qtpl:10
	qw422016.N().S(`</title>
		`)
//line templates/http_stuff.qtpl:11
	for _, el := range headElements {
//line templates/http_stuff.qtpl:11
		qw422016.N().S(el)
//line templates/http_stuff.qtpl:11
	}
//line templates/http_stuff.qtpl:11
	qw422016.N().S(`
	</head>
	<body>
		<header>
			<nav class="header-links">
				<ul class="header-links__list">
`)
//line templates/http_stuff.qtpl:17
	for _, link := range util.HeaderLinks {
//line templates/http_stuff.qtpl:17
		qw422016.N().S(`					<li class="header-links__entry"><a class="header-links__link" href="`)
//line templates/http_stuff.qtpl:18
		qw422016.E().S(link.Href)
//line templates/http_stuff.qtpl:18
		qw422016.N().S(`">`)
//line templates/http_stuff.qtpl:18
		qw422016.E().S(link.Display)
//line templates/http_stuff.qtpl:18
		qw422016.N().S(`</a></li>
`)
//line templates/http_stuff.qtpl:19
	}
//line templates/http_stuff.qtpl:19
	qw422016.N().S(`					`)
//line templates/http_stuff.qtpl:20
	qw422016.N().S(userMenuHTML(u))
//line templates/http_stuff.qtpl:20
	qw422016.N().S(`
				</ul>
			</nav>
		</header>
		`)
//line templates/http_stuff.qtpl:24
	qw422016.N().S(body)
//line templates/http_stuff.qtpl:24
	qw422016.N().S(`
	</body>
</html>
`)
//line templates/http_stuff.qtpl:27
}

//line templates/http_stuff.qtpl:27
func WriteBaseHTML(qq422016 qtio422016.Writer, title, body string, u *user.User, headElements ...string) {
//line templates/http_stuff.qtpl:27
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/http_stuff.qtpl:27
	StreamBaseHTML(qw422016, title, body, u, headElements...)
//line templates/http_stuff.qtpl:27
	qt422016.ReleaseWriter(qw422016)
//line templates/http_stuff.qtpl:27
}

//line templates/http_stuff.qtpl:27
func BaseHTML(title, body string, u *user.User, headElements ...string) string {
//line templates/http_stuff.qtpl:27
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/http_stuff.qtpl:27
	WriteBaseHTML(qb422016, title, body, u, headElements...)
//line templates/http_stuff.qtpl:27
	qs422016 := string(qb422016.B)
//line templates/http_stuff.qtpl:27
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/http_stuff.qtpl:27
	return qs422016
//line templates/http_stuff.qtpl:27
}

//line templates/http_stuff.qtpl:29
func StreamHyphaListHTML(qw422016 *qt422016.Writer, tbody string, pageCount int) {
//line templates/http_stuff.qtpl:29
	qw422016.N().S(`
<main>
	<h1>List of hyphae</h1>
	<p>This wiki has `)
//line templates/http_stuff.qtpl:32
	qw422016.N().D(pageCount)
//line templates/http_stuff.qtpl:32
	qw422016.N().S(` hyphae.</p>
	<table>
		<thead>
			<tr>
				<th>Full name</th>
				<th>Binary part type</th>
			</tr>
		</thead>
		<tbody>
			`)
//line templates/http_stuff.qtpl:41
	qw422016.N().S(tbody)
//line templates/http_stuff.qtpl:41
	qw422016.N().S(`
		</tbody>
	</table>
</main>
`)
//line templates/http_stuff.qtpl:45
}

//line templates/http_stuff.qtpl:45
func WriteHyphaListHTML(qq422016 qtio422016.Writer, tbody string, pageCount int) {
//line templates/http_stuff.qtpl:45
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/http_stuff.qtpl:45
	StreamHyphaListHTML(qw422016, tbody, pageCount)
//line templates/http_stuff.qtpl:45
	qt422016.ReleaseWriter(qw422016)
//line templates/http_stuff.qtpl:45
}

//line templates/http_stuff.qtpl:45
func HyphaListHTML(tbody string, pageCount int) string {
//line templates/http_stuff.qtpl:45
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/http_stuff.qtpl:45
	WriteHyphaListHTML(qb422016, tbody, pageCount)
//line templates/http_stuff.qtpl:45
	qs422016 := string(qb422016.B)
//line templates/http_stuff.qtpl:45
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/http_stuff.qtpl:45
	return qs422016
//line templates/http_stuff.qtpl:45
}

//line templates/http_stuff.qtpl:47
func StreamHyphaListRowHTML(qw422016 *qt422016.Writer, hyphaName, binaryMime string, binaryPresent bool) {
//line templates/http_stuff.qtpl:47
	qw422016.N().S(`
			<tr>
				<td><a href="/page/`)
//line templates/http_stuff.qtpl:49
	qw422016.E().S(hyphaName)
//line templates/http_stuff.qtpl:49
	qw422016.N().S(`">`)
//line templates/http_stuff.qtpl:49
	qw422016.E().S(hyphaName)
//line templates/http_stuff.qtpl:49
	qw422016.N().S(`</a></td>
			`)
//line templates/http_stuff.qtpl:50
	if binaryPresent {
//line templates/http_stuff.qtpl:50
		qw422016.N().S(`
				<td>`)
//line templates/http_stuff.qtpl:51
		qw422016.E().S(binaryMime)
//line templates/http_stuff.qtpl:51
		qw422016.N().S(`</td>
			`)
//line templates/http_stuff.qtpl:52
	} else {
//line templates/http_stuff.qtpl:52
		qw422016.N().S(`
				<td></td>
			`)
//line templates/http_stuff.qtpl:54
	}
//line templates/http_stuff.qtpl:54
	qw422016.N().S(`
			</tr>
`)
//line templates/http_stuff.qtpl:56
}

//line templates/http_stuff.qtpl:56
func WriteHyphaListRowHTML(qq422016 qtio422016.Writer, hyphaName, binaryMime string, binaryPresent bool) {
//line templates/http_stuff.qtpl:56
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/http_stuff.qtpl:56
	StreamHyphaListRowHTML(qw422016, hyphaName, binaryMime, binaryPresent)
//line templates/http_stuff.qtpl:56
	qt422016.ReleaseWriter(qw422016)
//line templates/http_stuff.qtpl:56
}

//line templates/http_stuff.qtpl:56
func HyphaListRowHTML(hyphaName, binaryMime string, binaryPresent bool) string {
//line templates/http_stuff.qtpl:56
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/http_stuff.qtpl:56
	WriteHyphaListRowHTML(qb422016, hyphaName, binaryMime, binaryPresent)
//line templates/http_stuff.qtpl:56
	qs422016 := string(qb422016.B)
//line templates/http_stuff.qtpl:56
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/http_stuff.qtpl:56
	return qs422016
//line templates/http_stuff.qtpl:56
}

//line templates/http_stuff.qtpl:58
func StreamAboutHTML(qw422016 *qt422016.Writer) {
//line templates/http_stuff.qtpl:58
	qw422016.N().S(`
<main>
	<section>
		<h1>About `)
//line templates/http_stuff.qtpl:61
	qw422016.E().S(util.SiteName)
//line templates/http_stuff.qtpl:61
	qw422016.N().S(`</h1>
		<ul>
			<li><b><a href="https://mycorrhiza.lesarbr.es">MycorrhizaWiki</a> version:</b> β 0.12 indev</li>
`)
//line templates/http_stuff.qtpl:64
	if user.AuthUsed {
//line templates/http_stuff.qtpl:64
		qw422016.N().S(`			<li><b>User count:</b> `)
//line templates/http_stuff.qtpl:65
		qw422016.N().D(user.Count())
//line templates/http_stuff.qtpl:65
		qw422016.N().S(`</li>
			<li><b>Home page:</b> <a href="/">`)
//line templates/http_stuff.qtpl:66
		qw422016.E().S(util.HomePage)
//line templates/http_stuff.qtpl:66
		qw422016.N().S(`</a></li>
			<li><b>Administrators:</b>`)
//line templates/http_stuff.qtpl:67
		for i, username := range user.ListUsersWithGroup("admin") {
//line templates/http_stuff.qtpl:68
			if i > 0 {
//line templates/http_stuff.qtpl:68
				qw422016.N().S(`<span aria-hidden="true">, </span>
`)
//line templates/http_stuff.qtpl:69
			}
//line templates/http_stuff.qtpl:69
			qw422016.N().S(`				<a href="/page/`)
//line templates/http_stuff.qtpl:70
			qw422016.E().S(util.UserHypha)
//line templates/http_stuff.qtpl:70
			qw422016.N().S(`/`)
//line templates/http_stuff.qtpl:70
			qw422016.E().S(username)
//line templates/http_stuff.qtpl:70
			qw422016.N().S(`">`)
//line templates/http_stuff.qtpl:70
			qw422016.E().S(username)
//line templates/http_stuff.qtpl:70
			qw422016.N().S(`</a>`)
//line templates/http_stuff.qtpl:70
		}
//line templates/http_stuff.qtpl:70
		qw422016.N().S(`</li>
`)
//line templates/http_stuff.qtpl:71
	} else {
//line templates/http_stuff.qtpl:71
		qw422016.N().S(`			<li>This wiki does not use authorization</li>
`)
//line templates/http_stuff.qtpl:73
	}
//line templates/http_stuff.qtpl:73
	qw422016.N().S(`		</ul>
		<p>See <a href="/list">/list</a> for information about hyphae on this wiki.</p>
	</section>
</main>
`)
//line templates/http_stuff.qtpl:78
}

//line templates/http_stuff.qtpl:78
func WriteAboutHTML(qq422016 qtio422016.Writer) {
//line templates/http_stuff.qtpl:78
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/http_stuff.qtpl:78
	StreamAboutHTML(qw422016)
//line templates/http_stuff.qtpl:78
	qt422016.ReleaseWriter(qw422016)
//line templates/http_stuff.qtpl:78
}

//line templates/http_stuff.qtpl:78
func AboutHTML() string {
//line templates/http_stuff.qtpl:78
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/http_stuff.qtpl:78
	WriteAboutHTML(qb422016)
//line templates/http_stuff.qtpl:78
	qs422016 := string(qb422016.B)
//line templates/http_stuff.qtpl:78
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/http_stuff.qtpl:78
	return qs422016
//line templates/http_stuff.qtpl:78
}