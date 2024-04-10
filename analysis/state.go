package analysis

import (
	"fmt"
	"golsp/lsp"
	"strings"
)

type State struct {
    // map of file names to contents
    Documents map[string]string
}

func NewState() State {
    return State{Documents: map[string]string{}}
}

func getDiagnosticsForFile(text string) []lsp.Diagnostic {
    diagnostics := []lsp.Diagnostic{}
    for row, line := range strings.Split(text, "\n") {
        if strings.Contains(line, "VS Code") {
            idx := strings.Index(line, "VS Code")
            diagnostics = append(diagnostics, lsp.Diagnostic{
            	Range:    LineRange(row, idx, idx+len("VS Code")),
            	Severity: 1,
            	Source:   "Common Sense",
            	Message:  "Forbidden editor detected",
            })
        }

        if strings.Contains(line, "Neovim") {
            idx := strings.Index(line, "Neovim")
            diagnostics = append(diagnostics, lsp.Diagnostic{
            	Range:    LineRange(row, idx, idx+len("Neovim")),
            	Severity: 3,
            	Source:   "Common Sense",
            	Message:  "Wise choice EZ",
            })
        }
    }

    return diagnostics
}

func (s *State) OpenDocument(uri, text string)  []lsp.Diagnostic {
    s.Documents[uri] = text

    return getDiagnosticsForFile(text)
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostic {
    s.Documents[uri] = text

    return getDiagnosticsForFile(text)
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

func (s *State) TextDocumentCodeAction(id int, uri string) lsp.TextDocumentCodeActionResponse {
    text := s.Documents[uri]

    actions := []lsp.CodeAction{}
    for row, line := range strings.Split(text, "\n") {
        idx := strings.Index(line, "VS Code")
        if idx >= 0 {
            replaceChange := map[string][]lsp.TextEdit{}
            replaceChange[uri] = []lsp.TextEdit {
                {
                    Range: LineRange(row, idx, idx+len("VS Code")),
                    NewText: "Neovim",
                },
            }

            actions = append(actions, lsp.CodeAction{
                Title: "Replace VS C*de with a superior editor",
                Edit: &lsp.WorkspaceEdit{Changes: replaceChange},
            })

            censorChange := map[string][]lsp.TextEdit{}
            censorChange[uri] = []lsp.TextEdit {
                {
                    Range: LineRange(row, idx, idx+len("VS Code")),
                    NewText: "VS C*de",
                },
            }

            actions = append(actions, lsp.CodeAction{
                Title: "Censor to VS C*de",
                Edit: &lsp.WorkspaceEdit{Changes: censorChange},
            })
        }
    }

    response := lsp.TextDocumentCodeActionResponse {
        Response: lsp.Response{
            RPC: "2.0",
            ID: &id,
        },
        Result: actions,
    }

    return response
}

func (s *State) TextDocumentCompletion(id int, uri string) lsp.CompletionResponse {

    // Ask your static analysis tools for good completions 
    items := []lsp.CompletionItem{
        {
            Label: "Neovim (BTW)",
            Detail: "Very cool editor",
            Documentation: "You are about to die of skill issues",
        },
    }

    response := lsp.CompletionResponse {
        Response: lsp.Response{
            RPC: "2.0",
            ID: &id,
        },
        Result: items,
    }

    return response
}

func LineRange(line, start, end int) lsp.Range {
    return lsp.Range {
        Start: lsp.Position{
            Line: line,
            Character: start,
        },
        End: lsp.Position{
            Line: line,
            Character: end,
        },
    }
}
