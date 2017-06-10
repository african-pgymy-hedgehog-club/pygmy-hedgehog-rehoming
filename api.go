package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/sc7639/sendmail"
)

const servername = "176.58.104.35:25"

type EmailTemplate struct {
	Title    string
	FormData map[string]template.HTML
	Order    []string
}

type APIResponse struct {
	Data    map[string]string `json:"data"`
	Success bool              `json:"success"`
	Error   string            `json:"error"`
}

// Transform string to have uppercase at the start of each word
func ucwords(str string) string {
	strArr := strings.Split(str, " ")
	for i := 0; i < len(strArr); i++ {
		strPart := strArr[i]
		strStart := strPart[:1]
		strEnd := strPart[1:]

		strArr[i] = strings.ToUpper(strStart) + strEnd
	}

	return strings.Join(strArr, " ")
}

// Transform form data for email template
func transformFormData(formData map[string][]string) map[string]template.HTML {
	var fd = make(map[string]template.HTML) // Make map to hold foorm data
	for key, val := range formData {        // Transform key to uppercase the bginning of words and replace _ with space
		key = strings.Replace(key, "_", " ", -1)
		if len(key) <= 3 {
			key = strings.ToUpper(key)
		} else {
			key = ucwords(key)
		}

		val[0] = strings.Replace(val[0], "\n", "<br />", -1)

		html := template.HTML([]byte(val[0]))

		fd[key] = html
	}

	return fd
}

// Function to create the connection to email server and sent the email
func sendEmail(from, subject, body string) (bool, error) {
	// Send email
	host, _, _ := net.SplitHostPort(servername)

	// Create new SendMail instance
	mail := sendmail.New(servername, from, subject, body)
	_, err := mail.AddToAddress("Admin", "admin@pygmyhedgehogrehoming.co.uk")
	if err != nil {
		return false, err
	}

	// Authenticate with smtp server
	mail.Auth("admin@pygmyhedgehogrehoming.co.uk", "rehoming123", host)
	mail.AddHeader("Content-Type", "text/html")
	return mail.Send()
}

// Parse api template
func apiTemplate(tmpl string) (*template.Template, error) {
	tmpl = filepath.Join(templateFolder, "email.html")
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Handle adoption api route calls
func adoptionHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(2 << 10); err != nil {
		clientError(w, err)
		return
	}

	fd := transformFormData(r.Form)

	t, err := apiTemplate("email")
	if err != nil {
		clientError(w, err)
		return
	}

	// Format date from Y-m-d to d/m/Y
	dob, err := time.Parse("2006-01-02", string(fd["DOB"]))
	if err != nil {
		clientError(w, err)
		return
	}
	fd["DOB"] = template.HTML(dob.Format("02/01/2006"))

	var eb bytes.Buffer           // Store email body
	t.Execute(&eb, EmailTemplate{ // Execute template with email template data struct
		Title:    "Adoption Form",
		FormData: fd,
		Order: []string{
			"Name",
			"Address",
			"DOB",
			"Email",
			"Number",
			"Colour",
			"Hedgehog Name",
			"Hedgehog Gender",
			"Hedgehog Colour",
			"Currently Own A Hedgehog",
			"Experience",
			"Plan To Keep In",
			"Plan To Feed",
		},
	})

	from := string(fd["Name"]) + " <" + string(fd["Email"]) + ">"

	ok, err := sendEmail(from, "Adoption Form", eb.String())
	if err != nil {
		clientError(w, err)
		return
	} else if !ok {
		clientError(w, err)
		return
	}

	// Create json from api response
	var resp = APIResponse{
		Success: true,
	}
	js, err := json.Marshal(resp)
	if err != nil {
		clientError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Handle contact api route calls
func contactHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(2 << 10); err != nil {
		clientError(w, err)
		return
	}

	fd := transformFormData(r.Form)

	from := string(fd["Name"]) + " <" + string(fd["Email"]) + ">"

	ok, err := sendEmail(from, string(fd["Subject"]), string(fd["Message"]))
	if err != nil {
		clientError(w, err)
		return
	} else if !ok {
		clientError(w, err)
		return
	}

	// Create json from api response
	var resp = APIResponse{
		Success: true,
	}
	js, err := json.Marshal(resp)
	if err != nil {
		clientError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Handle home for hog api route calls
func homeForHogHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(2 << 10); err != nil {
		clientError(w, err)
		return
	}

	fd := transformFormData(r.Form)

	t, err := apiTemplate("email")
	if err != nil {
		clientError(w, err)
		return
	}

	var eb bytes.Buffer // Store email body
	t.Execute(&eb, EmailTemplate{
		Title:    "Find a Home For Your Hog",
		FormData: fd,
		Order: []string{
			"Name",
			"Address",
			"Email",
			"Number",
			"Hedgehog Name",
			"Hedgehog Age",
			"Hedgehog Colour",
			"Additional Info",
		},
	})

	from := string(fd["Name"]) + " <" + string(fd["Email"]) + ">"

	ok, err := sendEmail(from, "Find a Home For Your Hog", eb.String())
	if err != nil {
		clientError(w, err)
		return
	} else if !ok {
		clientError(w, err)
		return
	}

	// Create json form api respose
	var resp = APIResponse{
		Success: true,
	}
	js, err := json.Marshal(resp)
	if err != nil {
		clientError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Handle forster carer api route calls
func fosterCarerHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(2 << 10); err != nil {
		clientError(w, err)
	}

	fd := transformFormData(r.Form)

	t, err := apiTemplate("email")
	if err != nil {
		clientError(w, err)
	}

	var eb bytes.Buffer // Store email body
	t.Execute(&eb, EmailTemplate{
		Title:    "Foster Carer Application",
		FormData: fd,
		Order: []string{
			"Name",
			"Address",
			"Email",
			"Number",
			"Why Would You Like To Foster",
			"What Setup Do You Have",
			"Any Further Information",
			"Over 18",
		},
	})

	from := string(fd["Name"]) + " <" + string(fd["Email"]) + ">"

	ok, err := sendEmail(from, "Foster Carer Application", eb.String())
	if err != nil {
		clientError(w, err)
		return
	} else if !ok {
		clientError(w, err)
	}

	// Create json from api response
	var resp = APIResponse{
		Success: true,
	}
	js, err := json.Marshal(resp)
	if err != nil {
		clientError(w, err)
		return
	}

	w.Header().Set("Conte-Type", "application/json")
	w.Write(js)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.Replace(r.URL.Path, "/api", "", 1)

	switch p {
	case "/adoption":
		adoptionHandler(w, r)
	case "/home-for-hog":
		homeForHogHandler(w, r)
	case "/foster-carer":
		fosterCarerHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	default:
		clientError(w, errors.New("This api route doesn't exist"))
	}
}
