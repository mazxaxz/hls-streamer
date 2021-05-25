package hls

type Service struct {
	hlsDir string
}

func NewService(hlsDirectory string) *Service {
	c := Service{
		hlsDir: hlsDirectory,
	}
	return &c
}
