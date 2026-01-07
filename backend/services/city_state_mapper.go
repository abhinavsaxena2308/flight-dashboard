package services

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
)

// CityStateMapper handles city to state mapping
type CityStateMapper struct {
	cityToStateMap map[string]string
}

// Global instance of the city state mapper
var cityStateMapper *CityStateMapper
var mapperOnce sync.Once

// GetCityStateMapper returns a singleton instance of CityStateMapper
func GetCityStateMapper() *CityStateMapper {
	mapperOnce.Do(func() {
		cityStateMapper = &CityStateMapper{}
		// Initialize the mapping
		cityStateMapper.loadCityStateMap()
	})
	return cityStateMapper
}

// loadCityStateMap loads the city to state mapping from JSON file or initializes it
func (csm *CityStateMapper) loadCityStateMap() {
	// Try to load from JSON file first
	data, err := os.ReadFile("data/city_state_map.json")
	if err != nil {
		log.Println("Could not load city-state map from JSON file, using default mapping:", err)
		// Use default mapping if file doesn't exist
		csm.cityToStateMap = createDefaultCityStateMap()
	} else {
		// Parse the JSON file
		var rawMap map[string][]string
		if err := json.Unmarshal(data, &rawMap); err != nil {
			log.Printf("Error parsing city-state map JSON: %v, using default mapping", err)
			csm.cityToStateMap = createDefaultCityStateMap()
		} else {
			// Convert the raw map to a city-to-state mapping
			csm.cityToStateMap = make(map[string]string)
			for state, cities := range rawMap {
				for _, city := range cities {
					normalizedCity := strings.ToLower(strings.TrimSpace(city))
					csm.cityToStateMap[normalizedCity] = strings.ToLower(state)
				}
			}
		}
	}

	log.Printf("Loaded city-to-state mapping for %d cities", len(csm.cityToStateMap))
}

// GetStateForCity returns the state for a given city
func (csm *CityStateMapper) GetStateForCity(city string) (string, bool) {
	if city == "" {
		return "", false
	}

	// Normalize the city name for lookup
	normalizedCity := normalizeCityName(city)

	// Direct lookup
	if state, exists := csm.cityToStateMap[normalizedCity]; exists {
		return state, true
	}

	// Try common aliases
	if alias := getCityAlias(normalizedCity); alias != "" {
		if state, exists := csm.cityToStateMap[alias]; exists {
			return state, true
		}
	}

	// If not found, return empty string and false
	log.Printf("City not found in mapping: %s (normalized: %s)", city, normalizedCity)
	return "", false
}

// normalizeCityName normalizes city names for consistent lookup
func normalizeCityName(city string) string {
	normalized := strings.ToLower(strings.TrimSpace(city))

	// Handle common variations
	switch normalized {
	case "bombay":
		return "mumbai"
	case "new delhi":
		return "delhi"
	case "calcutta":
		return "kolkata"
	case "bangalore":
		return "bengaluru"
	}

	return normalized
}

// getCityAlias returns alternative names for cities that might be used
func getCityAlias(city string) string {
	// Common aliases for Indian cities
	aliases := map[string]string{
		"mumbai":      "mumbai",
		"bombay":      "mumbai",
		"delhi":       "delhi",
		"new delhi":   "delhi",
		"kolkata":     "kolkata",
		"calcutta":    "kolkata",
		"bengaluru":   "bengaluru",
		"bangalore":   "bengaluru",
		"madras":      "chennai",
		"chennai":     "chennai",
		"hyderabad":   "hyderabad",
		"pondy":       "puducherry",
		"ponducherry": "puducherry",
		"puducherry":  "puducherry",
	}

	if alias, exists := aliases[city]; exists {
		return alias
	}
	return ""
}

// createDefaultCityStateMap creates a default mapping of major Indian cities to states
func createDefaultCityStateMap() map[string]string {
	cityToState := make(map[string]string)

	// Andhra Pradesh
	andhraCities := []string{"amaravati", "visakhapatnam", "vijayawada", "guntur", "nellore", "kurnool", "rajahmundry", "tirupati", "kakinada", "kadapa", "anantapur", "eluru", "ongole", "kadiri", "hindupur", "proddatur", "bhimavaram", "gudivada", "rajampet", "tadepalligudem", "srikakulam", "anakapalle", "nandyal", "suriapet", "adoni", "chittoor", "machilipatnam", "bapatla", "nagari", "narsapur", "tanuku", "yemmiganur", "sullurpeta", "palacole", "parvathipuram", "ramachandrapuram", "samalkot", "sattenapalle", "tadpatri", "tiruvuru", "venkatagiri"}
	for _, city := range andhraCities {
		cityToState[strings.ToLower(city)] = "andhra pradesh"
	}

	// Arunachal Pradesh
	arunachalCities := []string{"itanagar", "naharlagun", "pasighat", "tawang", "bomdila", "tezu", "khonsa", "anini", "dambuk", "miao", "roing", "silapathar", "sagalee", "parang", "seppa", "bhalukpong", "changlang", "hawai", "jairampur", "koloriang", "lathao", "mohendraganj", "namsai", "pangin", "phassang", "ramsoh", "saiha", "sakoli", "sakrabaari", "tikabali", "zangla", "zirang", "ziro"}
	for _, city := range arunachalCities {
		cityToState[strings.ToLower(city)] = "arunachal pradesh"
	}

	// Assam
	assamCities := []string{"dispur", "guwahati", "dibrugarh", "silchar", "tezpur", "jorhat", "nagaon", "tinsukia", "dhubri", "diphu", "north lakhimpur", "barpeta", "lakhimpur", "nalgonda", "sibsagar", "goalpara", "hailakandi", "dhemaji", "teok", "lumding", "mangaldoi", "marigaon", "narkatiaganj", "sadiya", "udalguri", "badarpur", "bilasipara", "kharupatia", "lanka", "morigaon", "razampur", "sorbhog", "tangla"}
	for _, city := range assamCities {
		cityToState[strings.ToLower(city)] = "assam"
	}

	// Bihar
	biharCities := []string{"patna", "gaya", "bhagalpur", "muzaffarpur", "darbhanga", "begusarai", "chapra", "katihar", "munger", "purnia", "saharsa", "hajipur", "sasaram", "dehri", "aurangabad", "nawada", "jamalpur", "sitamarhi", "danapur", "madhubani", "siwan", "chhapra", "araria", "kishanganj", "madhepura", "arrah", "mokama", "sultanganj"}
	for _, city := range biharCities {
		cityToState[strings.ToLower(city)] = "bihar"
	}

	// Chhattisgarh
	chhattisgarhCities := []string{"raipur", "bhilai", "korba", "bilaspur", "raigarh", "jagdalpur", "rajnandgaon", "ambikapur", "dhamtari", "chirmiri", "bhatapara", "sakti", "jashpur", "mahasamund", "dantewada", "bijapur", "narayanpur", "kanker", "kondagaon", "sukma", "balod", "baloda bazar", "bemetara", "gariaband", "gaurela pendra marwahi", "kabirdham"}
	for _, city := range chhattisgarhCities {
		cityToState[strings.ToLower(city)] = "chhattisgarh"
	}

	// Goa
	goaCities := []string{"panaji", "margao", "mapusa", "mormugao", "bicholim", "ponda", "sanguem", "canacona", "pale", "quepem", "salcette", "cortalim", "cuncolim", "cunha", "goa velha", "jua", "kharebudr", "majorda", "mardol", "maria"}
	for _, city := range goaCities {
		cityToState[strings.ToLower(city)] = "goa"
	}

	// Gujarat
	gujaratCities := []string{"gandhinagar", "ahmedabad", "surat", "vadodara", "rajkot", "bhavnagar", "jamnagar", "nadiad", "veraval", "gandhidham", "bharuch", "junagadh", "bhuj", "navsari", "botad", "dahod", "devbhoomi dwarka", "gir somnath", "kheda", "mehsana", "morbi", "narmada", "panchmahal", "patan", "surendranagar", "tapi", "vadodara", "valsad"}
	for _, city := range gujaratCities {
		cityToState[strings.ToLower(city)] = "gujarat"
	}

	// Haryana
	haryanaCities := []string{"chandigarh", "faridabad", "gurgaon", "hisar", "rohtak", "panipat", "karnal", "sonipat", "yamunanagar", "bhiwani", "sirsa", "bahadurgarh", "jind", "thanesar", "kaithal", "palwal", "bawal", "charkhi dadri", "fatehabad", "gohana", "jagadhri", "kalka", "meham", "mewat", "narwana", "narnaul", "narnaund", "panchkula", "pundri", "radaur", "rajgarh", "safidon", "shahbad"}
	for _, city := range haryanaCities {
		cityToState[strings.ToLower(city)] = "haryana"
	}

	// Himachal Pradesh
	himachalCities := []string{"shimla", "mandi", "solan", "nahan", "bilaspur", "kullu", "dharamshala", "palampur", "baddi", "nagrota", "una", "chamba", "pangi", "lahaul", "spiti", "kangra", "kinnaur", "hamirpur"}
	for _, city := range himachalCities {
		cityToState[strings.ToLower(city)] = "himachal pradesh"
	}

	// Jharkhand
	jharkhandCities := []string{"ranchi", "jamshedpur", "dhanbad", "bokaro", "hazaribagh", "giridih", "deoghar", "chaibasa", "chatra", "dumka", "gumla", "pakur", "sahebganj", "simdega", "palamu", "latehar", "khunti", "littipara", "madhupur", "mihijam"}
	for _, city := range jharkhandCities {
		cityToState[strings.ToLower(city)] = "jharkhand"
	}

	// Karnataka
	karnatakaCities := []string{"bengaluru", "mysore", "mangalore", "hubli", "davanagere", "belgaum", "gulbarga", "bellary", "bijapur", "shimoga", "tumkur", "mandya", "gadag", "raichur", "hassan", "dharmavaram", "chitradurga", "kolar", "udupi", "hospet", "bhatkal", "gokak", "madikeri", "ranibennur", "shahabad", "tarikere"}
	for _, city := range karnatakaCities {
		cityToState[strings.ToLower(city)] = "karnataka"
	}

	// Kerala
	keralaCities := []string{"thiruvananthapuram", "kochi", "kollam", "kottayam", "palakkad", "alappuzha", "thrissur", "kannur", "kozhikode", "malappuram", "wayanad", "kasaragod", "pathanamthitta", "idukki"}
	for _, city := range keralaCities {
		cityToState[strings.ToLower(city)] = "kerala"
	}

	// Madhya Pradesh
	madhyaCities := []string{"bhopal", "indore", "jabalpur", "gwalior", "ujjain", "sagar", "dewas", "satna", "rewa", "morena", "hoshangabad", "bhind", "damoh", "khargone", "mandsaur", "neemuch", "shahdol", "chhindwara", "guna", "tikamgarh", "sehore", "vijaypur", "ashoknagar", "shajapur", "seoni"}
	for _, city := range madhyaCities {
		cityToState[strings.ToLower(city)] = "madhya pradesh"
	}

	// Maharashtra
	maharashtraCities := []string{"mumbai", "pune", "nagpur", "nashik", "aurangabad", "solapur", "thane", "jalgaon", "kolhapur", "amravati", "latur", "sangli", "nanded", "satara", "akola", "parbhani", "malegaon", "osmanabad", "nandurbar", "ahmednagar", "chandrapur", "dhule", "gondia", "hinganghat", "jalna", "khamgaon", "khopoli"}
	for _, city := range maharashtraCities {
		cityToState[strings.ToLower(city)] = "maharashtra"
	}

	// Manipur
	manipurCities := []string{"imphal", "thoubal", "bishnupur", "churachandpur", "senapati", "tamenglong", "ukhrul", "kakching", "kangpokpi", "noney", "phungyar", "tengnoupal"}
	for _, city := range manipurCities {
		cityToState[strings.ToLower(city)] = "manipur"
	}

	// Meghalaya
	meghalayaCities := []string{"shillong", "tura", "jowai", "nongstoin", "baghmara", "resubelpara", "williamnagar", "cherrapunji", "mairang", "mawkyrwat", "sohra", "nongpoh"}
	for _, city := range meghalayaCities {
		cityToState[strings.ToLower(city)] = "meghalaya"
	}

	// Mizoram
	mizoramCities := []string{"aizawl", "lunglei", "champhai", "kolasib", "serchhip", "mamit", "saiha", "dampa", "hachhek", "tawi", "thenzawl"}
	for _, city := range mizoramCities {
		cityToState[strings.ToLower(city)] = "mizoram"
	}

	// Nagaland
	nagalandCities := []string{"kohima", "dimapur", "mokokchung", "tuensang", "wokha", "zunheboto", "mon", "phek", "kiphire", "longleng"}
	for _, city := range nagalandCities {
		cityToState[strings.ToLower(city)] = "nagaland"
	}

	// Odisha
	odishaCities := []string{"bhubaneswar", "cuttack", "rourkela", "sambalpur", "berhampur", "puri", "balasore", "bhadrak", "baripada", "kendrapara", "anugul", "bargarh", "baleshwar", "balangir", "boudh", "bhawanipatna", "bijapur", "bolangir", "dhenkanal", "jagatsinghpur", "jajpur", "jharsuguda", "kalahandi", "kapoorthala", "kendujhar", "koraput", "malkangiri", "mayurbhanj", "nabarangpur", "nayagarh", "nuapada", "phulbani", "rayagada", "subarnapur", "sundargarh", "tumudibandha"}
	for _, city := range odishaCities {
		cityToState[strings.ToLower(city)] = "odisha"
	}

	// Punjab
	punjabCities := []string{"chandigarh", "ludhiana", "amritsar", "jalandhar", "patiala", "bathinda", "hoshiarpur", "moga", "mohali", "firozpur", "malerkotla", "mandi", "gobindgarh", "khanna", "fatehgarh sahib", "sangrur", "sunam", "dhuri", "zira", "fazilka", "kharar", "rajpura", "sirhind", "barnala", "jagraon", "kotkapura", "muktsar", "phagwara", "rampura", "tarn", "tarsikka", "dhariwal", "fatehgarh churian", "gurdaspur", "kapurthala", "rupnagar", "sas nagar", "sri muktsar sahib"}
	for _, city := range punjabCities {
		cityToState[strings.ToLower(city)] = "punjab"
	}

	// Rajasthan
	rajasthanCities := []string{"jaipur", "jodhpur", "kota", "bikaner", "ajmer", "bhilwara", "alwar", "sikar", "sawai madhopur", "pali", "ganganagar", "bharatpur", "barmer", "tonk", "chittorgarh", "dungarpur", "sri ganganagar", "banswara", "dhaulpur", "dholpur", "karauli", "pratapgarh", "rajsamand", "udaipur", "hanumangarh", "jaisalmer", "jalore", "jhalawar", "jhunjhunu", "nagaur", "pratapgarh", "sirohi", "todalgarh"}
	for _, city := range rajasthanCities {
		cityToState[strings.ToLower(city)] = "rajasthan"
	}

	// Sikkim
	sikkimCities := []string{"gangtok", "namchi", "gyalshing", "mangan", "soreng", "rajgung", "rhenock"}
	for _, city := range sikkimCities {
		cityToState[strings.ToLower(city)] = "sikkim"
	}

	// Tamil Nadu
	tamilnaduCities := []string{"chennai", "coimbatore", "madurai", "tiruchirappalli", "salem", "tirunelveli", "tiruppur", "vellore", "thoothukudi", "erode", "tiruvannamalai", "pollachi", "rajapalayam", "ramanathapuram", "kanchipuram", "nagercoil", "dindigul", "karur", "nagapattinam", "kovilpatti", "karaikudi", "vaniyambadi", "sivakasi", "tiruchengode", "tirupattur", "ranipet", "tindivanam", "udumalaipettai", "virudhachalam", "virudhunagar"}
	for _, city := range tamilnaduCities {
		cityToState[strings.ToLower(city)] = "tamil nadu"
	}

	// Telangana
	telanganaCities := []string{"hyderabad", "warangal", "nizamabad", "karimnagar", "ramagundam", "khammam", "mahbubnagar", "nalgonda", "suryapet", "miryalaguda", "siddipet", "adilabad", "sangareddy", "sircilla", "peddapalli", "bodhan", "mancherial", "kamareddy", "nirmal", "kotagiri"}
	for _, city := range telanganaCities {
		cityToState[strings.ToLower(city)] = "telangana"
	}

	// Tripura
	tripuraCities := []string{"agartala", "udaipur", "dharmanagar", "pratapgarh", "kailasahar", "belonia", "ampinagar", "khowai", "phuldungri", "jagtial"}
	for _, city := range tripuraCities {
		cityToState[strings.ToLower(city)] = "tripura"
	}

	// Uttar Pradesh
	upCities := []string{"lucknow", "kanpur", "agra", "varanasi", "meerut", "allahabad", "gorakhpur", "noida", "ghaziabad", "bareilly", "aligarh", "saharanpur", "mathura", "firozabad", "muzaffarnagar", "moradabad", "gwalior", "bhopal", "indore", "jabalpur", "raipur", "bilaspur", "durg", "raigarh", "jagdalpur", "rajnandgaon", "ambikapur", "dhamtari", "chirmiri", "bhatapara"}
	for _, city := range upCities {
		cityToState[strings.ToLower(city)] = "uttar pradesh"
	}

	// Uttarakhand
	uttarakhandCities := []string{"dehradun", "haridwar", "rishikesh", "udupi", "haldwani", "kathgodam", "kashipur", "rudrapur", "khatima", "sitarganj", "jaspur", "pauri", "chakrata", "chamoli", "dakpathar", "devprayag", "dhandhera", "dharasu", "dhumak", "dwarahat", "gairsain", "gangaikhera", "gangotri", "gauchar", "gaurikund", "guptkashi", "hardwar", "harsil", "jawalmukhi", "jhandi", "joshimath", "kalsi", "kanalich", "kanda", "karnaprayag", "kashipur", "khirshn", "kotdwar", "laksar", "lalkuan", "lansdowne", "lohardaga", "manali", "mandi", "manglaur", "mukteshwar", "nagla", "nainital", "nandaprayag", "narendranagar", "paddal", "padampur", "pauri", "phata", "pilkha", "pithoragarh", "pratapnagar", "prayag", "purola", "raipur", "raithal", "ranikhet", "roorkee", "rudraprayag", "uttarkashi", "vanspurnagar", "vijaypur", "vikasnagar", "virbhadra", "yamkeshwar", "yamunotri"}
	for _, city := range uttarakhandCities {
		cityToState[strings.ToLower(city)] = "uttarakhand"
	}

	// West Bengal
	wbCities := []string{"kolkata", "siliguri", "durgapur", "asansol", "malda", "raiganj", "kharagpur", "jalpaiguri", "cooch behar", "bankura", "darjeeling", "krishnanagar", "berhampore", "bally", "budge budge", "dhulian", "dankuni", "haldia", "kulti", "kamarhati", "medinipur", "nabadwip", "purulia", "shantipur", "suri", "tamluk", "alipurduar"}
	for _, city := range wbCities {
		cityToState[strings.ToLower(city)] = "west bengal"
	}

	// Delhi
	delhiCities := []string{"new delhi", "delhi", "north delhi", "south delhi", "east delhi", "west delhi", "central delhi", "north west delhi", "south west delhi", "north east delhi", "shahdara", "new delhi", "palam", "rohini", "pitampura", "karol bagh", "connaught place", "defence colony", "greater kailash", "hauz khas", "karkardooma", "lajpat nagar", "mayur vihar", "narela", "noida", "pandav nagar", "paschim vihar", "rajouri garden", "saket", "vasant kunj", "vishwas nagar", "yamuna vihar"}
	for _, city := range delhiCities {
		cityToState[strings.ToLower(city)] = "delhi"
	}

	// Puducherry
	puducherryCities := []string{"puducherry", "karaikal", "mahe", "yanaon", "pondicherry", "oussudu", "kannigapuram", "thattanchavady", "mannadipet", "mudaliarpet"}
	for _, city := range puducherryCities {
		cityToState[strings.ToLower(city)] = "puducherry"
	}

	// Andaman and Nicobar Islands
	andamanCities := []string{"port blair", "car nicobar", "coco island", "havelock island", "interview island", "jerry point", "katchal", "landfall island", "little andaman", "mahiya", "nancowry", "north and middle andaman", "north and south tillanchong", "north sentinel island", "phoenix bay", "ross island", "saddle peak", "south andaman", "swaraj dweep", "tillanchong", "viper island"}
	for _, city := range andamanCities {
		cityToState[strings.ToLower(city)] = "andaman and nicobar islands"
	}

	// Dadra and Nagar Haveli and Daman and Diu
	dnhDiuCities := []string{"daman", "dadar", "nagar haveli", "silvassa", "dnh", "dnh and dd", "dnh dd", "dnh diu"}
	for _, city := range dnhDiuCities {
		cityToState[strings.ToLower(city)] = "dadra and nagar haveli and daman and diu"
	}

	// Lakshadweep
	lakshadweepCities := []string{"kavaratti", "agatti", "andrott", "bitra", "chethlath", "kadmath", "kalpeni", "kiltan", "minicoy", "muhassar", "pandarani", "thinnakara"}
	for _, city := range lakshadweepCities {
		cityToState[strings.ToLower(city)] = "lakshadweep"
	}

	// Ladakh
	ladakhCities := []string{"leh", "kargil", "drass", "padum", "zanskar", "nubra", "pangong", "tso", "tso kar", "tso moriri", "changthang", "nyoma", "urud"}
	for _, city := range ladakhCities {
		cityToState[strings.ToLower(city)] = "ladakh"
	}

	return cityToState
}
