// Question: Develop a weather monitoring system where multiple display units (e.g., temperature display, humidity display) need to be updated whenever there is a change in weather data.
// Scenario: Create a WeatherStation class that notifies all registered observers whenever there is a change in weather data.


package main

import "fmt"

// Observer interface
type Observer interface{
	update(temperature float64, humidity float64)
	display()
}

// WeatherStation struct
type WeatherStation struct{
	observers []Observer
	temperature float64
	humidity float64
}

// RegisterObserver function
func (ws *WeatherStation) RegisterObserver(Observer Observer){
	ws.observers = append(ws.observers, Observer)
}

// RemoveObserver function
func (ws *WeatherStation) RemoveObserver(observer Observer){
	for i, o := range ws.observers{
		if o == observer{
			ws.observers = append(ws.observers[:i], ws.observers[i+1:]...)
		}
	}
}

// NotifyObservers function
func (ws *WeatherStation) NotifyObservers(){
	for _, o := range ws.observers{
		o.update(ws.temperature, ws.humidity)
	}
}

// SetMeasurements function
func (ws *WeatherStation) SetMeasurements(temperature float64, humidity float64){
	ws.temperature = temperature
	ws.humidity = humidity
	ws.NotifyObservers()
}

// TemperatureDisplay struct
type TemperatureDisplay struct{
	Temperature float64
	Humidity float64
}

// Update function
func (td *TemperatureDisplay) update(temperature float64, humidity float64){
	td.Temperature = temperature
	td.Humidity = humidity
}

// Display function
func (td *TemperatureDisplay) display(){
	fmt.Printf("Temperature : %f\n", td.Temperature)
}

// HumidityDisplay struct
type HumidityDisplay struct{
	Temperature float64
	Humidity float64
}

// Update function
func (hd *HumidityDisplay) update(temperature float64, humidity float64){
	hd.Temperature = temperature
	hd.Humidity = humidity
}

// Display function
func (hd *HumidityDisplay) display(){
	fmt.Printf("Humidity : %f\n", hd.Humidity)
}

func main(){
	weatherStation := &WeatherStation{}
	temperatureDisplay	:= &TemperatureDisplay{}
	humidityDisplay := &HumidityDisplay{}

	weatherStation.RegisterObserver(temperatureDisplay)
	weatherStation.RegisterObserver(humidityDisplay)

	weatherStation.SetMeasurements(30.0, 65.3)

	temperatureDisplay.display()
	humidityDisplay.display()

}
