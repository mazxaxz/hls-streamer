package hls

type Options struct {
	VideoBitrate string
	VideoMaxRate string
	VideoBufSize string

	AudioBitrate      string
	AudioSamplingRate int
}
