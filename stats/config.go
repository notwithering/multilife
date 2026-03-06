package stats

type Config struct {
	Enabled  bool
	Interval int
	Infinite bool

	Basic struct {
		Enabled     bool
		TotalFrames int
		FPS         int
	}
	Timings struct {
		Enabled  bool
		Interval int
	}
	Ecosystem struct {
		Enabled  bool
		Interval int
	}
}
