package config

func (c Config) RequiredMetadata() []string {
	return c.Metadata.Required
}

func (c Config) GetWordCountRange() (int, int) {
	return c.Word.Min, c.Word.Max
}
