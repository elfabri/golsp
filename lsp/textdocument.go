package lsp

type TextDocumentItem struct {
    /**
	 * The text document's URI.
	 */
     URI string `json:"uri"`

	/**
	 * The text document's language identifier.
	 */
	LanguageID string `json:"languageid"`

	/**
	 * The version number of this document (it will increase after each
	 * change, including undo/redo).
	 */
	Version int `json:"version"`

	/**
	 * The content of the opened text document.
	 */
    Text string `json:"text"`
}

type TextDocumentIdentifier struct {
    URI string `json:"uri"`
}

type VersionDocumentIdentifier struct {
    TextDocumentIdentifier
    Version int `json:"version"`
}
