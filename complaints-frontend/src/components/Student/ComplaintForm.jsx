import React from 'react';
import {useState, useEffect } from 'react';
import axios from 'axios'

export default function ComplaintForm() {
  const [courses, setCourses] = useState([])
  const [isLoaded, setIsLoaded] = useState(false);
  const [error, setError] = useState(null)
  const [course_concerned, setCourseConcerned] = useState("")
  const [request_details, setRequestDetails] = useState("")
  const [test_score, setTestScore] = useState("")
  const [file, setFile] = useState(null)
  const [successMessage, setSuccessMessage] = useState("")
  const [errorMessage, setErrorMessage] = useState("")
  const userID = localStorage.getItem("userID")
  const token = localStorage.getItem("token")

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
  return response.json()
    })
    .then((json) => {
      setCourses(json.courses);
      setIsLoaded(true)
    })
    .catch((error) => {
      setIsLoaded(true);
    setError(error);
    })
  }, [token, userID])
  function handleSubmit(e) {
  e.preventDefault();

  const testScoreInt = parseInt(test_score, 10);
    if (isNaN(testScoreInt) || testScoreInt < 0 || testScoreInt > 30) {
      setErrorMessage("Test score must be between 0 and 30.");
      return;
    }

  const formData = new FormData()
  formData.append("course_concerned", course_concerned)
  formData.append("request_details", request_details)
  formData.append("test_score", testScoreInt)
  if (file) {
    formData.append("file", file)
  }

  // Log formData entries for debugging
  for (let pair of formData.entries()) {
    console.log(pair[0] + ': ' + pair[1]);
  }
  
  axios.post(`http://localhost:4000/complaint`, formData, {
    headers: {
      "Content-Type": "multipart/form-data",
      Authorization: token,
    },
  })
  .then((res) => {
    console.log(res)
    setSuccessMessage("Complaint submitted successfully")
    setCourseConcerned("")
    setTestScore("")
    setRequestDetails("")
    setFile(null)
  })
  .catch((err) => {
    if (err.response) {
      // The request was made and the server responded with a status code
      // that falls out of the range of 2xx
      console.log(err.response.data);
      console.log(err.response.status);
      console.log(err.response.headers);
      setErrorMessage("Failed to submit complaint. Server responded with: " + JSON.stringify(err.response.data));
    } else if (err.request) {
      // The request was made but no response was received
      console.log(err.request);
      setErrorMessage("Failed to submit complaint. No response received from server.");
    } else {
       // Something happened in setting up the request that triggered an Error
      console.log('Error', err.message);
      setErrorMessage("Failed to submit complaint. Error: " + err.message);
    }
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
      <label className="block mb-2"
      htmlFor='course_concerned'>Course Concerned</label>
      <br />
      
      <select
      className="block w-full border border-gray-300 rounded px-3 py-2 mb-4"
      name='course_concerned'
      id='course_concerned'
      onChange={(e) => setCourseConcerned(e.target.value)}
      value={course_concerned}
      >
      <option value="">--- SELECT COURSE ---</option>
      {courses.map((course) => 
      <option value= {course}>{course}</option>
      )}
      </select>

      <br />
      <br />

      <label className="block mb-2"
      htmlFor='request_details'>Details</label>
      <textarea
      id = "request_details"
      className="block w-full border border-gray-300 rounded px-3 py-2 mb-4"
      value={request_details}
      onChange={(e) => setRequestDetails(e.target.value)}
      />

      <br />
      <br />

      <label htmlFor='test_score' className="block mb-2">Test Score</label>
      <input
      type = "number"
      id = 'test_score'
      min="0"
      max="30"
      value = {test_score}
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

      <button type="submit" className="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600">Submit</button>
      </form>
    </div>
  ) 
}
}