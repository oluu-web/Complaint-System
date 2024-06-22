import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

export default function ComplaintForm() {
  const [courses, setCourses] = useState([]);
  const [isLoaded, setIsLoaded] = useState(false);
  const [error, setError] = useState(null);
  const [courseConcerned, setCourseConcerned] = useState("");
  const [requestDetails, setRequestDetails] = useState("");
  const [testScore, setTestScore] = useState("");
  const [file, setFile] = useState(null);
  const [successMessage, setSuccessMessage] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const navigate = useNavigate()
  const userID = sessionStorage.getItem("userID");
  const token = sessionStorage.getItem("token");

  useEffect(() => {
    fetch(`http://localhost:4000/courses/${userID}`, {
      headers: {
        Authorization: token
      }
    })
      .then((response) => {
        if (response.status !== 200) {
          let err = new Error();
          err.message = "Invalid response code: " + response.status;
          throw err;
        }
        return response.json();
      })
      .then((json) => {
        setCourses(json.courses);
        setIsLoaded(true);
      })
      .catch((error) => {
        setIsLoaded(true);
        setError(error);
      });
  }, [token, userID]);

  function handleSubmit(e) {
    e.preventDefault();
    setIsSubmitting(true);  // Start the spinner

    const testScoreInt = parseInt(testScore, 10);
    if (isNaN(testScoreInt) || testScoreInt < 0 || testScoreInt > 30) {
      setErrorMessage("Test score must be between 0 and 30.");
      setIsSubmitting(false);  // Stop the spinner
      return;
    }

    const formData = new FormData();
    formData.append("course_concerned", courseConcerned);
    formData.append("request_details", requestDetails);
    formData.append("test_score", testScoreInt);
    if (file) {
      formData.append("file", file);
    }

    axios.post(`http://localhost:4000/complaint`, formData, {
      headers: {
        "Content-Type": "multipart/form-data",
        Authorization: token,
      },
    })
    .then((res) => {
      console.log(res);
      setSuccessMessage("Complaint submitted successfully");
      setCourseConcerned("");
      setTestScore("");
      setRequestDetails("");
      setFile(null);
      setIsSubmitting(false);  // Stop the spinner
      navigate("/student-dashboard")
    })
    .catch((err) => {
      if (err.response) {
        setErrorMessage("Failed to submit complaint. Server responded with: " + JSON.stringify(err.response.data));
      } else if (err.request) {
        setErrorMessage("Failed to submit complaint. No response received from server.");
      } else {
        setErrorMessage("Failed to submit complaint. Error: " + err.message);
      }
      setIsSubmitting(false);  // Stop the spinner
    });
  }

  if (error) {
    return <div>Error: {error.message}</div>;
  } else if (!isLoaded) {
    return <p>Loading...</p>;
  } else {
    return (
      <div className="mx-auto max-w-md">
        <h1 className="text-2xl font-bold mb-4">Submit a Complaint</h1>
        {successMessage && <p className="text-green-500">{successMessage}</p>}
        {errorMessage && <p className="text-red-500">{errorMessage}</p>}
        <form onSubmit={handleSubmit} encType='multipart/form-data'>
          <label className="block mb-2" htmlFor='course_concerned'>Course Concerned</label>
          <select
            className="block w-full border border-gray-300 rounded px-3 py-2 mb-4"
            name='course_concerned'
            id='course_concerned'
            onChange={(e) => setCourseConcerned(e.target.value)}
            value={courseConcerned}
          >
            <option value="">--- SELECT COURSE ---</option>
            {courses.map((course) => 
              <option key={course} value={course}>{course}</option>
            )}
          </select>

          <label className="block mb-2" htmlFor='request_details'>Details</label>
          <textarea
            id="request_details"
            className="block w-full border border-gray-300 rounded px-3 py-2 mb-4"
            value={requestDetails}
            onChange={(e) => setRequestDetails(e.target.value)}
          />

          <label htmlFor='test_score' className="block mb-2">Test Score</label>
          <input
            type="number"
            id='test_score'
            min="0"
            max="30"
            value={testScore}
            onChange={(e => setTestScore(e.target.value))}
            className="block w-full border border-gray-300 rounded px-3 py-2 mb-4"
          />

          <label htmlFor="file" className="block mb-2">Upload File</label>
          <input
            type="file"
            id="file"
            onChange={(e) => setFile(e.target.files[0])}
            className="block w-full border border-gray-300 rounded px-3 py-2 mb-4"
          />

          <button type="submit" className="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600">
            {isSubmitting ? (
              <svg
                className="animate-spin h-5 w-5 text-white"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  className="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  strokeWidth="4"
                ></circle>
                <path
                  className="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"
                ></path>
              </svg>
            ) : (
              "Submit"
            )}
          </button>
        </form>
      </div>
    );
  }
}
