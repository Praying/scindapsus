package app

type App struct {
}
type Exchange interface {
}
type AppBuilder struct {
}

func (this *App) Builder() AppBuilder {
	return AppBuilder{}
}
