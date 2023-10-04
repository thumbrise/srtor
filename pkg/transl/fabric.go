package transl

import "srtor/pkg/util"

const EnvTranslateDebug = "TRANSLATE_DEBUG"

type Translator interface {
	Translate(source string, sourceLang string, targetLang string) (string, error)
}

func NewEnvBasedTranslator() Translator {
	if util.EnvGetBool(EnvTranslateDebug) {
		return DevTranslator{}
	}

	return GoogleTranslator{}
}
