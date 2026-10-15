package geo

import (
	"fmt"
	"math"
	"strings"
)

const earthKm = 6371.0

// HaversineKm menghitung jarak (km) dari koordinat manual. Tanpa API/GPS.
func HaversineKm(lat1, lng1, lat2, lng2 float64) float64 {
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*math.Sin(dLng/2)*math.Sin(dLng/2)
	return 2 * earthKm * math.Asin(math.Sqrt(a))
}

// FormatID memformat km gaya Indonesia: "0,8 km".
func FormatID(km float64) string {
	s := fmt.Sprintf("%.1f km", km)
	return strings.Replace(s, ".", ",", 1)
}
