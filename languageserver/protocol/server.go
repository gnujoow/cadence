/*
 * Cadence - The resource-oriented smart contract programming language
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package protocol

import (
	"github.com/onflow/cadence/languageserver/jsonrpc2"
)

// Server implements a layer between the JSONPRC2 server (jsonrpc2 package) and
// the business logic of the language server (server package).
//
// It handles unmarshalling inputs to the appropriate parameter types.
type Server struct {
	Handler        Handler
	conn           Conn
	jsonrpc2Server *jsonrpc2.Server
}

// Conn represents the connection to the language server client. It allows
// the language server to push various types of messages to the client.
// https://microsoft.github.io/language-server-protocol/specifications/specification-3-14
type Conn interface {
	Notify(method string, params interface{}) error
	ShowMessage(params *ShowMessageParams)
	LogMessage(params *LogMessageParams)
	PublishDiagnostics(params *PublishDiagnosticsParams) error
	RegisterCapability(params *RegistrationParams) error
}

type connection struct {
	jsonrpc2Server *jsonrpc2.Server
}

// ShowMessage displays a notification message to the client. It is visible to
// the user.
func (conn *connection) ShowMessage(params *ShowMessageParams) {
	_ = conn.Notify("window/showMessage", params)
}

// LogMessage logs a message to the Cadence terminal in VS Code. It isn't
// visible to the user unless they go looking for it.
func (conn *connection) LogMessage(params *LogMessageParams) {
	_ = conn.Notify("window/logMessage", params)
}

// PublishDiagnostics is used to report errors for a document, typically syntax
// or semantic errors in the code.
func (conn *connection) PublishDiagnostics(params *PublishDiagnosticsParams) error {
	return conn.Notify("textDocument/publishDiagnostics", params)
}

// Notify sends a notification to the client.
func (conn *connection) Notify(method string, params interface{}) error {
	return conn.jsonrpc2Server.Notify(method, params)
}

// RegisterCapability is used to dynamically inform the client that the server
// supports a particular API.
func (conn *connection) RegisterCapability(params *RegistrationParams) error {
	return conn.jsonrpc2Server.Call("client/registerCapability", params)
}

// Handler defines the subset of the Language Server Protocol we support.
type Handler interface {
	Initialize(conn Conn, params *InitializeParams) (*InitializeResult, error)
	DidOpenTextDocument(conn Conn, params *DidOpenTextDocumentParams) error
	DidChangeTextDocument(conn Conn, params *DidChangeTextDocumentParams) error
	Hover(conn Conn, params *TextDocumentPositionParams) (*Hover, error)
	Definition(conn Conn, params *TextDocumentPositionParams) (*Location, error)
	SignatureHelp(conn Conn, params *TextDocumentPositionParams) (*SignatureHelp, error)
	DocumentHighlight(conn Conn, params *TextDocumentPositionParams) ([]*DocumentHighlight, error)
	Rename(conn Conn, params *RenameParams) (*WorkspaceEdit, error)
	CodeAction(conn Conn, params *CodeActionParams) ([]*CodeAction, error)
	CodeLens(conn Conn, params *CodeLensParams) ([]*CodeLens, error)
	Completion(conn Conn, params *CompletionParams) ([]*CompletionItem, error)
	ResolveCompletionItem(conn Conn, item *CompletionItem) (*CompletionItem, error)
	ExecuteCommand(conn Conn, params *ExecuteCommandParams) (interface{}, error)
	DocumentSymbol(conn Conn, params *DocumentSymbolParams) ([]*DocumentSymbol, error)
	Shutdown(conn Conn) error
	Exit(conn Conn) error
}

func NewServer(handler Handler) *Server {
	jsonrpc2Server := jsonrpc2.NewServer()

	conn := &connection{
		jsonrpc2Server: jsonrpc2Server,
	}

	server := &Server{
		Handler:        handler,
		conn:           conn,
		jsonrpc2Server: jsonrpc2Server,
	}

	jsonrpc2Server.Methods["initialize"] =
		server.handleInitialize

	jsonrpc2Server.Methods["textDocument/didOpen"] =
		server.handleDidOpenTextDocument

	jsonrpc2Server.Methods["textDocument/didChange"] =
		server.handleDidChangeTextDocument

	jsonrpc2Server.Methods["textDocument/hover"] =
		server.handleHover

	jsonrpc2Server.Methods["textDocument/definition"] =
		server.handleDefinition

	jsonrpc2Server.Methods["textDocument/signatureHelp"] =
		server.handleSignatureHelp

	jsonrpc2Server.Methods["textDocument/codeLens"] =
		server.handleCodeLens

	jsonrpc2Server.Methods["textDocument/documentHighlight"] =
		server.handleDocumentHighlight

	jsonrpc2Server.Methods["textDocument/rename"] =
		server.handleRename

	jsonrpc2Server.Methods["textDocument/codeAction"] =
		server.handleCodeAction

	jsonrpc2Server.Methods["textDocument/completion"] =
		server.handleCompletion

	jsonrpc2Server.Methods["completionItem/resolve"] =
		server.handleCompletionItemResolve

	jsonrpc2Server.Methods["workspace/executeCommand"] =
		server.handleExecuteCommand

	jsonrpc2Server.Methods["textDocument/documentSymbol"] =
		server.handleDocumentSymbol

	jsonrpc2Server.Methods["shutdown"] =
		server.handleShutdown

	jsonrpc2Server.Methods["exit"] =
		server.handleExit

	return server
}

func (s *Server) Start(stream jsonrpc2.ObjectStream) <-chan struct{} {
	return s.jsonrpc2Server.Start(stream)
}

func (s *Server) Stop() error {
	return s.jsonrpc2Server.Stop()
}
