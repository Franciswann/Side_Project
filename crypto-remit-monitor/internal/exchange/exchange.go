package exchange

type Exchange interface {
	GetName() string
	GetTicker(pair string) (float64, error)
}
