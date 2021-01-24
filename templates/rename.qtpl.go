// Code generated by qtc from "rename.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line templates/rename.qtpl:1
package templates

//line templates/rename.qtpl:1
import "net/http"

// This dialog is to be shown to a user when they try to rename a hypha.

//line templates/rename.qtpl:3
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line templates/rename.qtpl:3
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line templates/rename.qtpl:3
func StreamRenameAskHTML(qw422016 *qt422016.Writer, rq *http.Request, hyphaName string, isOld bool) {
//line templates/rename.qtpl:3
	qw422016.N().S(`
`)
//line templates/rename.qtpl:4
	streamnavHTML(qw422016, rq, hyphaName, "rename-ask")
//line templates/rename.qtpl:4
	qw422016.N().S(`
<main>
`)
//line templates/rename.qtpl:6
	if isOld {
//line templates/rename.qtpl:6
		qw422016.N().S(`	<section>
		<h1>Rename `)
//line templates/rename.qtpl:8
		qw422016.E().S(hyphaName)
//line templates/rename.qtpl:8
		qw422016.N().S(`</h1>
		<form action="/rename-confirm/`)
//line templates/rename.qtpl:9
		qw422016.E().S(hyphaName)
//line templates/rename.qtpl:9
		qw422016.N().S(`" method="post" enctype="multipart/form-data">
			<fieldset>
				<legend>New name</legend>
				<input type="text" value="`)
//line templates/rename.qtpl:12
		qw422016.E().S(hyphaName)
//line templates/rename.qtpl:12
		qw422016.N().S(`" required autofocus id="new-name" name="new-name"/>
			</fieldset>

			<fieldset>
				<legend>Settings</legend>
				<input type="checkbox" id="recursive" name="recursive" value="true" checked/>
				<label for="recursive">Keep subhyphae</label>
			</fieldset>

			<p>If you rename this hypha, all incoming links and all relative outcoming links will break. You will also lose all history for the new name. Rename carefully.</p>
			<input type="submit"/>
		</form>
	</section>
`)
//line templates/rename.qtpl:25
	} else {
//line templates/rename.qtpl:25
		qw422016.N().S(`	`)
//line templates/rename.qtpl:26
		streamcannotRenameDueToNonExistence(qw422016, hyphaName)
//line templates/rename.qtpl:26
		qw422016.N().S(`
`)
//line templates/rename.qtpl:27
	}
//line templates/rename.qtpl:27
	qw422016.N().S(`</main>
`)
//line templates/rename.qtpl:29
}

//line templates/rename.qtpl:29
func WriteRenameAskHTML(qq422016 qtio422016.Writer, rq *http.Request, hyphaName string, isOld bool) {
//line templates/rename.qtpl:29
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/rename.qtpl:29
	StreamRenameAskHTML(qw422016, rq, hyphaName, isOld)
//line templates/rename.qtpl:29
	qt422016.ReleaseWriter(qw422016)
//line templates/rename.qtpl:29
}

//line templates/rename.qtpl:29
func RenameAskHTML(rq *http.Request, hyphaName string, isOld bool) string {
//line templates/rename.qtpl:29
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/rename.qtpl:29
	WriteRenameAskHTML(qb422016, rq, hyphaName, isOld)
//line templates/rename.qtpl:29
	qs422016 := string(qb422016.B)
//line templates/rename.qtpl:29
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/rename.qtpl:29
	return qs422016
//line templates/rename.qtpl:29
}

//line templates/rename.qtpl:31
func streamcannotRenameDueToNonExistence(qw422016 *qt422016.Writer, hyphaName string) {
//line templates/rename.qtpl:31
	qw422016.N().S(`
	<section>
		<h1>Cannot rename `)
//line templates/rename.qtpl:33
	qw422016.E().S(hyphaName)
//line templates/rename.qtpl:33
	qw422016.N().S(`</h1>
		<p>This hypha does not exist.</p>
		<p><a href="/page/`)
//line templates/rename.qtpl:35
	qw422016.E().S(hyphaName)
//line templates/rename.qtpl:35
	qw422016.N().S(`">Go back</a></p>
	</section>
`)
//line templates/rename.qtpl:37
}

//line templates/rename.qtpl:37
func writecannotRenameDueToNonExistence(qq422016 qtio422016.Writer, hyphaName string) {
//line templates/rename.qtpl:37
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/rename.qtpl:37
	streamcannotRenameDueToNonExistence(qw422016, hyphaName)
//line templates/rename.qtpl:37
	qt422016.ReleaseWriter(qw422016)
//line templates/rename.qtpl:37
}

//line templates/rename.qtpl:37
func cannotRenameDueToNonExistence(hyphaName string) string {
//line templates/rename.qtpl:37
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/rename.qtpl:37
	writecannotRenameDueToNonExistence(qb422016, hyphaName)
//line templates/rename.qtpl:37
	qs422016 := string(qb422016.B)
//line templates/rename.qtpl:37
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/rename.qtpl:37
	return qs422016
//line templates/rename.qtpl:37
}