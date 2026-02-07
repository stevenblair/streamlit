# All Widgets Showcase

Comprehensive example demonstrating **all available widgets** in the Streamlit Go backend.

## Features

This example showcases:

### 📝 Text Elements
- `st.Write()`, `st.Writef()` - Display text and formatted content
- `st.Title()`, `st.Header()`, `st.Subheader()` - Headings at different levels
- `st.Markdown()` - Rich markdown formatting with lists, bold, italic, code
- `st.Code()` - Syntax-highlighted code blocks
- `st.Divider()` - Visual separators

### ✅ Status Messages
- `st.Success()` - Green success alerts
- `st.Info()` - Blue informational alerts
- `st.Warning()` - Yellow warning alerts
- `st.Error()` - Red error alerts

### 🎛️ Input Widgets
- `st.TextInput()` - Single-line text input
- `st.Slider()` - Numeric slider with min/max/default
- `st.Checkbox()` - Boolean checkbox widget
- `st.Selectbox()` - Dropdown selection menu
- `st.Button()` - Clickable action button

### 📐 Layout Components
- `st.Columns()` - Multi-column responsive layouts
- `st.Metric()` - KPI displays with delta indicators
- `st.Sidebar()` - Sidebar container for controls

### 📊 Data Visualization
- `st.LineChart()` - Interactive line charts

### 🎉 Animations
- `st.Balloons()` - Celebration animation
- `st.Snow()` - Snow fall animation

## Running

```bash
go build -o all-widgets.exe .
.\all-widgets.exe
```

Then open http://localhost:8501

## What You'll See

1. **Comprehensive widget catalog** - Every available widget demonstrated
2. **Interactive sidebar** - Settings and configuration panel
3. **Live state management** - Widget values persist across interactions
4. **Real-time charts** - Adjustable data visualization
5. **Metrics dashboard** - KPI displays in multi-column layout
6. **Code samples** - Contextual examples for each language selection

## Code Structure

The example is organized into sections:
- Text elements demonstration
- Status/alert messages
- Input widgets with state management
- Button actions and animations
- Layout demonstrations
- Metrics and dashboards
- Data visualization
- Sidebar controls

## Use Cases

Perfect reference for:
- Learning all available Streamlit Go widgets
- Building your own dashboard applications
- Understanding widget state management
- Creating multi-column responsive layouts
- Implementing interactive data visualizations
