package gx

import "go.uber.org/fx"

func Named(constructor interface{}, name string) interface{} {
	return Annotate(constructor, ResultTags(`name:"`+name+`"`))
}

func ProvideAnnotated(constructor interface{}, anns ...Annotation) fx.Option {
	return Provide(
		Annotate(constructor, anns...),
	)
}

func ProvideNamed(constructor interface{}, name string) fx.Option {
	return Provide(
		Named(constructor, name),
	)
}
