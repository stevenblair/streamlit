# Waveform Generator Example

This example demonstrates how to perform complex processing in Go code and display the results in the Streamlit UI. It generates mathematical waveforms based on user input parameters and performs real-time calculations for visualization and analysis.

## Features

- **Sidebar Controls**: All input parameters conveniently located in the sidebar
- **Real-time Go Processing**: All waveform calculations are performed in Go code
- **Interactive Parameters**: Control frequency, amplitude, phase, and waveform type
- **Statistical Analysis**: Automatically calculates min, max, mean, energy, peak-to-peak, and RMS values
- **ASCII Visualization**: Generates an ASCII art plot of the waveform
- **Raw Data Display**: Optionally view the calculated waveform points

## What This Example Demonstrates

1. **Sidebar Usage**: Shows how to use `st.Sidebar()` to organize input controls
2. **Complex Calculations in Go**: The waveform generation functions (`generateWaveform`, `calculateStats`, `calculateEnergy`, `calculateRMS`) show how to perform mathematical operations in Go
3. **Data Flow**: User adjusts parameters → Go recalculates waveform → UI updates with new values
4. **Multiple Widget Types**: Sliders, selectbox, checkbox, columns, and metrics all working together
5. **Real-time Updates**: Changes to any parameter immediately trigger recalculation

## Waveform Types

- **Sine**: Classic sinusoidal wave
- **Square**: Digital square wave
- **Triangle**: Linear triangular wave
- **Sawtooth**: Asymmetric sawtooth wave

## Running the Example

```bash
# Build all examples
./scripts/build.sh  # or build.bat on Windows

# Run the waveform example
./bin/waveform  # or bin\waveform.exe on Windows
```

Then open your browser to http://localhost:8501

## Code Structure

- `generateWaveform()`: Core calculation function that generates waveform points based on parameters
- `calculateStats()`: Computes min, max, and mean values
- `calculateEnergy()`: Calculates signal energy
- `calculateRMS()`: Computes root mean square value
- `generateASCIIPlot()`: Creates a visual representation using ASCII characters
- `app()`: Main Streamlit app function that creates the UI and orchestrates processing

## Key Concepts

This example showcases the power of doing computation in Go:

1. **Performance**: Go's performance makes it suitable for real-time calculations
2. **Type Safety**: Go's type system ensures correctness of mathematical operations
3. **Concurrency Ready**: The code structure supports future concurrent processing
4. **Native Data Types**: Direct use of Go slices and structs for efficient data handling

## Extending This Example

You could extend this example to:

- Add more waveform types (cosine, noise, custom functions)
- Implement Fourier transforms for frequency analysis
- Add waveform mixing/modulation capabilities
- Save/load waveform configurations
- Export waveform data to files
- Add real-time audio synthesis
