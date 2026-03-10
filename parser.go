package gojsonparser

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

type Parser struct {
	scanner scanner.Scanner
	peeked  *rune
}

func ParseJson(src string) (any, error) {
	var s scanner.Scanner
	s.Init(strings.NewReader(src))
	s.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanStrings
	p := &Parser{scanner: s}
	return p.parseValue()
}

func (p *Parser) advance() rune {
	if p.peeked != nil {
		t := *p.peeked
		p.peeked = nil
		return t
	}
	return p.scanner.Scan()
}

func (p *Parser) peek() rune {
	if p.peeked == nil {
		t := p.scanner.Scan()
		p.peeked = &t
	}
	return *p.peeked
}

func (p *Parser) expect(expected rune) error {
	tok := p.advance()
	if tok != expected {
		return fmt.Errorf("expected '%c', got '%c'", expected, tok)
	}
	return nil
}

func (p *Parser) parseValue() (any, error) {
	switch p.peek() {
	case scanner.EOF:
		return nil, nil
	case scanner.Ident:
		return p.parseBooleanOrNull()
	case scanner.Int:
		return p.parseNumber()
	case scanner.String:
		return p.parseString()
	case '[':
		return p.parseArray()
	case '{':
		return p.parseObject()
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.scanner.TokenText())
	}
}

func (p *Parser) parseBooleanOrNull() (any, error) {
	if err := p.expect(scanner.Ident); err != nil {
		return nil, err
	}

	iden := p.scanner.TokenText()
	switch iden {
	case "true":
		return true, nil
	case "false":
		return false, nil
	case "null":
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected identifier: %s", iden)
	}
}

func (p *Parser) parseNumber() (int, error) {
	if err := p.expect(scanner.Int); err != nil {
		return 0, err
	}

	value, err := strconv.Atoi(p.scanner.TokenText())
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (p *Parser) parseString() (string, error) {
	if err := p.expect(scanner.String); err != nil {
		return "", err
	}

	unquoted, err := strconv.Unquote(p.scanner.TokenText())
	if err != nil {
		return "", err
	}
	return unquoted, nil
}

func (p *Parser) parseArray() ([]any, error) {
	if err := p.expect('['); err != nil {
		return nil, err
	}

	var arr []any
	
	if p.peek() != ']' {
		for {
			value, err := p.parseValue()
			if err != nil {
				return nil, err
			}
			arr = append(arr, value)

			if p.peek() == ',' {
				p.advance()
			} else {
				break
			}
		}
	}

	if err := p.expect(']'); err != nil {
		return nil, err
	}
	return arr, nil
}

func (p *Parser) parseObject() (map[string]any, error) {
	if err := p.expect('{'); err != nil {
		return nil, err
	}

	obj := make(map[string]any)

	if p.peek() != '}' {
		for {
			key, err := p.parseString()
			if err != nil {
				return nil, err
			}

			if err := p.expect(':'); err != nil {
				return nil, err
			}

			value, err := p.parseValue()
			if err != nil {
				return nil, err
			}

			obj[key] = value

			if p.peek() == ',' {
				p.advance()
			} else {
				break
			}
		}
	}
	if err := p.expect('}'); err != nil {
		return nil, err
	}

	return obj, nil
}