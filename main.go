package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

// Manufacturer struct represents a car manufacturer
type Manufacturer struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	FoundingYear int    `json:"foundingYear"`
}

// Category struct represents a car category
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Specifications struct represents the specifications of a car model
type Specifications struct {
	Engine       string `json:"engine"`
	Horsepower   int    `json:"horsepower"`
	Transmission string `json:"transmission"`
	Drivetrain   string `json:"drivetrain"`
}

// CarModel struct represents a car model
type CarModel struct {
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	ManufacturerID int            `json:"manufacturerId"`
	CategoryID     int            `json:"categoryId"`
	Year           int            `json:"year"`
	Specifications Specifications `json:"specifications"`
	Image          string         `json:"image"`
}

func fetchData(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, &target)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Define the port for the API
	const PORT = "3000"
	baseURL := "http://localhost:" + PORT + "/api"

	// Channels to receive data from goroutines
	carModelsCh := make(chan []CarModel)
	manufacturersCh := make(chan []Manufacturer)
	categoriesCh := make(chan []Category)

	// Fetch car models concurrently
	go func() {
		var carModels []CarModel
		err := fetchData(baseURL+"/models", &carModels)
		if err != nil {
			log.Printf("Error fetching car models: %v", err)
			carModelsCh <- nil
			return
		}
		carModelsCh <- carModels
	}()

	// Fetch manufacturers concurrently
	go func() {
		var manufacturers []Manufacturer
		err := fetchData(baseURL+"/manufacturers", &manufacturers)
		if err != nil {
			log.Printf("Error fetching manufacturers: %v", err)
			manufacturersCh <- nil
			return
		}
		manufacturersCh <- manufacturers
	}()

	// Fetch categories concurrently
	go func() {
		var categories []Category
		err := fetchData(baseURL+"/categories", &categories)
		if err != nil {
			log.Printf("Error fetching categories: %v", err)
			categoriesCh <- nil
			return
		}
		categoriesCh <- categories
	}()

	// Receive data from channels
	carModels := <-carModelsCh
	manufacturers := <-manufacturersCh
	categories := <-categoriesCh

	// Check if any of the fetch operations failed
	if carModels == nil || manufacturers == nil || categories == nil {
		log.Panic("Failed to fetch necessary data")
	}

	// Define routes to serve HTML content
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Parse the HTML template file
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			http.Error(w, "Oops! Something went wrong while rendering the page. Please try again later.", http.StatusInternalServerError)
			return
		}

		// Map to store manufacturer name by ID
		manufacturerNames := make(map[int]string)
		for _, m := range manufacturers {
			manufacturerNames[m.ID] = m.Name
		}

		// Map to store manufacturer country
		manufacturerCountry := make(map[int]string)
		for _, n := range manufacturers {
			manufacturerCountry[n.ID] = n.Country
		}

		// map to store manufacturer founding year
		manufacturerFoundingYear := make(map[int]int)
		for _, o := range manufacturers {
			manufacturerFoundingYear[o.ID] = o.FoundingYear
		}

		// Map to store category names by ID
		categoryNames := make(map[int]string)
		for _, c := range categories {
			categoryNames[c.ID] = c.Name
		}

		// Execute the template and pass the fetched data
		err = tmpl.Execute(w, struct {
			CarModels                []CarModel
			Manufacturers            []Manufacturer
			Categories               []Category
			ManufacturerNames        map[int]string
			ManufacturerCountry      map[int]string
			ManufacturerFoundingYear map[int]int
			CategoryMap              map[int]string
		}{
			CarModels:                carModels,
			Manufacturers:            manufacturers,
			Categories:               categories,
			ManufacturerNames:        manufacturerNames,
			ManufacturerCountry:      manufacturerCountry,
			ManufacturerFoundingYear: manufacturerFoundingYear,
			CategoryMap:              categoryNames,
		})
		if err != nil {
			http.Error(w, "Oops! Something went wrong while rendering the page. Please try again later.", http.StatusInternalServerError)
			return
		}
	})

	// Serve images from the API endpoint
	http.HandleFunc("/api/images/", func(w http.ResponseWriter, r *http.Request) {
		// Extract the image name from the request URL
		imageName := r.URL.Path[len("/api/images/"):]
		// Make a request to fetch the image
		resp, err := http.Get(baseURL + "/images/" + imageName)
		if err != nil {
			http.Error(w, "Image not found", http.StatusNotFound)
			return
		}
		defer resp.Body.Close()

		// Copy the image data to the response writer
		if _, err := io.Copy(w, resp.Body); err != nil {
			http.Error(w, "Error serving image", http.StatusInternalServerError)
			return
		}
	})

	// serve static files from the static directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// start server
	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
