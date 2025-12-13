package stats

type Config struct {
	Infinite bool
	Basic    struct {
		Enabled     bool
		Interval    int
		TotalFrames int
	}
	Ecosystem struct {
		Enabled  bool
		Interval int
	}
}
