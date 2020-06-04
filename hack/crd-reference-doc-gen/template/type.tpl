{{ define "type" }}
<h3 id="{{ anchorIDForType . }}">
    {{- .Name.Name }}
    {{ if eq .Kind "Alias" }}(<code>{{.Underlying}}</code> alias){{ end -}}
</h3>

{{ with (typeReferences .) }}
<p>
    (<em>Appears on:</em>
    {{- $prev := "" -}}
    {{- range . -}}
        {{- if $prev -}}, {{ end -}}
        {{ $prev = . }}
        <a href="{{ linkForType . }}">{{ typeDisplayName . }}</a>
    {{- end -}}
    )
</p>
{{ end }}

{{ with .CommentLines }}
{{ safe (renderComments .) }}
{{ end }}

{{ if .Members }}
<div class="md-typeset__scrollwrap">
    <div class="md-typeset__table">
        <table>
            <thead>
            <tr>
                <th>Field</th>
                <th>Description</th>
            </tr>
            </thead>
            <tbody>
            {{ if isExportedType . }}
                <tr>
                    <td>
                        <b>apiVersion</b><br>
                        string</td>
                    <td>
                        <b>{{ apiGroup . }}</b>
                    </td>
                </tr>
                <tr>
                    <td>
                        <b>kind</b><br>
                        string
                    </td>
                    <td>
                        <b>{{ .Name.Name }}</b>
                    </td>
                </tr>
            {{ end }}
            {{ template "members" . }}
            </tbody>
        </table>
    </div>
</div>
{{ end }}
{{ end }}
