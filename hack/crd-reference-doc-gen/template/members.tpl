{{ define "members" }}
    {{ range .Members }}
        {{ if not (hiddenMember .)}}
            <tr>
                <td>
                    <b>{{ fieldName . }}</b><br>
                    <em>
                        {{ if linkForType .Type }}
                            <a href="{{ linkForType .Type }}">
                                {{ typeDisplayName .Type }}
                            </a>
                        {{ else }}
                            {{ typeDisplayName .Type }}
                        {{ end }}
                    </em>
                </td>
                <td>
                    {{ if fieldEmbedded . }}
                        <p>
                            (Members of <b>{{ fieldName . }}</b> are embedded into this type.)
                        </p>
                    {{ end}}

                    {{ if isOptionalMember .}}
                        <em>(Optional)</em>
                    {{ end }}

                    {{ safe (renderComments .CommentLines) }}

                    {{ if and (eq (.Type.Name.Name) "ObjectMeta") }}
                        Refer to the Kubernetes API documentation for the fields of the
                        <b>metadata</b> field.
                    {{ end }}

                    {{ if or (eq (fieldName .) "spec") (eq (fieldName .) "status") (eq (fieldName .) "chart") }}
                        <table>
                            {{ template "members" .Type }}
                        </table>
                    {{ end }}
                </td>
            </tr>
        {{ end }}
    {{ end }}
{{ end }}
