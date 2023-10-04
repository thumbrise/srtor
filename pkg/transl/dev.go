package transl

type DevTranslator struct{}

func (t DevTranslator) Translate(source, sourceLang, targetLang string) (string, error) {
	return source, nil
}
