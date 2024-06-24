package controllers

import (
	"complaints/cmd/api/models"
	"complaints/cmd/api/utilities"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
}

var jwtKeyEncoded = os.Getenv("JWTKEY")

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}

	//hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//store in the database
	user.Password = string(hashedPassword)
	_, err = models.Register(user)
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, user, "user")
}

func Login(w http.ResponseWriter, r *http.Request) {

	var credentials models.LoginCredentials
	type jsonResp struct {
		OK      bool   `json:"ok"`
		Message string `json:"message"`
		UserID  string `json:"user_id"`
		Token   string `json:"token"`
		Role    string `json:"role"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bad := jsonResp{
		OK:      false,
		Message: "Invalid username or password",
		UserID:  "null",
	}

	//find user by ID
	user, err := models.GetUserByUserID(credentials.Username)
	if err != nil {
		utilities.WriteJSON(w, http.StatusUnauthorized, bad, "response")
		return
	}

	//compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	//Generste JWT token
	jwtKey, err := base64.URLEncoding.DecodeString(jwtKeyEncoded)
	if err != nil {
		http.Error(w, "Error decoding JWT key", http.StatusInternalServerError)
		fmt.Println("Error decoding JWT key:", err)
		return
	}
	type CustomClaims struct {
		jwt.StandardClaims
		Role string `json:"role,omitempty"`
	}
	expirationTime := time.Now().Add(4320 * time.Minute)
	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Subject:   user.ID.Hex(),
			Issuer:    credentials.Username,
		},
		Role: user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ok := jsonResp{
		OK:      true,
		Message: "Login successful",
		UserID:  credentials.Username,
		Token:   tokenString,
		Role:    user.Role,
	}

	//send token in response
	w.Header().Set("Authorization", tokenString)
	utilities.WriteJSON(w, http.StatusOK, ok, "response")
}

func NewComplaint(w http.ResponseWriter, r *http.Request) {
	// Limit the size of the incoming file
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB max
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	courseConcerned := r.FormValue("course_concerned")
	requestDetails := r.FormValue("request_details")
	testScore, err := strconv.Atoi(r.FormValue("test_score"))
	if err != nil {
		fmt.Println(err)
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}
	defer file.Close()

	//ensure upload directory exists
	uploadsDir := "uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err = os.Mkdir(uploadsDir, os.ModePerm)
		if err != nil {
			utilities.ErrorJSON(w, err)
			return
		}
	}

	//create new fle name and save the file
	dst, err := os.Create(filepath.Join(uploadsDir, handler.Filename))
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		utilities.ErrorJSON(w, err)
		return
	}

	studentId, ok := r.Context().Value("userID").(string)
	if !ok {
		utilities.ErrorJSON(w, errors.New("unable to get student ID from context"))
		return
	}

	complaint := models.Complaint{
		RequestingStudent: studentId,
		CourseConcerned:   courseConcerned,
		RequestDetails:    requestDetails,
		TestScore:         testScore,
		StudentProof:      fmt.Sprintf("/uploads/%s", handler.Filename),
		Status:            "Pending",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	course, err := models.GetCourseByCourseCode(string(complaint.CourseConcerned))
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}
	if len(course.Lecturers) == 0 {
		utilities.ErrorJSON(w, errors.New("no lecturers found for the course"))
		return
	}
	respondingLecturer := course.Lecturers[rand.Intn(len(course.Lecturers))]

	complaint.RespondingLecturer = respondingLecturer
	exists, err := models.ComplaintAlreadyExists(studentId, courseConcerned)
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}
	if !exists {
		_, err = models.CreateNewComplaint(complaint)
		if err != nil {
			utilities.ErrorJSON(w, err)
			return
		}
		body := fmt.Sprintf("You have a new revalidation request from %s concerning %s. \n Log in to your portal to see the full details", studentId, courseConcerned)
		lecturer, err := models.GetUserByUserID(respondingLecturer)
		if err != nil {
			utilities.ErrorJSON(w, err)
			return
		}
		email := lecturer.Email
		subject := fmt.Sprintf("New revalidation request for %s", courseConcerned)
		err = SendEmail(body, email, subject)
		if err != nil {
			utilities.ErrorJSON(w, err)
		}
	} else {
		utilities.ErrorJSON(w, fmt.Errorf("you already have an existing complaint for this course"))
	}
}

func ExtractEmailFromRequest(r *http.Request) string {
	var requestData struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		return ""
	}
	return requestData.Email
}

func GetComplaintByObjectID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}

	complaint, err := models.GetComplaintByObjectId(objID)
	if err != nil {
		fmt.Println("Unable to get complaint")
		utilities.ErrorJSON(w, err)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, complaint, "complaint")
}

func GetComplaintsByStaffID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	complaints, err := models.GetComplaintsByStaffId(id)
	if err != nil {
		fmt.Println("Unable to get complaints", err)
		utilities.ErrorJSON(w, err)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, complaints, "complaints")
}

func GetComplaintsByStudentID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	complaints, err := models.GetComplaintsByStudentId(id)
	if err != nil {
		fmt.Println("Unable to get complaints", err)
		utilities.ErrorJSON(w, err)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, complaints, "complaints")
}

func ChangeComplaintStatusByLecturer(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")
	var updatedComplaint models.Complaint

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB max
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	// Extract reason from form data
	reason := r.FormValue("reason")
	if reason == "" {
		http.Error(w, "Reason is required", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}
	defer file.Close()

	//ensure upload directory exists
	uploadsDir := "uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err = os.Mkdir(uploadsDir, os.ModePerm)
		if err != nil {
			utilities.ErrorJSON(w, err)
			return
		}
	}

	//create new file name and save the file
	dst, err := os.Create(filepath.Join(uploadsDir, handler.Filename))
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		utilities.ErrorJSON(w, err)
		return
	}

	updatedComplaint.Reason = reason

	lecturerProof := fmt.Sprintf("/uploads/%s", handler.Filename)
	updatedComplaint.LecturerProof = lecturerProof

	err = models.ChangeComplaintStatusLecturer(id, "Approved By Lecturer", updatedComplaint.Reason, updatedComplaint.LecturerProof)
	if err != nil {
		utilities.ErrorJSON(w, err)
	}

	body := fmt.Sprintf("Your revalidation request, concerning %s, which was assigned to %s has been approved by the lecturer. \n Log in to your portal to see more details", updatedComplaint.CourseConcerned, updatedComplaint.RespondingLecturer)
	student, err := models.GetUserByUserID(updatedComplaint.RequestingStudent)
	if err != nil {
		utilities.ErrorJSON(w, err)
	}
	email := student.Email
	subject := fmt.Sprintf("Update Concerning Revalidation Request for %s", updatedComplaint.CourseConcerned)
	err = SendEmail(body, email, subject)
	if err != nil {
		utilities.ErrorJSON(w, err)
	}

	body2 := fmt.Sprintf("You have a new revalidation request through %s concerning %s for %s. \n Log in to your portal to see the full details", updatedComplaint.RespondingLecturer, updatedComplaint.CourseConcerned, updatedComplaint.RequestingStudent)
	hod, err := models.GetUserByUserID("23001")
	if err != nil {
		utilities.ErrorJSON(w, err)
	}
	email2 := hod.Email
	subject2 := fmt.Sprintf("Revalidation Request from %s concerning %s", updatedComplaint.RespondingLecturer, updatedComplaint.CourseConcerned)
	err = SendEmail(body2, email2, subject2)
	if err != nil {
		utilities.ErrorJSON(w, err)
	}

	utilities.WriteJSON(w, http.StatusOK, "Status Updated Successfully", "Success")
}

func ChangeComplaintStatusByAdvisor(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	var updatedComplaint models.Complaint
	err := json.NewDecoder(r.Body).Decode(&updatedComplaint)
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}

	err = models.ChangeComplaintStatus(id, "Approved By Course Advisor")
	if err != nil {
		utilities.ErrorJSON(w, err)
	}

	utilities.WriteJSON(w, http.StatusOK, "Status Updated Successfully", "Success")
}

func ChangeComplaintStatusByHOD(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	var updatedComplaint models.Complaint
	err := json.NewDecoder(r.Body).Decode(&updatedComplaint)
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}

	err = models.ChangeComplaintStatus(id, "Approved By HOD")
	if err != nil {
		utilities.ErrorJSON(w, err)
	}
	body := fmt.Sprintf("Your revalidation request, concerning %s, which was assigned to %s has been approved by your HOD. \n Log in to your portal to see more details", updatedComplaint.CourseConcerned, updatedComplaint.RespondingLecturer)
	student, err := models.GetUserByUserID(updatedComplaint.RequestingStudent)
	if err != nil {
		utilities.ErrorJSON(w, err)
	}
	email := student.Email
	subject := fmt.Sprintf("Update Concerning Revalidation Request for %s", updatedComplaint.CourseConcerned)
	err = SendEmail(body, email, subject)
	if err != nil {
		utilities.ErrorJSON(w, err)
	}

	body2 := fmt.Sprintf("You have a new revalidation request through %s concerning %s for %s. \n Log in to your portal to see the full details", updatedComplaint.RespondingLecturer, updatedComplaint.CourseConcerned, updatedComplaint.RequestingStudent)
	senate, err := models.GetUserByUserID("19201")
	if err != nil {
		utilities.ErrorJSON(w, err)
	}
	email2 := senate.Email
	subject2 := fmt.Sprintf("Revalidation Request from %s concerning %s", updatedComplaint.RespondingLecturer, updatedComplaint.CourseConcerned)
	err = SendEmail(body2, email2, subject2)
	if err != nil {
		utilities.ErrorJSON(w, err)
	}

	utilities.WriteJSON(w, http.StatusOK, "Status Updated Successfully", "Success")
}

func ChangeComplaintStatusBySenate(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	var updatedComplaint models.Complaint
	err := json.NewDecoder(r.Body).Decode(&updatedComplaint)
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}

	err = models.ChangeComplaintStatus(id, "Approved By Senate")
	if err != nil {
		utilities.ErrorJSON(w, err)
	}
	body := fmt.Sprintf("Your revalidation request, concerning %s, which was assigned to %s has been approved by the senate. \n Log in to your portal to see more details", updatedComplaint.CourseConcerned, updatedComplaint.RespondingLecturer)
	student, err := models.GetUserByUserID(updatedComplaint.RequestingStudent)
	if err != nil {
		utilities.ErrorJSON(w, err)
	}
	email := student.Email
	subject := fmt.Sprintf("Update Concerning Revalidation Request for %s", updatedComplaint.CourseConcerned)
	err = SendEmail(body, email, subject)
	if err != nil {
		utilities.ErrorJSON(w, err)
	}

	utilities.WriteJSON(w, http.StatusOK, "Status Updated Successfully", "Success")
}

func GetStudentByID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	student, err := models.GetStudentById(id)
	if err != nil {
		fmt.Println("Unable to get student")
		utilities.ErrorJSON(w, err)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, student, "student")
}

func GetCoursesByStudentID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	student, err := models.GetStudentById(id)
	if err != nil {
		fmt.Println("Unable to get student")
		utilities.ErrorJSON(w, err)
		return
	}

	courses := student.Courses

	utilities.WriteJSON(w, http.StatusOK, courses, "courses")
}

func GetCoursesByStaffID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	lecturer, err := models.GetStaffById(id)
	if err != nil {
		fmt.Println("Unable to get student")
		utilities.ErrorJSON(w, err)
		return
	}

	courses := lecturer.CoursesTaken

	utilities.WriteJSON(w, http.StatusOK, courses, "courses")
}

func DeclineRequest(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	var updatedComplaint models.Complaint
	err := json.NewDecoder(r.Body).Decode(&updatedComplaint)
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}

	err = models.ChangeComplaintStatus(id, "Declined")
	if err != nil {
		utilities.ErrorJSON(w, err)
	}

	utilities.WriteJSON(w, http.StatusOK, "Status Updated Successfully", "Success")
}

func GetComplaintsForHOD(w http.ResponseWriter, r *http.Request) {
	complaints, err := models.GetComplaintsByStatus("Approved By Lecturer")
	if err != nil {
		fmt.Println("Unable to get complaints", err)
		utilities.ErrorJSON(w, err)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, complaints, "complaints")
}

func GetComplaintsForSenate(w http.ResponseWriter, r *http.Request) {
	complaints, err := models.GetComplaintsByStatus("Approved By HOD")
	if err != nil {
		fmt.Println("Unable to get complaints", err)
		utilities.ErrorJSON(w, err)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, complaints, "complaints")
}

func GetComplaintsByCourseCode(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")
	selectedCourse := r.URL.Query().Get("course")

	complaints, err := models.GetComplaintsByCourseCode(id, selectedCourse)
	if err != nil {
		fmt.Println("Unable to get complaints", err)
		utilities.ErrorJSON(w, err)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, complaints, "complaints")
}

func GetComplaintByCourseCode(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")
	selectedCourse := r.URL.Query().Get("course")

	complaint, err := models.GetComplaintByCourseCode(id, selectedCourse)
	if err != nil {
		fmt.Println("Unable to get complaint", err)
		utilities.ErrorJSON(w, err)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, complaint, "complaint")
}

func SendEmail(body string, recipient string, subject string) error {
	from := os.Getenv("EMAIL_SENDER")
	pass := os.Getenv("EMAIL_PASS")
	to := recipient

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return fmt.Errorf("smtp error: %s", err)
	}

	log.Print("sent, visit http://foobarbazz.mailinator.com")
	return err
}
