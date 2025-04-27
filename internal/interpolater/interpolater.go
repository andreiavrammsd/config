package interpolater

import (
	"strings"
	"unicode"
)

type value struct {
	content string
	i       int
}

// endsWithDollarSign tests for: `NAME=text + $`.
func (ip *value) endsWithDollarSign() bool {
	return ip.atEnd() && !isEscape(ip.peek(-1))
}

func (ip *value) atEnd() bool {
	return ip.i == len(ip.content)-1
}

// nextVarIsDoubleEscaped detects: `\\$“.
func (ip *value) nextVarIsDoubleEscaped() bool {
	return isEscape(ip.current()) && isEscape(ip.peek(1)) && isDolar(ip.peek(2))
}

// nextVarIsEscaped detects: `\$“.
func (ip *value) nextVarIsEscaped() bool {
	return isEscape(ip.current()) && isDolar(ip.peek(1))
}

func (ip *value) current() rune {
	return rune(ip.content[ip.i])
}

func (ip *value) peek(steps int) rune {
	if ip.i+steps < 0 || ip.i+steps >= len(ip.content) {
		return 0
	}

	return rune(ip.content[ip.i+steps])
}

func (ip *value) isOpenBrace() bool {
	return ip.current() == '{'
}

func (ip *value) isCloseBrace() bool {
	return ip.current() == '}'
}

func (ip *value) isAtSpace() bool {
	return unicode.IsSpace(ip.current())
}

type variable struct {
	name  []rune
	value []rune
}

type Interpolater struct {
	// All variables to be interpolated.
	vars map[string]string

	// Value before interpolation of current analyzed variable: `ABC $VAR TEXT`.
	rawValue value

	// Variable that is currently being interpolated
	interpolatedVar variable
}

func (ip *Interpolater) Interpolate(vars map[string]string) {
	ip.vars = vars

	for key, value := range ip.vars {
		ip.rawValue.content = value

		if !ip.containsVar() {
			continue
		}

		atVar := false // Notifies we're in the context of a variable: `text $IT_IS_HERE text`.
		ip.interpolatedVar.name = nil
		ip.interpolatedVar.value = nil

		for ip.rawValue.i = 0; ip.rawValue.i < len(ip.rawValue.content); ip.rawValue.i++ {
			if ip.varStarts() {
				atVar = true

				// If value ends in $, literally append $. Do not interpret as variable.
				if ip.rawValue.endsWithDollarSign() {
					ip.appendCurrentCharacterToNewValue()
				}

				continue
			}

			// ${VARIABLE}
			if atVar && (ip.rawValue.isOpenBrace() || ip.rawValue.isCloseBrace()) {
				continue
			}

			if !atVar {
				if ip.rawValue.nextVarIsDoubleEscaped() {
					ip.appendCurrentCharacterToNewValue()

					continue
				}

				if ip.rawValue.nextVarIsEscaped() {
					continue
				}

				// Append literal.
				ip.appendCurrentCharacterToNewValue()

				continue
			}

			// Variable ends when a space is found. Append literal.
			if ip.rawValue.isAtSpace() {
				ip.appendAllToNewValue()
				ip.appendCurrentCharacterToNewValue()
				ip.interpolatedVar.name = nil
				atVar = false

				continue
			}

			ip.appendToName()
		}

		if atVar {
			ip.appendAllToNewValue()
		}

		vars[key] = string(ip.interpolatedVar.value)
	}
}

func (ip *Interpolater) containsVar() bool {
	return strings.IndexByte(ip.rawValue.content, '$') != -1
}

func (ip *Interpolater) varStarts() bool {
	// Normal variable: `text $VAR text`. Will use its value.
	if !isDolar(ip.rawValue.current()) {
		return false
	}

	// Variable is double escaped: `text \\$VAR text`. The escape character is actually escaped.
	// Will use an escape character and variable's value.
	if isEscape(ip.rawValue.peek(-2)) && isEscape(ip.rawValue.peek(-1)) {
		return true
	}

	// Variable is escaped: `text \$VAR text`. Actually not a variable. Will use it literally.
	if isEscape(ip.rawValue.peek(-1)) {
		return false
	}

	next := ip.rawValue.peek(1)
	if unicode.IsSpace(next) || next == '"' || next == '\'' {
		return false
	}

	return true
}

func (ip *Interpolater) appendToName() {
	ip.interpolatedVar.name = append(ip.interpolatedVar.name, ip.rawValue.current())
}

func (ip *Interpolater) appendCurrentCharacterToNewValue() {
	ip.interpolatedVar.value = append(ip.interpolatedVar.value, ip.rawValue.current())
}

func (ip *Interpolater) appendAllToNewValue() {
	ip.interpolatedVar.value = append(ip.interpolatedVar.value, []rune(ip.vars[string(ip.interpolatedVar.name)])...)
}

func isEscape(r rune) bool {
	return r == '\\'
}

func isDolar(r rune) bool {
	return r == '$'
}

func New() *Interpolater {
	return &Interpolater{}
}
