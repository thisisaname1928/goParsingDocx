package docx

const (
	Italic    = 6
	Bold      = 2
	Marked    = 3
	Underline = 4
	ImgSource = 5
)

type Prop struct {
	Type  int
	Value string
}

type FluidProperty struct {
	Start    int
	End      int
	Property []Prop
}

type FluidString struct {
	Text       string
	Properties []FluidProperty
}

func calLen(inp string) int {
	return len([]rune(inp))
}

func isAnswerKey(ch rune, ch2 rune) bool {
	if ch == 'A' || ch == 'B' || ch == 'C' || ch == 'D' {
		if ch2 == '.' {
			return true
		}
	}

	if ch == 'a' || ch == 'b' || ch == 'c' || ch == 'd' {
		if ch2 == ')' {
			return true
		}
	}

	return false
}

func removeRune(s []rune, i int) []rune {
	return append(s[:i], s[i+1:]...)
}

func String2Fluid(content string) []FluidString {
	var res []FluidString
	aRune := []rune(content)

	for i := 0; i < len(aRune); {
		curRes := ""

		for aRune[i] != '\n' {
			curRes += string(aRune[i])
			i++

			if i >= len(aRune) {
				break
			}
		}

		aaRune := []rune(curRes)
		var prop []FluidProperty

		for k := 0; k < len(aaRune)-2; {
			if aaRune[k] == '*' && isAnswerKey(aaRune[k+1], aaRune[k+2]) {
				if k+3 >= len(aaRune) || aaRune[k+3] == ' ' || aaRune[k+3] == '\t' || aaRune[k+3] == ':' || aaRune[k+3] == ')' {
					aaRune = removeRune(aaRune, k)
					prop = append(prop, FluidProperty{Start: k, End: k + 1, Property: []Prop{{Marked, "yellow"}}})
					k += 2
					continue
				}
			}
			k++
		}

		curRes = string(aaRune)

		res = append(res, FluidString{Text: curRes, Properties: prop})

		i++
	}

	return res
}

func MakeFluidStringInstance(fluid FluidString) FluidString {
	var res FluidString
	res.Text = fluid.Text[:]

	res.Properties = make([]FluidProperty, len(fluid.Properties))
	for i := range res.Properties {
		res.Properties[i].Property = make([]Prop, len(fluid.Properties[i].Property))
		res.Properties[i].Start = fluid.Properties[i].Start
		res.Properties[i].End = fluid.Properties[i].End

		for j := range res.Properties[i].Property {
			res.Properties[i].Property[j].Type = fluid.Properties[i].Property[j].Type
			res.Properties[i].Property[j].Value = fluid.Properties[i].Property[j].Value
		}
	}

	return res
}

func Parse2Fluid(path string) ([]FluidString, error) {
	var str []FluidString

	doc, e := Parse(path)

	if e != nil {
		return nil, e
	}

	paras := doc.Body.Paragraphs
	rIDTab := GetRID(path)

	for _, v := range paras { // parse per line
		var currentLine FluidString
		for _, r := range v.Runs {
			if r.Drawing != nil { // parse drawing
				for _, drawing := range *r.Drawing {
					var currentProperty FluidProperty
					currentLine.Text += ""
					currentProperty.Start = calLen(currentLine.Text)
					currentProperty.End = currentProperty.Start
					var p Prop
					p.Type = ImgSource
					p.Value = rIDTab[ParseDrawing(&drawing)]

					currentProperty.Property = append(currentProperty.Property, p)
					currentLine.Properties = append(currentLine.Properties, currentProperty)
				}
			}
			if r.Text.Text != nil {
				if r.RunProperties.Properties != nil { // parse properties
					var currentProperty FluidProperty
					currentProperty.Start = calLen(currentLine.Text)
					currentProperty.End = calLen(*r.Text.Text) + currentProperty.Start

					for _, pr := range *r.RunProperties.Properties {
						var p Prop
						switch pr.XMLName.Local {
						case "b":
							p.Type = Bold
							if pr.Val != nil {
								p.Value = *pr.Val
							}
						case "i":
							p.Type = Italic
							if pr.Val != nil {
								p.Value = *pr.Val
							}
						case "u":
							p.Type = Underline
							if pr.Val != nil {
								p.Value = *pr.Val
							}
						case "shd":
							p.Type = Marked
							if pr.Fill != nil {
								p.Value = *pr.Fill
							}
						case "highlight":
							p.Type = Marked
							if pr.Fill != nil {
								p.Value = *pr.Fill
							}
						}
						currentProperty.Property = append(currentProperty.Property, p)

					}
					currentLine.Properties = append(currentLine.Properties, currentProperty)
				}
				currentLine.Text += *r.Text.Text

			}
		}
		currentLine.Text += " "
		str = append(str, currentLine)

	}

	return str, nil
}

func DelFirstCharacterStr(str *string) {
	r := []rune(*str)
	*str = string(r[1:])
}

func DelNCharacterStr(str *string, n int) {
	for i := 0; i < n; i++ {
		DelFirstCharacterStr(str)
	}
}

func DelFirstCharacterRune(str *[]rune) {
	tmp := *str
	*str = tmp[1:]
}

func DelNCharacterRune(str *[]rune, n int) {
	for i := 0; i < n; i++ {
		DelFirstCharacterRune(str)
	}
}

func DelFirstCharacter(str *FluidString) {
	tmpText := ""
	r := []rune(str.Text) // go use utf8 by default so this is the safe way to delete
	if len(r) > 0 {
		tmpText = string(r[1:])
	}

	for i := range str.Properties {
		if str.Properties[i].Start > 0 {
			str.Properties[i].Start--
		}

		if str.Properties[i].End > 0 {
			str.Properties[i].End--
		}
	}

	str.Text = tmpText
}

func DelNCharacter(str *FluidString, n int) {
	for i := 0; i < n; i++ {
		DelFirstCharacter(str)
	}
}

func parseFluid2HtmlInternal(str FluidString, allowMark bool) string {
	output := ""
	addLabel := false

	text := []rune(str.Text)
	for _, v := range text {
		if v != ' ' {
			addLabel = true
			break
		}
	}

	for i := 0; i <= len(text); i++ {
		for _, prop := range str.Properties {
			if prop.End == i && i != 0 {
				for _, pr := range prop.Property {
					if pr.Type == Bold && pr.Value != "false" {
						output += "</b>"
					}
					if pr.Type == Italic && pr.Value != "false" {
						output += "</i>"
					}
					if pr.Type == Underline && pr.Value != "false" {
						output += "</u>"
					}
					if allowMark && pr.Type == Marked && pr.Value != "auto" {
						output += "</mark>"
					}
				}
			}

			if prop.Start == i {
				for _, pr := range prop.Property {
					if pr.Type == ImgSource {
						output += "<img src=\"" + pr.Value + "\">"
					}
					if pr.Type == Bold && pr.Value != "false" {
						output += "<b>"
					}
					if pr.Type == Italic && pr.Value != "false" {
						output += "<i>"
					}
					if pr.Type == Underline && pr.Value != "false" {
						output += "<u>"
					}
					if allowMark && pr.Type == Marked && pr.Value != "auto" {
						output += "<mark>"
					}
				}
			}
		}

		if i < len(text) {
			output += string(text[i])
		}
	}

	if addLabel {
		output = "<label class=\"ques_content\">" + output + "</label>"
	}
	return output
}

func ParseFluid2Html(str FluidString) string {
	return parseFluid2HtmlInternal(str, true)
}

func ParseFluid2HtmlNonMark(str FluidString) string {
	return parseFluid2HtmlInternal(str, false)
}

// a1, a2 as source index, b1 b2 as chop index
func chopRange(a1 int, a2 int, b1 int, b2 int) (int, int) {
	c1, c2 := -1, -1

	if b1 >= a1 && b2 <= a2 {
		c1, c2 = b1, b2
	} else if b1 < a1 && b2 <= a2 && b2 >= a1 { // b2 in range [a1,a2] and b1 not in
		c1, c2 = a1, b2

	} else if b1 >= a1 && b1 <= a2 && b2 > a2 { // b1 in range [a1,a2] and b2 not in
		c1, c2 = b1, a2

	} else if b1 < a1 && b2 > a2 {
		c1, c2 = a1, a2

	}

	return c1, c2
}

func CopyFluid(fluid FluidString, beginIndex int, endIndex int) FluidString {
	var res FluidString
	aRune := []rune(fluid.Text)
	if beginIndex < 0 {
		beginIndex = 0
	}
	if endIndex >= len(aRune) {
		endIndex = len(aRune) - 1
	}
	if beginIndex <= endIndex && len(aRune) > 0 {
		res.Text = string(aRune[beginIndex : endIndex+1])
	}
	nonCopyLen := beginIndex

	// check 4 property
	for i := range fluid.Properties {
		b, e := chopRange(fluid.Properties[i].Start, fluid.Properties[i].End, beginIndex, endIndex)
		if min(b, e) != -1 {
			var prop FluidProperty
			prop.Property = make([]Prop, len(fluid.Properties[i].Property))
			copy(prop.Property, fluid.Properties[i].Property)
			prop.Start, prop.End = b-nonCopyLen, e-nonCopyLen
			res.Properties = append(res.Properties, prop)
		}
	}

	return res
}

func ConcatFluid(fluid1 FluidString, fluid2 FluidString) FluidString {
	concatIndex := calLen(fluid1.Text)
	var res FluidString
	res.Text = fluid1.Text[:] + fluid2.Text[:]

	// add fluid 1 properties
	for i := range fluid1.Properties {
		var prop = make([]Prop, len(fluid1.Properties[i].Property))

		copy(prop, fluid1.Properties[i].Property)

		res.Properties = append(res.Properties, FluidProperty{fluid1.Properties[i].Start, fluid1.Properties[i].End, prop})
	}

	// add fluid 2 properties
	for i := range fluid2.Properties {
		var prop = make([]Prop, len(fluid2.Properties[i].Property))

		Start := fluid2.Properties[i].Start + concatIndex
		End := fluid2.Properties[i].End + concatIndex

		copy(prop, fluid2.Properties[i].Property)

		res.Properties = append(res.Properties, FluidProperty{Start, End, prop})
	}

	return res
}
