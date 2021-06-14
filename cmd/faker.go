package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/tidwall/gjson"
)

var (
	ErrInvalidBlankJSON  = errors.New("glazer: faker JSON cannot be empty")
	ErrInvalidJSONFormat = errors.New("glazer: faker data is not a validation JSON format")
	ErrInvalidObjectJSON = errors.New("glazer: faker data is not a validation object JSON")
)

func generateFakerData(raw json.RawMessage) (faker map[string]interface{}, err error) {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 {
		return nil, ErrInvalidBlankJSON
	}

	if !gjson.ValidBytes(raw) {
		return nil, ErrInvalidJSONFormat
	}

	parsed := gjson.ParseBytes(raw)
	if !parsed.IsObject() {
		return nil, ErrInvalidObjectJSON
	}

	return generateFakerObject(parsed), nil
}

func generateFakerObject(o gjson.Result) (faker map[string]interface{}) {
	if !o.IsObject() {
		return nil
	}

	objs := o.Map()
	faker = make(map[string]interface{}, len(objs))
	for k, v := range objs {
		faker[k] = parseResult(v)
	}

	return faker
}

func generateFakerArray(a gjson.Result) (faker []interface{}) {
	if !a.IsArray() {
		return nil
	}

	arrs := a.Array()
	if len(arrs) == 0 {
		return make([]interface{}, 0)
	}

	faker = make([]interface{}, 0, len(arrs))
	for _, r := range arrs {
		faker = append(faker, parseResult(r))
	}
	return faker
}

func parseResult(r gjson.Result) interface{} {
	var ret interface{}
	switch r.Type {
	case gjson.String:
		ret = doFaker([]byte(r.Str))
	case gjson.JSON:
		if r.IsArray() {
			ret = generateFakerArray(r)
		}
		if r.IsObject() {
			ret = generateFakerObject(r)
		}
	default:
		ret = r.Value()
	}
	return ret
}

func doFaker(name []byte) string {
	if fn, isGen := fakerGenerators[string(name)]; isGen {
		return fn()
	}
	return string(name)
}

const (
	pwLower     = true
	pwUpper     = true
	pwNumeric   = true
	pwSpecial   = false
	pwSpace     = false
	pwNum       = 10
	sentenceLen = 10
	logLevel    = "info"
)

var fakerGenerators = map[string]func() string{
	// Person
	"@fake:Name":           gofakeit.Name,
	"@fake:NamePrefix":     gofakeit.NamePrefix,
	"@fake:NameSuffix":     gofakeit.NameSuffix,
	"@fake:FirstName":      gofakeit.FirstName,
	"@fake:LastName":       gofakeit.LastName,
	"@fake:Gender":         gofakeit.Gender,
	"@fake:SSN":            gofakeit.SSN,
	"@fake:Email":          gofakeit.Email,
	"@fake:Phone":          gofakeit.Phone,
	"@fake:PhoneFormatted": gofakeit.PhoneFormatted,

	// Auth
	"@fake:Username": gofakeit.Username,
	"@fake:Password": func() string { return gofakeit.Password(pwLower, pwUpper, pwNumeric, pwSpecial, pwSpace, pwNum) },

	// Address
	"@fake:City":         gofakeit.City,
	"@fake:Country":      gofakeit.Country,
	"@fake:CountryAbr":   gofakeit.CountryAbr,
	"@fake:State":        gofakeit.State,
	"@fake:StateAbr":     gofakeit.StateAbr,
	"@fake:Street":       gofakeit.Street,
	"@fake:StreetName":   gofakeit.StreetName,
	"@fake:StreetNumber": gofakeit.StreetNumber,
	"@fake:StreetPrefix": gofakeit.StreetPrefix,
	"@fake:StreetSuffix": gofakeit.StreetSuffix,
	"@fake:Zip":          gofakeit.Zip,
	"@fake:Latitude":     func() string { return fmt.Sprintf("%.2f", gofakeit.Latitude()) },
	"@fake:Longitude":    func() string { return fmt.Sprintf("%.2f", gofakeit.Longitude()) },
	// LatitudeInRange(min, max float64) (float64, error)  // TODO
	// LongitudeInRange(min, max float64) (float64, error) // TODO

	// Game
	"@fake:Gamertag": gofakeit.Gamertag,

	// Beer
	"@fake:BeerAlcohol": gofakeit.BeerAlcohol,
	"@fake:BeerBlg":     gofakeit.BeerBlg,
	"@fake:BeerHop":     gofakeit.BeerHop,
	"@fake:BeerIbu":     gofakeit.BeerIbu,
	"@fake:BeerMalt":    gofakeit.BeerMalt,
	"@fake:BeerName":    gofakeit.BeerName,
	"@fake:BeerStyle":   gofakeit.BeerStyle,
	"@fake:BeerYeast":   gofakeit.BeerYeast,

	// Car
	"@fake:CarMaker":            gofakeit.CarMaker,
	"@fake:CarModel":            gofakeit.CarModel,
	"@fake:CarType":             gofakeit.CarType,
	"@fake:CarFuelType":         gofakeit.CarFuelType,
	"@fake:CarTransmissionType": gofakeit.CarTransmissionType,

	// Words
	"@fake:Noun":               gofakeit.Noun,
	"@fake:Verb":               gofakeit.Verb,
	"@fake:Adverb":             gofakeit.Adverb,
	"@fake:Preposition":        gofakeit.Preposition,
	"@fake:Adjective":          gofakeit.Adjective,
	"@fake:Word":               gofakeit.Word,
	"@fake:Sentence":           func() string { return gofakeit.Sentence(sentenceLen) },
	"@fake:LoremIpsumWord":     gofakeit.LoremIpsumWord,
	"@fake:LoremIpsumSentence": func() string { return gofakeit.LoremIpsumSentence(sentenceLen) },
	"@fake:Question":           gofakeit.Question,
	"@fake:Quote":              gofakeit.Quote,
	"@fake:Phrase":             gofakeit.Phrase,
	// "@fake:Paragraph": func() string { return gofakeit.Paragraph() }, // TODO
	// "@fake:LoremIpsumParagraph": func() string { return	gofakeit.LoremIpsumParagraph() }, // TODO

	// Foods
	"@fake:Fruit":     gofakeit.Fruit,
	"@fake:Vegetable": gofakeit.Vegetable,
	"@fake:Breakfast": gofakeit.Breakfast,
	"@fake:Lunch":     gofakeit.Lunch,
	"@fake:Dinner":    gofakeit.Dinner,
	"@fake:Snack":     gofakeit.Snack,
	"@fake:Dessert":   gofakeit.Dessert,

	// Misc
	"@fake:UUID":      gofakeit.UUID,
	"@fake:FlipACoin": gofakeit.FlipACoin,

	// Colors
	"@fake:Color":     gofakeit.Color,
	"@fake:HexColor":  gofakeit.HexColor,
	"@fake:RGBColor":  func() string { return rgbString() },
	"@fake:SafeColor": gofakeit.SafeColor,

	// Internel
	"@fake:URL":                  gofakeit.URL,
	"@fake:DomainName":           gofakeit.DomainName,
	"@fake:DomainSuffix":         gofakeit.DomainSuffix,
	"@fake:IPv4Address":          gofakeit.IPv4Address,
	"@fake:IPv6Address":          gofakeit.IPv6Address,
	"@fake:MacAddress":           gofakeit.MacAddress,
	"@fake:HTTPStatusCode":       func() string { return strconv.Itoa(gofakeit.HTTPStatusCode()) },
	"@fake:HTTPStatusCodeSimple": func() string { return strconv.Itoa(gofakeit.HTTPStatusCodeSimple()) },
	"@fake:LogLevel":             func() string { return gofakeit.LogLevel(logLevel) },
	"@fake:HTTPMethod":           gofakeit.HTTPMethod,
	"@fake:UserAgent":            gofakeit.UserAgent,
	"@fake:ChromeUserAgent":      gofakeit.ChromeUserAgent,
	"@fake:FirefoxUserAgent":     gofakeit.FirefoxUserAgent,
	"@fake:OperaUserAgent":       gofakeit.OperaUserAgent,
	"@fake:SafariUserAgent":      gofakeit.SafariUserAgent,

	// Date/Time
	"@fake:Date":           func() string { return gofakeit.Date().String() },
	"@fake:NanoSecond":     func() string { return strconv.Itoa(gofakeit.NanoSecond()) },
	"@fake:Second":         func() string { return strconv.Itoa(gofakeit.Second()) },
	"@fake:Minute":         func() string { return strconv.Itoa(gofakeit.Minute()) },
	"@fake:Hour":           func() string { return strconv.Itoa(gofakeit.Hour()) },
	"@fake:Month":          func() string { return strconv.Itoa(gofakeit.Month()) },
	"@fake:MonthString":    gofakeit.MonthString,
	"@fake:Day":            func() string { return strconv.Itoa(gofakeit.Day()) },
	"@fake:WeekDay":        gofakeit.WeekDay,
	"@fake:Year":           func() string { return strconv.Itoa(gofakeit.Year()) },
	"@fake:TimeZone":       gofakeit.TimeZone,
	"@fake:TimeZoneAbv":    gofakeit.TimeZoneAbv,
	"@fake:TimeZoneFull":   gofakeit.TimeZoneFull,
	"@fake:TimeZoneOffset": func() string { return fmt.Sprintf("%.3f", gofakeit.TimeZoneOffset()) },
	"@fake:TimeZoneRegion": gofakeit.TimeZoneRegion,
	// "@fake:DateRange(start,DateRange end time.Time) time.Time // TODO
}

func rgbString() string {
	rgb := gofakeit.RGBColor()
	var bs bytes.Buffer
	for _, b := range rgb {
		bs.WriteString(strconv.Itoa(b))
		bs.WriteRune(',')
	}
	return bs.String()
}
