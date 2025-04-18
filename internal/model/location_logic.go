package model

import (
	"math"
	"strconv"
)

func GetHitCount(locations []Location) int {
	return len(locations)
}

func GetCommonAddress(locations []Location) string {
	if len(locations) == 0 {
		return ""
	}
	pref := locations[0].Prefecture
	city := locations[0].City
	townPrefix := locations[0].Town

	for _, loc := range locations[1:] {
		if loc.Prefecture != pref {
			pref = ""
		}
		if loc.City != city {
			city = ""
		}
		townPrefix = getCommonPrefix(townPrefix, loc.Town)
	}

	result := pref + city
	if townPrefix != "" {
		result += townPrefix
	}
	return result
}

func getCommonPrefix(a, b string) string {
	arune := []rune(a)
	brune := []rune(b)
	minLen := len(arune)
	if len(brune) < len(arune) {
		minLen = len(b)
	}
	i := 0
	for i < minLen && arune[i] == brune[i] {
		i++
	}
	return string(arune[:i])
}

func GetFromTokyoStation(locations []Location) float64 {
	var maxDistance float64 = 0
	for _, loc := range locations {
		x, _ := strconv.ParseFloat(loc.X, 64)
		y, _ := strconv.ParseFloat(loc.Y, 64)
		dis := calculateDistanceFromTokyoStation(x, y)
		if dis > maxDistance {
			maxDistance = dis
		}
	}
	return roundToDecimal(maxDistance)
}

func calculateDistanceFromTokyoStation(x, y float64) float64 {
	const (
		R  = 6371.0
		xt = 139.7673068
		yt = 35.6809591
	)
	distance := (math.Pi / 180) * R * math.Sqrt(
		math.Pow((x-xt)*math.Cos(math.Pi*(y+yt)/360), 2)+math.Pow(y-yt, 2),
	)
	return distance
}

func roundToDecimal(f float64) float64 {
	return math.Round(f*10) / 10
}
