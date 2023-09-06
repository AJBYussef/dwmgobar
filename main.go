package main
import (
    "strconv"
    "fmt"
    "time"
    "os"
    "os/exec"
    "encoding/json"
    "net/http"
    b "github.com/distatus/battery" 
)


func main()  {
    

    
    
    // var status string
    var tstring string
    var wstring string
    var bstring string

    currentTime := make(chan string)
    weather := make(chan string)
    battery := make(chan string)

    go getDate(currentTime)
    go getWeather(weather)
    // go getBattery(battery)

    for {
	select {
	case msg1 := <- currentTime:
	    tstring = msg1
	case msg2 := <- weather:
	    wstring = msg2
	case msg3 := <- battery:
	    bstring = msg3
	}
	cmd := exec.Command("xsetroot", "-name", bstring + wstring + tstring )
	err := cmd.Run()

	if err != nil {
	    fmt.Println(err)
	    os.Exit(-1)
	}

    }
    

	
}

func getWeather(weather chan string) {
    var i interface{}
    var output string
    for {
	wapi, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=39.23&longitude=9.12&hourly=temperature_2m,relativehumidity_2m,precipitation,windspeed_10m&current_weather=true&timezone=Europe%2FBerlin")

	if err != nil {
	    weather <- "Loading weather"
	    time.Sleep(time.Second * 10)
	    continue	
	}
	err = json.NewDecoder(wapi.Body).Decode(&i)

	djson := i.(map[string]interface{})
	current_weather := djson["current_weather"].(map[string]interface{})

	temp := fmt.Sprintf("%v",current_weather["temperature"])
	icon, _ := strconv.Atoi(fmt.Sprintf("%v", current_weather["weathercode"]))
	windspeed := fmt.Sprintf("%v", current_weather["windspeed"])

	output = fmt.Sprintf("%s°C🍃%skm/h%s" , temp, windspeed, selIcon(wmo_table(icon)))
	weather <- output

	time.Sleep(time.Hour)

    }
}

func getBattery(battery chan string) { 
    var output string
    var icon string
    for {
    bat, err := b.Get(0)
    if err != nil {
	battery <- "Could not get battery info!"
	time.Sleep(time.Minute * 1)
	continue	
    }
    percent := bat.Current / bat.Full * 100

    if bat.State == b.Discharging {
	if percent <= 20 { icon = "🪫"}  else { icon = "🔋" }
    } else { icon = "🔌" }

    output = fmt.Sprintf("%s%0.f%%",icon, percent)
    battery <- output
    time.Sleep(time.Second * 3)
    }
}

func getDate(currentTime chan string) {
    for {
	currentTime <- time.Now().Format("15:04:05")
	time.Sleep(time.Second)
    }
}



func selIcon(desc string) string {
    // h := time.Now().Hour()
    switch desc{
    case "Mist":
	return "🌫️"
    case "Snow":
	return "🌨️"
    case "Thunderstorm":
	return "⛈️"
    case "Rain":
	return "🌧️"
    case "Drizzle":
	return "🌧️"
    case "Clouds":
	return "☁"
    case "Clear":
	h := time.Now().Hour()
	if h >= 20 || h <= 6 {
	   return   "🌕"
	} else {
	    return "☀"
	}
    default:
	return "m"
    }
}


func wmo_table(code int) string{
    switch code{
	case 0: return "Clear"
	case 1,2,3: return "Clouds"
	case 45,48: return "Mist"
	case 51,53,55, 56, 57: return "Drizzle"
	case 61,63,65: return "Rain"
	case 71,73,75: return "Snow"
	case 95: return "Thunderstorm"
    default: return fmt.Sprintf("%d", code)
    }
}

