// TEMPORARY AUTOGENERATED FILE: easyjson stub code to make the package
// compilable during generation.

package api

import (
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

func (EmptySuccessResponse) MarshalJSON() ([]byte, error)       { return nil, nil }
func (*EmptySuccessResponse) UnmarshalJSON([]byte) error        { return nil }
func (EmptySuccessResponse) MarshalEasyJSON(w *jwriter.Writer)  {}
func (*EmptySuccessResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {}

type EasyJSON_exporter_EmptySuccessResponse *EmptySuccessResponse
