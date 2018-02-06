package brdgme

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Parser interface {
	Parse(input string, names []string) (Output, *ParseError)
	Expected(names []string) []string
	ToSpec() Spec
}

type Output struct {
	Value     interface{}
	Consumed  string
	Remaining string
}

type ParseError struct {
	Message  string
	Expected []string
	Offset   uint
}

func (e ParseError) Error() string {
	output := &bytes.Buffer{}
	if e.Message != "" {
		output.WriteString(e.Message)
		if len(e.Expected) > 0 {
			output.WriteString(", ")
		}
	}
	if len(e.Expected) > 0 {
		output.WriteString("expected ")
	}
	return output.String()
}

type Spec struct {
	Int    *Int    `json:",omitempty"`
	Token  *Token  `json:",omitempty"`
	Enum   *Enum   `json:",omitempty"`
	OneOf  *OneOf  `json:",omitempty"`
	Chain  *Chain  `json:",omitempty"`
	Many   *Many   `json:",omitempty"`
	Opt    *Opt    `json:",omitempty"`
	Doc    *Doc    `json:",omitempty"`
	Player *Player `json:",omitempty"`
	Space  *Space  `json:",omitempty"`
}

var _ Parser = Spec{}

func (s Spec) ToSpec() Spec {
	return s
}

func (s Spec) Parsers() []Parser {
	parsers := []Parser{}
	if s.Int != nil {
		parsers = append(parsers, s.Int)
	}
	if s.Token != nil {
		parsers = append(parsers, s.Token)
	}
	if s.Enum != nil {
		parsers = append(parsers, s.Enum)
	}
	if s.OneOf != nil {
		parsers = append(parsers, s.OneOf)
	}
	if s.Chain != nil {
		parsers = append(parsers, s.Chain)
	}
	if s.Many != nil {
		parsers = append(parsers, s.Many)
	}
	if s.Opt != nil {
		parsers = append(parsers, s.Opt)
	}
	if s.Doc != nil {
		parsers = append(parsers, s.Doc)
	}
	if s.Player != nil {
		parsers = append(parsers, s.Player)
	}
	if s.Space != nil {
		parsers = append(parsers, s.Space)
	}
	return parsers
}

func (s Spec) Parse(input string, names []string) (Output, *ParseError) {
	for _, p := range s.Parsers() {
		return p.Parse(input, names)
	}
	return Output{}, &ParseError{
		Message: "there are no available parsers",
	}
}

func (s Spec) Expected(names []string) []string {
	expected := []string{}
	for _, p := range s.Parsers() {
		expected = append(expected, p.Expected(names)...)
	}
	return expected
}

type Int struct {
	Min *int `json:"min,omitempty"`
	Max *int `json:"max,omitempty"`
}

var _ Parser = Int{}

func (i Int) ToSpec() Spec {
	return Spec{
		Int: &i,
	}
}

var IntRegexp = regexp.MustCompile("^-?[0-9]+")

func (i Int) Parse(input string, names []string) (Output, *ParseError) {
	match := IntRegexp.FindString(input)
	parsed, err := strconv.Atoi(match)
	if err != nil {
		return Output{}, &ParseError{
			Expected: i.Expected(names),
		}
	}
	if i.Min != nil && int(*i.Min) > parsed {
		return Output{}, &ParseError{
			Message:  fmt.Sprintf("%d is too low", parsed),
			Expected: i.Expected(names),
		}
	}
	if i.Max != nil && int(*i.Max) < parsed {
		return Output{}, &ParseError{
			Message:  fmt.Sprintf("%d is too high", parsed),
			Expected: i.Expected(names),
		}
	}
	return Output{
		Value:     parsed,
		Consumed:  match,
		Remaining: input[len(match):],
	}, nil
}

func (i Int) ExpectedOutput() string {
	switch {
	case i.Min != nil && i.Max != nil:
		return fmt.Sprintf("number between %d and %d", *i.Min, *i.Max)
	case i.Min != nil && i.Max == nil:
		return fmt.Sprintf("number %d or higher", *i.Min)
	case i.Max != nil && i.Min == nil:
		return fmt.Sprintf("number %d or lower", *i.Max)
	default:
		return "number"
	}
}

func (i Int) Expected(names []string) []string {
	return []string{i.ExpectedOutput()}
}

type Token string

var _ Parser = Token("")

func (t Token) ToSpec() Spec {
	return Spec{
		Token: &t,
	}
}

func (t Token) Parse(input string, names []string) (Output, *ParseError) {
	tLen := len(t)
	if strings.HasPrefix(strings.ToLower(input[:tLen]), strings.ToLower(string(t))) {
		return Output{
			Value:     t,
			Consumed:  input[:tLen],
			Remaining: input[tLen:],
		}, nil
	}
	return Output{}, &ParseError{
		Expected: t.Expected(names),
	}
}

func (t Token) Expected(names []string) []string {
	return []string{string(t)}
}

type Enum struct {
	Values []string `json:"values"`
	Exact  bool     `json:"exact"`
}

var _ Parser = Enum{}

func (e Enum) ToSpec() Spec {
	return Spec{
		Enum: &e,
	}
}

func sharedPrefix(s1, s2 string) int {
	until := len(s1)
	if s2Len := len(s2); s2Len < until {
		until = s2Len
	}

	for i := 0; i < until; i++ {
		if s1[i] != s2[i] {
			return i
		}
	}

	return until
}

func commaList(items []string, conj string) string {
	switch len(items) {
	case 0:
		return ""
	case 1:
		return items[0]
	case 2:
		return fmt.Sprintf("%s %s %s", items[0], conj, items[1])
	default:
		return fmt.Sprintf("%s, %s", items[0], commaList(items[1:], conj))
	}
}

func commaListAnd(items []string) string {
	return commaList(items, "and")
}

func commaListOr(items []string) string {
	return commaList(items, "or")
}

func (e Enum) Parse(input string, names []string) (Output, *ParseError) {
	inputLower := strings.ToLower(input)
	matched := []string{}
	matchLen := 0
	fullMatch := false

	for _, v := range e.Values {
		vLower := strings.ToLower(v)
		vLen := len(vLower)

		matching := sharedPrefix(inputLower, vLower)
		if e.Exact && matching < vLen {
			continue
		}

		if matching > 0 && matching >= matchLen && (!fullMatch || matching == vLen) {
			if matching == vLen {
				fullMatch = true
			}
			if matching > matchLen {
				matched = []string{v}
				matchLen = matching
			} else {
				matched = append(matched, v)
			}
		}
	}

	switch len(matched) {
	case 1:
		return Output{
			Value:     matched[0],
			Consumed:  input[:matchLen],
			Remaining: input[matchLen:],
		}, nil
	case 0:
		return Output{}, &ParseError{
			Expected: e.Expected(names),
		}
	default:
		return Output{}, &ParseError{
			Message: fmt.Sprintf(
				"matched %s, more input is required to uniquely match one",
				commaListAnd(matched),
			),
			Expected: e.Expected(names),
		}
	}
}

func (e Enum) Expected(names []string) []string {
	return e.Values
}

type OneOf []Spec

var _ Parser = OneOf{}

func (o OneOf) ToSpec() Spec {
	return Spec{
		OneOf: &o,
	}
}

func (o OneOf) Parse(input string, names []string) (Output, *ParseError) {
	errors := []ParseError{}
	errorConsumed := uint(0)

	for _, p := range o {
		output, err := p.Parse(input, names)
		if err == nil {
			return output, nil
		}
		if err.Offset > errorConsumed {
			errors = []ParseError{*err}
			errorConsumed = err.Offset
		} else if err.Offset == errorConsumed {
			errors = append(errors, *err)
		}
	}

	messages := []string{}
	expected := []string{}
	for _, e := range errors {
		expected = append(expected, e.Expected...)
		if e.Message != "" {
			messages = append(messages, e.Message)
		}
	}
	return Output{}, &ParseError{
		Message:  commaListOr(messages),
		Expected: expected,
		Offset:   errorConsumed,
	}
}

func (o OneOf) Expected(names []string) []string {
	expected := []string{}
	for _, spec := range o {
		expected = append(expected, spec.Expected(names)...)
	}
	return expected
}

type Chain []Spec

var _ Parser = Chain{}

func (c Chain) ToSpec() Spec {
	return Spec{
		Chain: &c,
	}
}

func (c Chain) Expected(names []string) []string {
	if len(c) == 0 {
		return []string{}
	}
	return c[0].Expected(names)
}

func (c Chain) Parse(input string, names []string) (Output, *ParseError) {
	return parseChain(input, names, c)
}

func parseChain(input string, names []string, specs []Spec) (Output, *ParseError) {
	sLen := len(specs)
	if sLen == 0 {
		return Output{
			Value:     []interface{}{},
			Consumed:  "",
			Remaining: input,
		}, nil
	}

	headOutput, headErr := specs[0].Parse(input, names)
	outputValue := []interface{}{headOutput.Value}
	if headErr != nil {
		headOutput.Value = outputValue
		return headOutput, headErr
	}

	tailOutput, tailErr := parseChain(headOutput.Remaining, names, specs[1:])
	outputValue = append(outputValue, tailOutput.Value.([]interface{})...)

	if tailErr != nil {
		tailErr.Offset += uint(len(headOutput.Consumed))
	}

	tailOutput.Value = outputValue
	tailOutput.Consumed = headOutput.Consumed + tailOutput.Consumed
	return tailOutput, tailErr
}

type Many struct {
	Spec  Spec
	Min   *uint
	Max   *uint
	Delim string
}

var _ Parser = Many{}

func (m Many) ToSpec() Spec {
	return Spec{
		Many: &m,
	}
}

func (m Many) ExpectedPrefix() string {
	switch {
	case m.Min != nil && m.Max != nil:
		return fmt.Sprintf("between %d and %d", *m.Min, *m.Max)
	case m.Min != nil:
		return fmt.Sprintf("%d or more", *m.Min)
	case m.Max != nil:
		return fmt.Sprintf("up to %d", *m.Max)
	default:
		return "any number of"
	}
}

func (m Many) Expected(names []string) []string {
	expected := []string{}
	prefix := m.ExpectedPrefix()
	for _, e := range m.Spec.Expected(names) {
		expected = append(expected, fmt.Sprintf("%s %s", prefix, e))
	}
	return expected
}

func (m Many) Parse(input string, names []string) (Output, *ParseError) {
	parsed := []interface{}{}
	if m.Max != nil && (*m.Max == 0 || m.Min != nil && *m.Min > *m.Max) {
		return Output{
			Value:     parsed,
			Remaining: input,
		}, nil
	}

	first := true
	offset := 0
	delim := Chain{
		Opt(Space{}.ToSpec()).ToSpec(),
		Token(m.Delim).ToSpec(),
		Opt(Space{}.ToSpec()).ToSpec(),
	}

	for {
		innerOffset := offset

		if !first {
			delimOutput, delimErr := delim.Parse(input[offset:], names)
			if delimErr != nil {
				break
			}
			innerOffset += len(delimOutput.Consumed)
		}
		first = false

		specOutput, specErr := m.Spec.Parse(input[innerOffset:], names)
		if specErr != nil {
			break
		}
		parsed = append(parsed, specOutput.Value)
		offset = innerOffset + len(specOutput.Consumed)

		if m.Max != nil && uint(len(parsed)) == *m.Max {
			break
		}
	}

	if m.Min != nil && uint(len(parsed)) < *m.Min {
		return Output{}, &ParseError{
			Message: fmt.Sprintf(
				"expected at least %d items but could only parse %d",
				*m.Min,
				len(parsed),
			),
			Offset: uint(offset),
		}
	}

	return Output{
		Value:     parsed,
		Consumed:  input[:offset],
		Remaining: input[offset:],
	}, nil
}

type Opt Spec

var _ Parser = Opt(Token("blah").ToSpec())

func (o Opt) ToSpec() Spec {
	return Spec{
		Opt: &o,
	}
}

func (o Opt) Expected(names []string) []string {
	expected := []string{}
	for _, e := range Spec(o).Expected(names) {
		expected = append(expected, fmt.Sprintf("optional %s", e))
	}
	return expected
}

func (o Opt) Parse(input string, names []string) (Output, *ParseError) {
	output, err := Spec(o).Parse(input, names)
	if err != nil {
		return Output{
			Value:     nil,
			Remaining: input,
		}, nil
	}
	return output, err
}

type Doc struct {
	Name string
	Desc string
	Spec Spec
}

var _ Parser = Doc{}

func (d Doc) ToSpec() Spec {
	return Spec{
		Doc: &d,
	}
}

func (d Doc) Expected(names []string) []string {
	return d.Spec.Expected(names)
}

func (d Doc) Parse(input string, names []string) (Output, *ParseError) {
	return d.Spec.Parse(input, names)
}

type Player struct{}

var _ Parser = Player{}

func (p Player) Parser(names []string) Enum {
	return Enum{
		Values: names,
	}
}

func (p Player) ToSpec() Spec {
	return Spec{
		Player: &p,
	}
}

func (p Player) Expected(names []string) []string {
	return p.Parser(names).Expected(names)
}

func (p Player) Parse(input string, names []string) (Output, *ParseError) {
	output, err := p.Parser(names).Parse(input, names)
	if err != nil {
		return output, err
	}
	name := output.Value.(string)
	output.Value = 0
	for k, n := range names {
		if n == name {
			output.Value = k
		}
	}
	return output, err
}

type Space struct{}

var _ Parser = Space{}

func (s Space) ToSpec() Spec {
	return Spec{
		Space: &s,
	}
}

func (s Space) Expected(names []string) []string {
	return []string{"whitespace"}
}

var SpaceRegexp = regexp.MustCompile(`^(\s+)`)

func (s Space) Parse(input string, names []string) (Output, *ParseError) {
	match := SpaceRegexp.FindString(input)
	if match == "" {
		return Output{
				Value:     "",
				Consumed:  "",
				Remaining: input,
			}, &ParseError{
				Message:  "expected whitespace",
				Expected: s.Expected(names),
			}
	}
	return Output{
		Value:     match,
		Consumed:  match,
		Remaining: input[len(match):],
	}, nil
}

func AfterSpace(spec Spec) Chain {
	return Chain{Space{}.ToSpec(), spec}
}
