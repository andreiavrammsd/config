package interpolator

import (
	"strings"
	"unicode"
)

type value struct {
	content string
	i       int
}

func (ip *value) atEnd() bool {
	return ip.peek(1) == 0
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

type Interpolator struct {
	// All variables to be interpolated.
	vars map[string]string

	// Value before interpolation of current analyzed variable: `ABC $VAR TEXT`.
	rawValue value

	// Variable that is currently being interpolated
	interpolatedVar variable
}

// Interpolate converts each variable usage ($VAR, ${VAR}) inside a variable (value of map)
// to the actual value of the variable.
// It considers scenarios such as escaped variables and literal dollar signs.
//
// From:
//
//	A=1
//	B=text $A
//	C=\$B
//
// To:
//
//	A=1
//	B=text 1
//	C=$B
func (ip *Interpolator) Interpolate(vars map[string]string) {
	ip.vars = vars

	for key, value := range ip.vars {
		ip.rawValue.content = value

		if ip.rawValueContainsVars() {
			ip.parseVars()
			vars[key] = string(ip.interpolatedVar.value)
		}
	}
}

func (ip *Interpolator) rawValueContainsVars() bool {
	return strings.IndexByte(ip.rawValue.content, '$') != -1
}

func (ip *Interpolator) parseVars() {
	atVar := false // Notifies we're in the context of a variable: `text $IT_IS_HERE more text`.
	ip.interpolatedVar.name = nil
	ip.interpolatedVar.value = nil

	for ip.rawValue.i = 0; ip.rawValue.i < len(ip.rawValue.content); ip.rawValue.i++ {
		// Variable starts now. Continue to get its name and value.
		if ip.varStarts() {
			atVar = true
			continue
		}

		// Variable is between braces: ${VARIABLE}. Continue ignoring braces.
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
}

func (ip *Interpolator) varStarts() bool {
	// Litteral dollar sign at the end: `text $`.
	if isDolar(ip.rawValue.current()) && ip.rawValue.atEnd() {
		return false
	}

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

func (ip *Interpolator) appendToName() {
	ip.interpolatedVar.name = append(ip.interpolatedVar.name, ip.rawValue.current())
}

func (ip *Interpolator) appendCurrentCharacterToNewValue() {
	ip.interpolatedVar.value = append(ip.interpolatedVar.value, ip.rawValue.current())
}

func (ip *Interpolator) appendAllToNewValue() {
	ip.interpolatedVar.value = append(ip.interpolatedVar.value, []rune(ip.vars[string(ip.interpolatedVar.name)])...)
}

func isEscape(r rune) bool {
	return r == '\\'
}

func isDolar(r rune) bool {
	return r == '$'
}

func New() *Interpolator {
	return &Interpolator{}
}
