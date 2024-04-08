package analysis

import (
	"fmt"
	"golsp/lsp"
)

type State struct {
    // map of file names to contents
    Documents map[string]string
}

func NewState() State {
    return State{Documents: map[string]string{}}
}

func (s *State) OpenDocument(uri, text string) {
    s.Documents[uri] = text
}

func (s *State) UpdateDocument(uri, text string) {
    s.Documents[uri] = text
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
    // in real life, this would look up the type in our type analysis code...
    // here we could have type analysis or something
    // we have direct information from the editor, no need to save the file 
    // for detecting changes

    doc := s.Documents[uri]

    return lsp.HoverResponse {
            Response: lsp.Response{
                RPC: "2.0",
                ID: &id,
            },
            Result: lsp.HoverResult{
                Contents: fmt.Sprintf("File: %s, Characters: %d", uri, len(doc)),
            },
        }
}

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
    // in real life, this would look up the definition
    // here we just move one line up

    return lsp.DefinitionResponse {
            Response: lsp.Response{
                RPC: "2.0",
                ID: &id,
            },
            Result: lsp.Location{
                URI: uri,
                Range: lsp.Range{
                    Start: lsp.Position{
                        Line: position.Line - 1,
                        Character: 0,
                    },
                    End: lsp.Position{
                        Line: position.Line - 1,
                        Character: 0,
                    },
                },
            },
        }
}
