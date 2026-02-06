package main

import (
	"fmt"
	"math"

	st "github.com/stevenblair/streamlit/go-backend"
)

func main() {
	st.Run(app)
}

func app() {
	st.Title("🌊 Waveform Generator")
	st.Write("Generate and visualize waveforms with real-time Go calculations")

	// Input parameters in sidebar
	sidebar := st.Sidebar()
	sidebar.Header("Parameters")

	frequency := sidebar.Slider("Frequency (Hz)", 1, 10, 2)
	amplitude := sidebar.Slider("Amplitude", 1, 10, 5)
	phase := sidebar.Slider("Phase (°)", 0, 360, 0)

	waveType := sidebar.Selectbox(
		"Waveform Type",
		[]string{"Sine", "Square", "Triangle", "Sawtooth"},
		0,
	)

	samples := sidebar.Slider("Number of Samples", 10, 5000, 1000)

	sidebar.Divider()
	sidebar.Subheader("About")
	sidebar.Write("This example demonstrates real-time processing in Go with UI updates.")

	st.Divider()

	// Generate waveform in Go
	st.Header("Generated Waveform")

	points := generateWaveform(waveType, frequency, amplitude, phase, int(samples))

	// Display statistics
	st.Subheader("Statistics")
	stats := calculateStats(points)

	metricCols := st.Columns(4)
	metricCols[0].Metric("Samples", fmt.Sprintf("%d", len(points)), "")
	metricCols[1].Metric("Min", fmt.Sprintf("%.2f", stats.min), "")
	metricCols[2].Metric("Max", fmt.Sprintf("%.2f", stats.max), "")
	metricCols[3].Metric("Mean", fmt.Sprintf("%.2f", stats.mean), "")

	st.Divider()

	// Display waveform data
	st.Subheader("Waveform Data")

	if st.Checkbox("Show raw data", false) {
		displayRawData(points)
	}

	// Visualization
	st.Subheader("Visualization")
	st.Write(fmt.Sprintf("📈 **%s Waveform** - %d samples | Frequency: %.0f Hz | Amplitude: %.0f",
		waveType, len(points), frequency, amplitude))

	// Display waveform chart
	st.LineChart(points)

	st.Divider()

	// Fourier analysis simulation
	st.Subheader("Analysis")

	energy := calculateEnergy(points)
	peakToPeak := stats.max - stats.min
	rms := calculateRMS(points)

	analysisCols := st.Columns(3)
	analysisCols[0].Metric("Energy", fmt.Sprintf("%.2f", energy), "")
	analysisCols[1].Metric("Peak-to-Peak", fmt.Sprintf("%.2f", peakToPeak), "")
	analysisCols[2].Metric("RMS", fmt.Sprintf("%.2f", rms), "")

	st.Info(fmt.Sprintf("💡 Generated %s wave at %.0f Hz with amplitude %.0f",
		waveType, frequency, amplitude))
}

// WaveformStats holds statistical information about a waveform
type WaveformStats struct {
	min  float64
	max  float64
	mean float64
}

// generateWaveform creates a waveform based on parameters
func generateWaveform(waveType string, freq, amp, phase float64, samples int) []float64 {
	points := make([]float64, samples)
	phaseRad := phase * math.Pi / 180.0

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(samples) * 2 * math.Pi * freq

		switch waveType {
		case "Sine":
			points[i] = amp * math.Sin(t+phaseRad)
		case "Square":
			if math.Sin(t+phaseRad) >= 0 {
				points[i] = amp
			} else {
				points[i] = -amp
			}
		case "Triangle":
			// Triangle wave using arcsin approximation
			val := (2.0 / math.Pi) * math.Asin(math.Sin(t+phaseRad))
			points[i] = amp * val
		case "Sawtooth":
			// Sawtooth wave
			val := 2.0 * (t/math.Pi - math.Floor(0.5+t/math.Pi))
			points[i] = amp * val
		default:
			points[i] = amp * math.Sin(t+phaseRad)
		}
	}

	return points
}

// calculateStats computes statistical measures
func calculateStats(points []float64) WaveformStats {
	if len(points) == 0 {
		return WaveformStats{}
	}

	stats := WaveformStats{
		min:  points[0],
		max:  points[0],
		mean: 0,
	}

	sum := 0.0
	for _, p := range points {
		if p < stats.min {
			stats.min = p
		}
		if p > stats.max {
			stats.max = p
		}
		sum += p
	}

	stats.mean = sum / float64(len(points))
	return stats
}

// calculateEnergy computes signal energy
func calculateEnergy(points []float64) float64 {
	energy := 0.0
	for _, p := range points {
		energy += p * p
	}
	return energy
}

// calculateRMS computes root mean square
func calculateRMS(points []float64) float64 {
	if len(points) == 0 {
		return 0
	}
	sumSquares := 0.0
	for _, p := range points {
		sumSquares += p * p
	}
	return math.Sqrt(sumSquares / float64(len(points)))
}

// displayRawData shows the waveform points
func displayRawData(points []float64) {
	st.Write(fmt.Sprintf("Showing %d data points:", len(points)))

	// Display in rows of 10
	text := ""
	for i, p := range points {
		text += fmt.Sprintf("%6.2f ", p)
		if (i+1)%10 == 0 {
			text += "\n"
		}
	}

	st.Code(text, "")
}
