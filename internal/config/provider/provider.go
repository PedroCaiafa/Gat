package provider

type Provider interface {
	Embed(text string) ([]float32, error)

	InferDependencies(code string) ([]string, error)
}
