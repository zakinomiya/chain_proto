package gateway

type Gateway struct {
}

func (g *Gateway) ServiceName() string {
	return "Gateway"
}

func (g *Gateway) Start() error {
	return nil
}

func (g *Gateway) Stop() {

}
