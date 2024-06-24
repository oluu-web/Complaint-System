import React, { useState, useEffect, Fragment } from "react";
import { useParams, useNavigate } from "react-router-dom";
import '../Home.css';

const StudentHome = () => {
  const { id } = useParams();
  const [complaint, setComplaint] = useState({
    id: id,
    matricNo: null,
    details: null,
    student_proof: null,
    lecturer_proof: null,
    reason: null,
    status: null,
  });
  const [isLoaded, setIsLoaded] = useState(false);
  const [error, setError] = useState(null);
  const [courses, setCourses] = useState([]);
  const [selectedCourse, setSelectedCourse] = useState("");
  const navigate = useNavigate();
  const token = sessionStorage.getItem("token");
  const userID = sessionStorage.getItem("userID");

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

  useEffect(() => {
    if (selectedCourse) {
      fetch(`http://localhost:4000/student-complaint/${userID}?course=${selectedCourse}`, {
      headers: {
        Authorization: token,
      },
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
        console.log(json.complaint)
        if (json.complaint !== null) {
        setComplaint({
          matricNo: json.complaint.requesting_student,
          details: json.complaint.request_details,
          student_proof: `http://localhost:4000/${json.complaint.student_proof}`,
          lecturer_proof: `http://localhost:4000/${json.complaint.lecturer_proof}`,
          reason: json.complaint.reason,
          status: json.complaint.status,
        });
      } else {
        setComplaint(null)
      }
        setIsLoaded(true);
      })
      .catch((error) => {
        setIsLoaded(true);
        setError(error);
      });
    } else {
      setIsLoaded(true)
    }
    
  }, [selectedCourse, token, userID]);

  const handleNewComplaint = () => {
    navigate('/new-complaint');
  };

  if (error) {
    return <div className="text-red-500 text-center mt-4">Error: {error.message}</div>;
  } else if (!isLoaded) {
    return <p className="text-center mt-4">Loading...</p>;
  } else {
    return (
      <Fragment>
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-2xl font-bold">My Complaints</h1>
        </div>

            <div className="mb-4">
            <label htmlFor="course_select" className="block mb-2">Select Course</label>
            <select
              id="course_select"
              className="block w-full border border-gray-300 rounded px-3 py-2"
              value={selectedCourse}
              onChange={(e) => setSelectedCourse(e.target.value)}
            >
              <option value="">--- SELECT COURSE ---</option>
              {courses.map((course) => (
                <option key={course} value={course}>{course}</option>
              ))}
            </select>
          </div>

        {selectedCourse === "" ? (
            <div className="flex justify-center">
              <button
                className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                onClick={handleNewComplaint}
              >
                New Complaint
              </button>
            </div>
          ) : complaint !== null ? (
            <div>
              <button
                className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded flex justify-between items-center mb-8"
                onClick={handleNewComplaint}
              >
                New Complaint
              </button>
              <div>
                <h1 className="text-2xl font-bold mb-4">Complaint Details</h1>
                <p className="text-lg font-semibold mb-2">Matric Number: <span className="font-normal">{complaint.matricNo}</span></p>
                <div className="bg-gray-100 p-4 rounded-lg mb-4">
                  <h2 className="text-xl font-semibold mb-2">Details</h2>
                  <p className="mb-4">{complaint.details}</p>
                  {complaint.student_proof && <img src={complaint.student_proof} alt='complaint' className="max-w-full h-auto rounded-lg" />}
                  <br />
                  <h2 className="text-xl font-semibold mb-2">Approval Details</h2>
                  <p><span className="font-semibold">Status</span>: {complaint.status}</p>
                  <p className="mb-4"><span className="font-semibold">Reason</span>: {complaint.reason}</p>
                  {complaint.status !== "Pending" ? (<img src={complaint.lecturer_proof} className="max-w-full h-auto rounded-lg" alt="lecturer proof"/>): (<div></div>)}
                </div>
              </div>
            </div>
          ) : (
            <div className="flex justify-center">
              <button
                className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                onClick={handleNewComplaint}
              >
                New Complaint
              </button>
            </div>
          )}
        </div>
      </Fragment>
    );
  }
};

export default StudentHome;
