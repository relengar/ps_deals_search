package serviceadapter

type ServiceAdapterConfig struct{}

type ServiceAdapter struct {
}

func (s ServiceAdapter) CallEncoder() [][]float64 {
	return [][]float64{}
}
