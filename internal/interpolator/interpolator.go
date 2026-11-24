package interpolator

import (
	"strings"
	"unicode"
	"unicode/utf8"
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
	return isEscape(ip.current()) && isEscape(ip.peek(1)) && isDolar(ip.peek(2)) //nolint: mnd
}

// nextVarIsEscaped detects: `\$“.
func (ip *value) nextVarIsEscaped() bool {
	return isEscape(ip.current()) && isDolar(ip.peek(1))
}

func (ip *value) current() rune {
	current, _ := utf8.DecodeRuneInString(ip.content[ip.i:])
	return current
}

func (ip *value) peek(steps int) rune {
	position := ip.i

	if steps < 0 {
		for range -steps {
			if position == 0 {
				return 0
			}

			position--
			for position > 0 && !utf8.RuneStart(ip.content[position]) {
				position--
			}
		}
	} else {
		for range steps {
			_, size := utf8.DecodeRuneInString(ip.content[position:])
			position += size
		}
	}

	if position >= len(ip.content) {
		return 0
	}

	next, _ := utf8.DecodeRuneInString(ip.content[position:])

	return next
}

func (ip *value) advance() {
	_, size := utf8.DecodeRuneInString(ip.content[ip.i:])
	ip.i += size
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

	for key, value := range vars {
		ip.rawValue.content = value

		if !ip.rawValueContainsVars() {
			continue
		}

		ip.parseVars()
		vars[key] = string(ip.interpolatedVar.value)
	}
}

func (ip *Interpolator) rawValueContainsVars() bool {
	return strings.IndexByte(ip.rawValue.content, '$') != -1
}

func (ip *Interpolator) parseVars() {
	ip.interpolatedVar.value = nil
	ip.parseValue(&ip.rawValue, 0)
}

func (ip *Interpolator) parseValue(rawValue *value, depth int) {
	atVar := false // Notifies we're in the context of a variable: `text $IT_IS_HERE more text`.

	var name []rune

	for rawValue.i = 0; rawValue.i < len(rawValue.content); rawValue.advance() {
		// Variable starts now. Continue to get its name and value.
		if ip.varStarts(rawValue) {
			atVar = true
			continue
		}

		// Variable is between braces: ${VARIABLE}. Continue ignoring braces.
		if atVar && (rawValue.isOpenBrace() || rawValue.isCloseBrace()) {
			continue
		}

		if !atVar {
			if rawValue.nextVarIsDoubleEscaped() {
				ip.appendCurrentCharacterToNewValue(rawValue.current())
				continue
			}

			if rawValue.nextVarIsEscaped() {
				continue
			}

			// Append literal.
			ip.appendCurrentCharacterToNewValue(rawValue.current())

			continue
		}

		// Variable ends when a space is found. Append literal.
		if rawValue.isAtSpace() {
			ip.appendVariable(name, depth)
			ip.appendCurrentCharacterToNewValue(rawValue.current())

			name = nil
			atVar = false

			continue
		}

		name = append(name, rawValue.current())
	}

	if atVar {
		ip.appendVariable(name, depth)
	}
}

func (ip *Interpolator) varStarts(rawValue *value) bool {
	// Litteral dollar sign at the end: `text $`.
	if isDolar(rawValue.current()) && rawValue.atEnd() {
		return false
	}

	// Normal variable: `text $VAR text`. Will use its value.
	if !isDolar(rawValue.current()) {
		return false
	}

	// Variable is double escaped: `text \\$VAR text`. The escape character is actually escaped.
	// Will use an escape character and variable's value.
	if isEscape(rawValue.peek(-2)) && isEscape(rawValue.peek(-1)) {
		return true
	}

	// Variable is escaped: `text \$VAR text`. Actually not a variable. Will use it literally.
	if isEscape(rawValue.peek(-1)) {
		return false
	}

	next := rawValue.peek(1)
	if unicode.IsSpace(next) || next == '"' || next == '\'' {
		return false
	}

	return true
}

func (ip *Interpolator) appendVariable(name []rune, depth int) {
	raw := ip.vars[string(name)]
	if depth >= len(ip.vars) || !strings.Contains(raw, "$") {
		ip.interpolatedVar.value = append(ip.interpolatedVar.value, []rune(raw)...)
		return
	}

	nestedValue := value{content: raw}
	ip.parseValue(&nestedValue, depth+1)
}

func (ip *Interpolator) appendCurrentCharacterToNewValue(current rune) {
	ip.interpolatedVar.value = append(ip.interpolatedVar.value, current)
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
