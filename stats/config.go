package stats

type Config struct {
	Basic struct {
		Enabled     bool
		Interval    int
		TotalFrames int
	}
	Ecosystem struct {
		Enabled  bool
		Interval int
	}
}
