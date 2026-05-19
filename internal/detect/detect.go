package detect

import "unicode"

type script int

const (
	scriptUnknown  script = iota
	scriptLatin
	scriptCyrillic
	scriptJapanese // Hiragana or Katakana — unambiguously Japanese
	scriptCJK      // Han characters — Chinese (or Japanese kanji without kana)
	scriptHangul
	scriptArabic
	scriptHebrew
)

// langScript maps ISO 639-1 codes to their primary identifying script.
var langScript = map[string]script{
	"ru": scriptCyrillic, "uk": scriptCyrillic, "bg": scriptCyrillic,
	"sr": scriptCyrillic, "mk": scriptCyrillic, "be": scriptCyrillic,
	"ja": scriptJapanese,
	"zh": scriptCJK,
	"ko": scriptHangul,
	"ar": scriptArabic, "fa": scriptArabic,
	"he": scriptHebrew,
	"en": scriptLatin, "es": scriptLatin, "fr": scriptLatin, "de": scriptLatin,
	"pt": scriptLatin, "it": scriptLatin, "nl": scriptLatin, "pl": scriptLatin,
	"cs": scriptLatin, "sk": scriptLatin, "ro": scriptLatin, "hr": scriptLatin,
}

func scriptForLang(lang string) script {
	if s, ok := langScript[lang]; ok {
		return s
	}
	return scriptUnknown
}

// dominantScript returns the primary non-Latin script found in s,
// falling back to scriptLatin or scriptUnknown when no other script is present.
// Hiragana/Katakana takes priority over Han to distinguish Japanese from Chinese.
func dominantScript(s string) script {
	counts := make(map[script]int)
	for _, r := range s {
		switch {
		case unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r):
			counts[scriptJapanese]++
		case unicode.Is(unicode.Cyrillic, r):
			counts[scriptCyrillic]++
		case unicode.Is(unicode.Hangul, r):
			counts[scriptHangul]++
		case unicode.Is(unicode.Han, r):
			counts[scriptCJK]++
		case unicode.Is(unicode.Arabic, r):
			counts[scriptArabic]++
		case unicode.Is(unicode.Hebrew, r):
			counts[scriptHebrew]++
		case (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z'):
			counts[scriptLatin]++
		}
	}
	for _, sc := range []script{scriptJapanese, scriptCyrillic, scriptHangul, scriptCJK, scriptArabic, scriptHebrew} {
		if counts[sc] > 0 {
			return sc
		}
	}
	if counts[scriptLatin] > 0 {
		return scriptLatin
	}
	return scriptUnknown
}

// Direction returns (from, to) language codes based on the dominant script of text
// relative to the configured language pair (lang1, lang2).
//
// When the detected script matches exactly one configured language, that language
// is the source. For ambiguous pairs (e.g. both Latin), defaults to lang2 → lang1.
func Direction(text, lang1, lang2 string) (from, to string) {
	detected := dominantScript(text)
	s1 := scriptForLang(lang1)
	s2 := scriptForLang(lang2)

	if detected == s1 && detected != s2 {
		return lang1, lang2
	}
	if detected == s2 && detected != s1 {
		return lang2, lang1
	}
	// Ambiguous or unknown — treat input as lang2, translate to lang1
	return lang2, lang1
}
