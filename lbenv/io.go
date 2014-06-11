package lbenv

import (
	"encoding/xml"
	"fmt"
	"io"
)

func Decode(r io.Reader) ([]Action, error) {
	var err error

	actions := make([]Action, 0)
	caller := ""
	if f, ok := r.(interface {
		Name() string
	}); ok {
		caller = f.Name()
	}

	dec := xml.NewDecoder(r)
	var tok xml.Token
	for {
		tok, err = dec.Token()
		if err != nil {
			break
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			var action Action
			switch tok.Name.Local {
			case "config":
				continue

			case "declare":
				var a DeclareVar
				action = &a
				for _, attr := range tok.Attr {
					switch attr.Name.Local {
					case "local":
						switch attr.Value {
						case "true":
							a.Local = true
						case "false":
							a.Local = false
						}

					case "type":
						switch attr.Value {
						case "list":
							a.Type = VarList
						case "scalar":
							a.Type = VarScalar
						}

					case "variable":
						a.Name = attr.Value
					}
				}

			case "unset":
				var a UnsetVar
				action = &a
				for _, attr := range tok.Attr {
					switch attr.Name.Local {
					case "variable":
						a.Name = attr.Value
					}
				}

			case "set":
				var a SetVar
				action = &a
				for _, attr := range tok.Attr {
					switch attr.Name.Local {
					case "variable":
						a.Name = attr.Value
					}
				}
				var vtok xml.Token
				vtok, err = dec.Token()
				if err != nil {
					return nil, err
				}
				a.Value = string(vtok.(xml.CharData))

			case "prepend":
				var a PrependVar
				action = &a
				for _, attr := range tok.Attr {
					switch attr.Name.Local {
					case "variable":
						a.Name = attr.Value
					}
				}
				var vtok xml.Token
				vtok, err = dec.Token()
				if err != nil {
					return nil, err
				}
				a.Value = string(vtok.(xml.CharData))

			case "append":
				var a AppendVar
				action = &a
				for _, attr := range tok.Attr {
					switch attr.Name.Local {
					case "variable":
						a.Name = attr.Value
					}
				}
				var vtok xml.Token
				vtok, err = dec.Token()
				if err != nil {
					return nil, err
				}
				a.Value = string(vtok.(xml.CharData))

			case "include":
				var a Include
				action = &a
				for _, attr := range tok.Attr {
					switch attr.Name.Local {
					case "hints":
						a.Hints = attr.Value
					}
				}
				var vtok xml.Token
				vtok, err = dec.Token()
				if err != nil {
					return nil, err
				}
				a.File = string(vtok.(xml.CharData))
				a.Caller = caller

			default:
				fmt.Printf(">>> name=%q\n", tok.Name)
				for _, attr := range tok.Attr {
					fmt.Printf("    attr=%q, value=%q\n", attr.Name, attr.Value)
				}
				panic(fmt.Errorf("unknown action %q", tok.Name.Local))
			}
			actions = append(actions, action)

			var endtok xml.Token
			endtok, err = dec.Token()
			if err != nil {
				return nil, err
			}

			_ = endtok.(xml.EndElement)

		case xml.CharData:
			//fmt.Printf(">>> chardata=%q\n", string(tok))

		case xml.ProcInst:
			// noop

		case xml.EndElement:
			// noop

		default:
			fmt.Printf("--- %v (%T)\n", tok, tok)
			panic(fmt.Errorf("unhandled toke: %[1]v %[1]T", tok))
		}
	}

	return actions, err
}

func Encode(w io.Writer, actions []Action) error {
	var err error
	_, err = fmt.Fprintf(
		w,
		`<?xml version="1.0" ?>
<env:config xmlns:env="EnvSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="EnvSchema ./EnvSchema.xsd ">
`)
	if err != nil {
		return err
	}

	for _, action := range actions {
		switch v := action.(type) {
		case *DeclareVar:
			vtype := ""
			switch v.Type {
			case VarScalar:
				vtype = "scalar"
			case VarList:
				vtype = "list"
			default:
				panic(fmt.Errorf("unknown variable type %[1]v (%[1]T)", v.Type, v.Type))
			}
			_, err = fmt.Fprintf(
				w, "<env:declare local=%q type=%q variable=%q/>\n",
				v.Local, vtype, v.Name,
			)
		case *SetVar:
			_, err = fmt.Fprintf(w, "<env:set variable=%q>%s</env:set>\n", v.Name, v.Value)
		case *UnsetVar:
			_, err = fmt.Fprintf(w, "<env:unset variable=%q/>\n", v.Name)
		case *AppendVar:
			_, err = fmt.Fprintf(w, "<env:append variable=%q>%s</env:append>\n",
				v.Name, v.Value,
			)
		case *PrependVar:
			_, err = fmt.Fprintf(w, "<env:prepend variable=%q>%s</env:prepend>\n",
				v.Name, v.Value,
			)
		case *Include:
			_, err = fmt.Fprintf(w, "<env:include hints=%q>%s</env:include>\n",
				v.Hints, v.File,
			)
		default:
			panic(fmt.Errorf("unknown Action type: %[1]v (type=%[1]T)", v))
		}
	}

	_, err = fmt.Fprintf(w, `</env:config>`)
	return err
}
